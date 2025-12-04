# Project To-Do List

## Feature Overview (Sorted by Priority)

| Done | Priority | Feature | Description |
|------|----------|---------|-------------|
| [ ] | 5 | [Folder structure](#folder-structure) | Base project layout and module separation |
| [ ] | 5 | [Config system](#config-system) | Central configuration with env & file support |
| [ ] | 5 | [JWT auth](#jwt-auth) | Secure authentication with tokens |
| [ ] | 5 | [Routing setup](#routing-setup) | Router initialization, groups, middleware |
| [ ] | 5 | [DB migrations](#db-migrations) | Versioned schema management |
| [ ] | 5 | [Panic recovery middleware](#panic-recovery-middleware) | Prevent server crashes, return safe errors |
| [ ] | 5 | [Structured logging](#structured-logging) | Consistent, machine-parseable logs |
| [ ] | 5 | [Central error handling](#central-error-handling) | Unified error responses & wrapping |
| [ ] | 5 | [README](#readme) | Root documentation |
| [ ] | 4 | [Request validation](#request-validation) | Input checking & sanitization |
| [ ] | 4 | [Request logging middleware](#request-logging-middleware) | Per-request log details |
| [ ] | 4 | [Correlation ID middleware](#correlation-id-middleware) | Trace requests end-to-end |
| [ ] | 4 | [File log output](#file-log-output) | Redirect logs to files |
| [ ] | 4 | [DB event logging](#db-event-logging) | Write user events to database |
| [ ] | 4 | [App error types](#app-error-types) | Custom domain error structures |
| [ ] | 4 | [Linters](#linters) | golangci-lint baseline |
| [ ] | 4 | [Graceful shutdown](#graceful-shutdown) | Safe shutdown on signals |
| [ ] | 4 | [API docs](#api-docs) | HTTP documentation (OpenAPI/Swagger) |
| [ ] | 3 | [Connection pool tuning](#connection-pool-tuning) | Configure pgx pool settings |
| [ ] | 3 | [Log rotation](#log-rotation) | Auto-rotate log files |
| [ ] | 3 | [Unit tests](#unit-tests) | Test isolated business logic |
| [ ] | 3 | [Staticcheck](#staticcheck) | Advanced static analysis |
| [ ] | 3 | [gosec](#gosec) | Security static analysis |
| [ ] | 3 | [govulncheck](#govulncheck) | Vulnerability audit |
| [ ] | 3 | [Metrics extension](#metrics-extension) | Add custom counters/histograms |
| [ ] | 3 | [Docker support](#docker-support) | Containerize the application |
| [ ] | 2 | [DB health check](#db-health-check) | Ping database endpoint |
| [ ] | 2 | [Rate limiting](#rate-limiting) | Protect endpoints from abuse |
| [ ] | 2 | [Integration tests](#integration-tests) | Test API + DB together |
| [ ] | 2 | [docker-compose](#docker-compose) | Local infra stack |
| [ ] | 2 | [ERD diagram](#erd-diagram) | Schema overview |
| [ ] | 1 | [Load tests](#load-tests) | Performance benchmarking |
| [ ] | 1 | [Distributed tracing](#distributed-tracing) | OpenTelemetry traces |


# Extended Feature Descriptions

## Folder structure
**Priority:** 5  
A clean modular layout.

**Subtasks**
- Create `/cmd`, `/internal`, `/pkg`
- Separate handler, service, repository layers
- Setup `/migrations`, `/configs`, `/scripts`

---

## Config system
**Priority:** 5  

**Subtasks**
- Load from `.env`
- Override via CLI flags
- Add structured config type
- Validate config on startup

---

## JWT auth
**Priority:** 5  

**Subtasks**
- Token generation
- Token validation middleware
- Refresh token flow
- Role/permission integration

---

## Routing setup
**Priority:** 5  

**Subtasks**
- Create router groups (`/auth`, `/api/v1`)
- Attach middlewares
- Add 404 and MethodNotAllowed handlers

---

## DB migrations
**Priority:** 5  

**Subtasks**
- Create `sql/migrations`
- Setup migrate tool integration
- Auto-run on startup (optional)

---

## Panic recovery middleware
**Priority:** 5  

**Subtasks**
- Recover from panics
- Log panic stack trace
- Return safe JSON 500

---

## Structured logging
**Priority:** 5  

**Subtasks**
- Use zerolog / slog / zap
- Add timestamp, request ID, latency fields
- Provide log wrappers

---

## Central error handling
**Priority:** 5  

**Subtasks**
- Convert internal errors → API errors
- Add error codes
- JSON error format
- Attach correlation IDs

---

## README
**Priority:** 5  

**Subtasks**
- Setup instructions  
- Dev workflow  
- Tech stack  
- Example requests  

---

## Request validation
**Priority:** 4  

**Subtasks**
- Validate payloads w/ structs
- Add sanitization (SQLi, XSS)
- Add reusable validator module

---

## Request logging middleware
**Priority:** 4  

**Subtasks**
- Log method, path, status, latency
- Mask sensitive data
- Include correlation ID

---

## Correlation ID middleware
**Priority:** 4  

**Subtasks**
- Generate UUID if missing
- Add to logs
- Return in response header

---

## File log output
**Priority:** 4  

**Subtasks**
- Configurable log file path
- Redirect structured logs to files
- Fail gracefully if path invalid

---

## DB event logging
**Priority:** 4  

**Subtasks**
- Insert user events into table
- Design event schema
- Add async buffer (optional)

---

## App error types
**Priority:** 4  

**Subtasks**
- Define typed errors
- Wrap DB errors
- Attach metadata (field errors, etc.)

---

## Linters
**Priority:** 4  

**Subtasks**
- Configure golangci-lint
- Enable vet, unused, style issues
- Add CI step

---

## Graceful shutdown
**Priority:** 4  

**Subtasks**
- Capture SIGINT/SIGTERM
- Close DB pool
- Finish active requests

---

## API docs
**Priority:** 4  

**Subtasks**
- OpenAPI spec
- Swagger UI
- Auto-generated handlers (optional)

---

## Connection pool tuning
**Priority:** 3  

**Subtasks**
- Configure min/max pool size  
- Tune max lifetime  
- Add metrics for pool usage  

---

## Log rotation
**Priority:** 3  

**Subtasks**
- File rotation by size/time
- Log compression
- Cleanup policy

---

## Unit tests
**Priority:** 3  

**Subtasks**
- Test handlers with mocks
- Test service logic
- Test validation

---

## Staticcheck
**Priority:** 3  

**Subtasks**
- Enable staticcheck rules
- Fix dead code
- Fix shadow vars and misuse patterns

---

## gosec
**Priority:** 3  

**Subtasks**
- Check for crypto misuse  
- Sensitive file permissions  
- Insecure random usage  

---

## govulncheck
**Priority:** 3  

**Subtasks**
- Scan for vulnerable dependencies  
- Add CI step  

---

## Metrics extension
**Priority:** 3  

**Subtasks**
- Add request duration histograms  
- Add error counters  
- Add DB metrics  

---

## Docker support
**Priority:** 3  

**Subtasks**
- Create Dockerfile  
- Multi-stage build  
- Small image size (Alpine/distroless)  

---

## DB health check
**Priority:** 2  

**Subtasks**
- `/health/db` endpoint  
- Pool ping checks  

---

## Rate limiting
**Priority:** 2  

**Subtasks**
- Token bucket  
- Global and per-IP  
- Configurable limits  

---

## Integration tests
**Priority:** 2  

**Subtasks**
- Spin up DB  
- Test real HTTP requests  
- Test error cases  

---

## docker-compose
**Priority:** 2  

**Subtasks**
- Run Postgres locally  
- Include adminer/pgadmin  
- Mount logs  

---

## ERD diagram
**Priority:** 2  

**Subtasks**
- Draw tables + relations  
- Add to README  
- Export PNG/SVG  

---

## Load tests
**Priority:** 1  

**Subtasks**
- k6 or vegeta scripts  
- Add scenario tests  
- SLA baseline  

---

## Distributed tracing
**Priority:** 1  

**Subtasks**
- OpenTelemetry integration  
- Trace IDs in logs  
- Export to Jaeger/Tempo  

