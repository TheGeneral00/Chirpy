# PostalWarden — Architecture & Service Overview

## 1️⃣ Purpose

PostalWarden is an **experimental, training-focused email security system**.
It demonstrates:

* Multi-language systems architecture (Go, C++, Rust)
* Safe interoperability and FFI design
* Security-oriented design and decision-making
* Microservice vs library boundaries

**Note:** Not intended for production, commercial, or open-source deployment.

---

## 2️⃣ High-Level Architecture

```
Incoming Email / User Request
            │
            ▼
   ┌───────────────────┐
   │ Go Orchestrator   │
   │ - Chi router      │
   │ - Task scheduler  │
   │ - Resource mgmt   │
   └─────────┬─────────┘
             │ HTTP/JSON
             ▼
   ┌───────────────────┐
   │ C++ Regex Engine  │
   │ - Microservice    │
   │ - Extract URLs &  │
   │   suspicious text │
   │ - Return JSON     │
   └─────────┬─────────┘
             │ JSON results
             ▼
   ┌───────────────────┐
   │ Rust Rule Engine  │
   │ - FFI / libyara   │
   │ - Evaluate rules  │
   │ - Decide verdict  │
   └─────────┬─────────┘
             │ Verdict / Action
             ▼
        Go Orchestrator
             │
             ▼
     Action / Response to User
```

---

## 3️⃣ Services & Responsibilities

| Service / Component     | Language | Role / Responsibilities                                                                                                 | Boundary / Integration                   |
| ----------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------- | ---------------------------------------- |
| Go Orchestrator         | Go       | - HTTP endpoints, routing<br>- Task orchestration & scheduling<br>- Resource mgmt<br>- Calls C++ service & Rust library | Core process; orchestrates workflow      |
| C++ Regex Engine        | C++      | - Extract URLs & suspicious patterns<br>- High-performance regex evaluation                                             | Microservice; communicates via JSON/HTTP |
| Rust Rule Engine + YARA | Rust     | - Evaluate extracted patterns<br>- Determine allow/flag/quarantine<br>- Wrap libyara via FFI                            | In-process library; called via FFI       |
| Data & Task Flow        | JSON     | - Microservice communication (Go ↔ C++)<br>- FFI for Rust ↔ Go                                                          | Clear contracts, schema-driven           |

---

## 4️⃣ Directory / Project Structure (Go repo)

```
PostalWarden/
├── app/
│   ├── assets/
│   └── templates/
├── internal/
│   ├── auth/         # User auth, session management
│   ├── database/     # DB connections, queries
│   ├── helpers/      # Utility functions
│   ├── server/       # HTTP routing, endpoints
│   └── rules/        # Rust FFI wrapper, rule evaluation
├── log/              # Logs
├── sql/              # Schema + queries
└── README.md
```

**Notes:**

* Rust library could live inside `internal/rules/` early on, later extracted as a separate repo.
* C++ Regex Engine is a **standalone microservice** in a separate repo.

---

## 5️⃣ Data Flow / Integration

1. **Go orchestrator receives email** → validates input → creates a task.
2. **C++ Regex Engine** processes email → extracts URLs/patterns → returns JSON list.
3. **Rust Rule Engine** evaluates patterns via FFI + YARA rules → returns decision (allow, flag, quarantine).
4. **Go orchestrator** executes action → responds to user/system.
5. **Logging and metrics** are maintained throughout.

---

## 6️⃣ MVP Scope

* **Implemented / immediate focus:**

  * Go orchestrator skeleton + Chi router
  * C++ regex engine basic URL extraction
  * Rust library skeleton calling YARA via FFI

* **Next Steps / Phase 1:**

  * Go → C++ JSON pipeline
  * Rust FFI integration & rule evaluation
  * Task scheduling & orchestration in Go

* **Phase 2 / Optional Extensions:**

  * User dashboards & event logging
  * Email filtering, prioritization
  * Hot-reload rule sets, advanced YARA rules

**Non-goals:**

* Full mail server or commercial-grade deployment
* Security/performance production guarantees
* Large-scale multi-user features initially

---

## 7️⃣ Learning Outcomes

* Multi-language integration & FFI patterns (Go ↔ Rust ↔ C)
* Safe interaction with low-level libraries (YARA in Rust)
* Designing microservices vs libraries thoughtfully
* Orchestrating asynchronous tasks in Go
* Rule-based security detection pipelines
* Structuring code and documentation for clarity and maintainability

---

## ✅ Summary

PostalWarden is a **learning-focused, systems-oriented security project** designed to demonstrate architectural thinking, language interoperability, and safe security tooling design.

Even if incomplete, it **communicates strong technical judgment, problem-solving, and systems design skills** — precisely the kind of work that stands out on a resume or portfolio.

