# PostalWarden — Architecture & Service Overview

## 1️⃣ Purpose

PostalWarden is an **experimental, training-focused email security system**.
It demonstrates:

* Multi-language systems architecture (Go + C++)
* Safe interoperability via FFI and microservice design
* Security-oriented design and decision-making
* Deliberate architectural boundary decisions (FFI vs microservice)
* Rule-based detection pipelines with reputation scoring

**Note:** Not intended for production, commercial, or open-source deployment.

---

## 2️⃣ High-Level Architecture

```
┌──────────────────────────────────────────────────────┐
│                 Client / API Consumer                │
└─────────────────────────┬────────────────────────────┘
                          │ HTTPS + API Token
                          ▼
┌──────────────────────────────────────────────────────┐
│                  Go Orchestrator                     │
│        Ingestion · Auth · Worker Management          │
│              Event Tracing · Response                │
└────────┬─────────────────────────────────────────────┘
         │ Go Channel (WorkItem)
         ▼
┌──────────────────────────────────────────────────────┐
│                 Enrichment Layer                     │
│           Worker Pool (config-bounded)               │
│                                                      │
│   C++ Regex Engine · Header Parser · DKIM/SPF        │
│      Attachment Hasher · [ future modules ]          │
└────────┬─────────────────────────────────────────────┘
         │ Go Channel (EnrichedMail)
         ▼
┌──────────────────────────────────────────────────────┐
│               Go Rule Engine (FFI)                   │
│           Judge Pool (config-bounded)                │
│                                                      │
│      IoC Rules · YARA · Trust Scoring · Verdict      │
│           ── Final court — owns the verdict ──       │
└────────┬─────────────────────────────────────────────┘
         │ Go Channel (Verdict)
         ▼
┌──────────────────────────────────────────────────────┐
│                  Go Orchestrator                     │
│    Matches verdict by composite ID · Logs event      │
│          Delivers complete response to client        │
└──────────────────────────────────────────────────────┘
```

> For a precise description of each stage, channel, worker pool, connection handling, and reconnection model see [data_flow.md](./data_flow.md).

---

## 3️⃣ Ingestion & API Design

PostalWarden uses an **API-first ingestion model** rather than implementing SMTP/POP3 directly. This decision was made because:

* It enables clean statistics collection and auditability from the entry point
* It decouples the orchestrator from mail transport concerns
* It allows both manual uploads and automated pipeline integration
* It simplifies authentication and per-user token management

**Ingestion flow:**

1. User or system POSTs a raw email (RFC 5322 / `.eml`) to the API
2. Request is authenticated via API token linked to a user account
3. Orchestrator assigns a UUID to the work item and logs the ingestion event
4. A `WorkItem` is created and passed into the processing pipeline
5. On completion, the API responds with `{ verdict, process_id, reason }`

**Key data structures:**

```go
type WorkItem struct {
    ProcessID  uuid.UUID      `json:"process_id"`
    ReceivedAt time.Time      `json:"received_at"`
    UserID     uuid.UUID      `json:"user_id"`
    RawMail    []byte         `json:"raw_mail"`
    Metadata   map[string]any `json:"metadata"`
}

type EnrichedMail struct {
    ProcessID   uuid.UUID
    Headers     map[string]string
    From        string
    To          []string
    Subject     string
    Body        BodyParts
    Attachments []Attachment
    Links       []string
    Hashes      []string
    SPFResult   string
    DKIMResult  string
}

type Verdict struct {
    ProcessID uuid.UUID `json:"process_id"`
    Result    string    `json:"verdict"`  // allow, flag, quarantine
    Reason    string    `json:"reason"`
    Score     float64   `json:"score"`
}
```

---

## 4️⃣ Services & Responsibilities

| Service / Component  | Language | Role / Responsibilities                                                                                      | Boundary / Integration                              |
| -------------------- | -------- | ------------------------------------------------------------------------------------------------------------ | --------------------------------------------------- |
| Go Orchestrator      | Go       | HTTP endpoints, API token auth, task orchestration, event tracing, verdict response                          | Core process; orchestrates the full pipeline        |
| C++ Regex Engine     | C++      | High-performance URL and pattern extraction from email content                                               | Standalone microservice; communicates via HTTP/JSON |
| Go Rule Engine       | Go       | IoC rule evaluation, YARA integration via FFI, trust scoring, verdict determination                          | In-process library; called via cgo FFI              |
| PostgreSQL           | —        | User database (credentials, JWT hashes), event database (traced by UUID)                                     | Persistent storage for auth and audit trail         |

---

## 5️⃣ Integration Pattern Rationale

PostalWarden deliberately uses **two different integration patterns**, each chosen for specific reasons:

### C++ Regex Engine — Microservice (HTTP/JSON)

The regex engine is deployed as a standalone microservice because:

* It is a natural service boundary — extraction is a discrete, stateless operation
* Independent deployment and scaling is desirable for a CPU-intensive component
* HTTP/JSON is a clean, inspectable contract between services
* Language independence is a feature, not a constraint

### Go Rule Engine — In-Process FFI (cgo)

The rule engine is integrated via FFI rather than as a microservice because:

* YARA (`libyara`) is a C library and already requires cgo — FFI is the natural boundary
* The rule engine is tightly coupled to the verdict logic; in-process execution avoids serialization overhead and network latency on the critical decision path
* FFI demonstrates a distinct and complementary skill set to the microservice pattern
* The cgo surface is intentionally kept narrow — all FFI is contained within `internal/rules/`, exposing only a clean Go interface to the rest of the codebase

This deliberate mix of patterns demonstrates architectural judgement rather than uniform adoption of one approach.

---

## 6️⃣ Rule Engine & Trust Scoring

### IoC-Based Rules

Rules are defined as sets of related indicators of compromise (IoCs). Each rule bundles:

* A named list of known-bad values (e.g., URL shorteners, suspicious TLDs)
* A set of pattern conditions (e.g., hash-busting regex patterns)
* A verdict category on match (e.g., `spam`, `phish`, `suspicious`)
* A weight contributing to the overall score

Example rule concept:

```
rule url_shortener_phish {
    ioc_list: known_shorteners
    pattern:  hash_bust_regex
    verdict:  phish
    weight:   0.8
}
```

Multiple rules can fire against a single enriched mail. Weighted scores are aggregated by the rule engine to produce a final verdict.

### Domain Trust & Reputation Scoring

A trust store tracks known sender domains, validated via DKIM and SPF. Each root domain accumulates a reputation grade over time:

* Clean verdicts increase the domain's score
* Spam or phish verdicts decrease it
* A weighted moving average is used so historical data does not permanently bias newer senders
* The higher the grade, the more likely incoming mail from that domain is clean

Trust scoring is designed as an additive layer — rules fire independently of trust, but trust grade influences the final weighted verdict.

### YARA Integration

YARA rules are evaluated via `libyara` through cgo FFI. YARA handles content and binary pattern matching against mail body and attachment data. The FFI surface is contained entirely within `internal/rules/` and is not exposed to other packages.

---

## 7️⃣ Event Tracing & Observability

All significant pipeline events are traced using a custom logger implemented in Go. The logger can be dropped into any function as a one-liner and is designed to be wrapped or replaced later (e.g., with `slog`, `zerolog`, or an external sink) without touching business logic.

Every event is linked to a process UUID, enabling full auditability of a work item's journey through the pipeline. The event database in PostgreSQL records these traces, making per-rule firing events and pipeline stages queryable after the fact.

---

## 8️⃣ Directory / Project Structure

```
PostalWarden/
├── app/
│   ├── assets/
│   └── templates/
├── internal/
│   ├── auth/         # User auth, API token validation, JWT management
│   ├── database/     # DB connections, queries (user DB + event DB)
│   ├── helpers/      # Utility functions
│   ├── server/       # HTTP routing via Chi, ingestion endpoints
│   └── rules/        # Go Rule Engine — cgo FFI wrapper, YARA integration, IoC evaluation
├── log/              # Logs
├── sql/              # Schema + queries
└── README.md
```

**Notes:**

* `internal/rules/` is the sole owner of all FFI and cgo code. No other package interacts with libyara directly.
* The C++ Regex Engine lives in a separate repository as a standalone service.
* The user database holds credentials and JWT validation hashes. The event database documents pipeline events traced by UUID.

---

## 9️⃣ Data Flow

The pipeline moves data through three sequential stages — Ingestion, Enrichment, and Rule Evaluation — connected by Go channels and managed by config-bounded worker pools. Verdicts are matched back to waiting connections via a composite ID and delivered as a complete response.

> For the full stage-by-stage breakdown including channel design, worker pool behaviour, connection handling, and the reconnection model see [data_flow.md](./data_flow.md).

---

## 🔟 Implementation Status & Roadmap

### Completed
* User database (credentials, JWT hashes)
* Event database (UUID-traced pipeline events)
* Custom event logger (drop-in, wrappable)
* Chi HTTP router skeleton
* Basic user auth and session handling

### Phase 1 — Core Pipeline
* API ingestion endpoint (raw mail upload, token auth)
* C++ Regex Engine — URL and pattern extraction
* Go → C++ HTTP/JSON pipeline
* `EnrichedMail` struct and enrichment logic
* cgo FFI integration with libyara
* Basic IoC rule evaluation and verdict output

### Phase 2 — Rule Engine Depth
* Full IoC rule set definitions
* YARA rule authoring
* Domain trust store and reputation scoring
* Weighted verdict aggregation

### Phase 3 — Extensions (Optional)
* User dashboard and event log UI
* Hot-reload rule sets
* Advanced YARA rules
* Statistics and reporting layer

### Non-Goals
* Full mail server or SMTP/POP3 implementation
* Production-grade security or performance guarantees
* Large-scale multi-user deployment

---

## 1️⃣1️⃣ Learning Outcomes

* API-first design and HTTP ingestion pipeline
* Multi-language integration — Go orchestrating C++ and cgo FFI
* Deliberate architectural boundary decisions (microservice vs in-process FFI)
* Safe interaction with low-level C libraries (libyara via cgo)
* Rule-based security detection with IoC sets and YARA
* Reputation scoring and trust modelling
* Event tracing and audit trail design
* Structuring code and documentation for clarity and maintainability

---

## ✅ Summary

PostalWarden is a **learning-focused, systems-oriented email security project** that demonstrates architectural thinking, language interoperability, and security tooling design.

It deliberately uses two integration patterns — **microservice for the C++ regex engine** and **in-process FFI for the Go rule engine** — with clear technical justification for each. The API-first ingestion model, IoC-based rule engine, and domain reputation scoring reflect real-world anti-spam and anti-phish design patterns.

Even if incomplete, it communicates **strong technical judgement, deliberate design decisions, and systems-level thinking** — precisely the kind of work that stands out in security engineering and backend roles.
