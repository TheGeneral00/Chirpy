# PostalWarden — Detailed Data Flow

## Overview

The PostalWarden pipeline is fully channel-driven. Work moves through three sequential stages connected by buffered Go channels, with each stage operated by a config-bounded worker pool. No stage blocks another — the orchestrator manages ingestion and delivery while the pipeline processes asynchronously.

```
Client
  │
  │  POST /analyze  (HTTPS + API Token)
  ▼
Orchestrator — generates composite ID — holds connection
  │
  │ chan WorkItem (buffered)
  ▼
Enrichment Worker Pool
  │
  │ chan EnrichedMail (buffered)
  ▼
Rule Engine Judge Pool
  │
  │ chan Verdict (buffered)
  ▼
Orchestrator — matches composite ID — delivers response — logs event
  │
  ▼
Client receives { composite_id, verdict, reason, score }
```

---

## Stage 1 — Ingestion

### Trigger
A client POSTs a raw email (RFC 5322 / `.eml`) to the API with a valid API token in the request header.

### Steps

1. **Authentication** — the orchestrator validates the API token against the user database. Invalid or expired tokens are rejected immediately with a `401`.

2. **Composite ID generation** — a deterministic composite ID is derived from a hash of:
   ```
   HMAC( user_id + session_id + timestamp + request_fingerprint )
   ```
   This ID is reproducible from the same inputs, allowing a client to reclaim a verdict if the connection drops mid-processing.

3. **WorkItem construction** — the orchestrator wraps the raw mail and metadata into a `WorkItem`:
   ```go
   type WorkItem struct {
       CompositeID uuid.UUID      // deterministic reconnection key
       ProcessID   uuid.UUID      // unique pipeline trace ID
       ReceivedAt  time.Time
       UserID      uuid.UUID
       RawMail     []byte
       Metadata    map[string]any
   }
   ```

4. **Event log** — ingestion event is written to the event database, keyed by `ProcessID`.

5. **Queue** — the `WorkItem` is pushed onto the buffered `WorkItem` channel. If the channel is full (all workers busy and queue at depth limit), the orchestrator returns a `429 Too Many Requests` immediately rather than blocking the client indefinitely.

6. **Connection held** — the HTTP handler blocks on a per-request verdict channel, keyed by `CompositeID`, waiting for the pipeline to complete.

---

## Stage 2 — Enrichment Layer

### Worker Pool

On startup the orchestrator spawns N enrichment workers, where N is determined by the config file (absolute count, CPU percentage, or named profile). Workers are long-lived goroutines that loop on the `WorkItem` channel — they are not spawned per request.

```
WorkItem Channel
    ├── Enrichment Worker 1 (idle → picks up WorkItem)
    ├── Enrichment Worker 2 (busy)
    ├── Enrichment Worker 3 (idle → picks up WorkItem)
    └── Enrichment Worker N ...
```

### Per-Worker Enrichment Steps

Each worker runs all enrichment modules against a single `WorkItem` and assembles the results into an `EnrichedMail`. Modules that are independent run concurrently within the worker via goroutines, coordinated with a `sync.WaitGroup`.

**Modules:**

| Module | Integration | Output |
|---|---|---|
| C++ Regex Engine | HTTP/JSON microservice | Extracted URLs, suspicious patterns |
| Header Parser | In-process Go | Parsed routing headers, sender metadata |
| DKIM/SPF Validator | In-process Go | Validation results per sender domain |
| Attachment Hasher | In-process Go | SHA-256 hashes per attachment |
| *(future modules)* | — | — |

**EnrichedMail structure:**
```go
type EnrichedMail struct {
    CompositeID uuid.UUID
    ProcessID   uuid.UUID
    Headers     map[string]string
    From        string
    To          []string
    Subject     string
    Body        BodyParts
    Attachments []Attachment   // includes hashes
    Links       []string       // extracted by C++ engine
    Patterns    []string       // suspicious patterns from C++ engine
    SPFResult   string
    DKIMResult  string
}
```

### Output

The completed `EnrichedMail` is pushed onto the buffered `EnrichedMail` channel for consumption by the rule engine judge pool.

Event log: enrichment complete event written, keyed by `ProcessID`.

---

## Stage 3 — Rule Engine (The Courthouse)

### Judge Pool

The rule engine operates a pool of judge workers — independent goroutines each capable of evaluating an `EnrichedMail` to a verdict simultaneously. The judge pool size is config-bounded separately from the enrichment pool, allowing independent tuning of the two stages.

```
EnrichedMail Channel
    ├── Judge 1 (evaluating)
    ├── Judge 2 (idle → picks up EnrichedMail)
    ├── Judge 3 (evaluating)
    └── Judge N ...
```

Each judge is autonomous — it receives a case, evaluates it fully, and issues a verdict without coordination with other judges.

### Evaluation Steps (per judge)

1. **IoC Rule Evaluation** — the enriched dataset is tested against all loaded IoC rule sets. Each rule that fires contributes a weighted score and a verdict category.

   ```
   rule url_shortener_phish {
       ioc_list: known_shorteners
       pattern:  hash_bust_regex
       verdict:  phish
       weight:   0.8
   }
   ```

2. **YARA Evaluation** — YARA rules are run via `libyara` through cgo FFI against mail body and attachment content. Matches contribute additional weighted scores. The FFI surface is contained entirely within `internal/rules/` — no other package touches cgo directly.

3. **Domain Trust Scoring** — the sender domain's reputation grade is looked up from the trust store. The grade acts as a modifier on the aggregated IoC/YARA score:
   - High-trust domains require stronger signal to reach a flag/quarantine verdict
   - Unknown or low-trust domains are held to a stricter threshold
   - Each verdict issued updates the domain's grade via a weighted moving average

4. **Verdict Aggregation** — all weighted scores are combined into a final score. The score maps to a verdict category:

   | Score Range | Verdict |
   |---|---|
   | High | `allow` |
   | Medium | `flag` |
   | Low | `quarantine` |

5. **Verdict construction:**
   ```go
   type Verdict struct {
       CompositeID uuid.UUID `json:"composite_id"`
       ProcessID   uuid.UUID `json:"process_id"`
       Result      string    `json:"verdict"`   // allow, flag, quarantine
       Reason      string    `json:"reason"`
       Score       float64   `json:"score"`
   }
   ```

### Output

The completed `Verdict` is pushed onto the buffered `Verdict` channel.

Event log: verdict event written, keyed by `ProcessID`, including score and reason.

---

## Stage 4 — Response Delivery

### Normal Path (connection alive)

The orchestrator maintains a `sync.Map` of per-request verdict channels, keyed by `CompositeID`. When a verdict arrives on the `Verdict` channel, the orchestrator looks up the matching per-request channel and delivers the verdict to the waiting HTTP handler, which writes the response to the client.

```
Verdict arrives → lookup CompositeID in sync.Map
    │
    ├── found → send to per-request channel → HTTP handler responds → connection closed
    └── not found → connection dropped (see reconnection path below)
```

Response payload:
```json
{
  "composite_id": "...",
  "process_id":   "...",
  "verdict":      "flag",
  "reason":       "URL shortener combined with hash-busting pattern detected",
  "score":        0.74
}
```

### Dropped Connection Path

If the client connection drops while the pipeline is still processing, the waiting HTTP handler unblocks and removes the per-request channel from the `sync.Map`. The verdict, when it arrives, finds no waiting handler and is instead written to a **verdict buffer** — a short-lived in-memory store keyed by `CompositeID` with a configurable TTL.

```
Verdict arrives → lookup CompositeID → not found
    │
    └── write to verdict buffer (CompositeID → Verdict, TTL starts)
```

### Reconnection Path

A client that lost its connection can reconnect and present its `CompositeID`. The orchestrator checks the verdict buffer:

```
Client reconnects with CompositeID
    │
    ├── verdict in buffer → deliver immediately, remove from buffer
    ├── verdict not yet ready → re-register per-request channel, wait as normal
    └── TTL expired → direct client to query event DB by process_id
```

The event database is the durable fallback — all verdicts are persisted there regardless of delivery status, queryable by `ProcessID` at any time.

---

## Configuration

Worker pool sizes and queue depths are set via a YAML config file on startup.

```yaml
profile: balanced   # minimal | balanced | performance
                    # explicit values below override the profile

resources:
  enrichment_workers: 4        # absolute count
  enrichment_workers: 50%      # or relative to CPU cores
  rule_workers: 4
  rule_workers: 25%
  workitem_queue_depth: 200
  enrichedmail_queue_depth: 200
  verdict_buffer_ttl: 300      # seconds before expired verdicts are dropped
```

**Predefined profiles:**

| Profile | Enrichment Workers | Rule Workers | Queue Depth |
|---|---|---|---|
| `minimal` | 25% CPU | 10% CPU | 50 |
| `balanced` | 50% CPU | 25% CPU | 200 |
| `performance` | 75% CPU | 50% CPU | 1000 |

---

## Channel Summary

| Channel | Type | Direction | Purpose |
|---|---|---|---|
| `workItemCh` | `chan WorkItem` | Orchestrator → Enrichment Pool | Queues incoming mail for enrichment |
| `enrichedMailCh` | `chan EnrichedMail` | Enrichment Pool → Judge Pool | Passes fully enriched mail to rule evaluation |
| `verdictCh` | `chan Verdict` | Judge Pool → Orchestrator | Returns completed verdicts for delivery |
| per-request channel | `chan Verdict` | Orchestrator → HTTP handler | Delivers verdict to the specific waiting connection |

All pipeline channels are buffered. Depth is config-controlled. Backpressure propagates naturally — if the judge pool is saturated, `enrichedMailCh` fills; if enrichment is saturated, `workItemCh` fills; if `workItemCh` is full, ingestion returns `429`.

---

## Event Tracing

Every stage writes to the event database keyed by `ProcessID`:

| Event | Stage |
|---|---|
| `ingestion.received` | Ingestion |
| `ingestion.queued` | Ingestion |
| `enrichment.started` | Enrichment |
| `enrichment.complete` | Enrichment |
| `rules.started` | Rule Engine |
| `rules.verdict` | Rule Engine |
| `response.delivered` | Response |
| `response.buffered` | Response (dropped connection) |
| `response.redelivered` | Response (reconnection) |
| `response.expired` | Response (TTL exceeded) |
