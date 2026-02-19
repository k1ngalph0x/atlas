# Atlas

A distributed error monitoring platform built as a hands-on study of distributed systems concepts — Kafka, RabbitMQ, microservices, and local LLM inference.

Atlas ingests application errors via an SDK, deduplicates them into issues by fingerprint, fires configurable alerts, and generates AI remediation insights using a locally-run language model.

---

## Architecture

```
Application + SDK
      |
      | HTTP (X-API-Key)
      v
ingestion-service  :8081
      |
      | Kafka: atlas-events
      v
issue-service      :8082  — fingerprints and deduplicates events into issues
      |
      | Kafka: issue-updates
      |
      +-------------------+
      |                   |
      v                   v
alert-service  :8084    intelligence-service  :8083
fires rules             RabbitMQ -> Ollama -> insights

identity-service  :8080  — auth (JWT), orgs, projects, API key management
```

---

## Design Highlights

- Event-driven architecture with Kafka as the messaging backbone
- Fingerprint-based deduplication (SHA-256) — identical errors consolidate into one issue with an incrementing count
- Idempotent issue creation via a unique index on `(project_id, fingerprint)`
- Asynchronous AI insight generation via RabbitMQ workers — ingestion is never blocked by inference latency
- API key hashing (SHA-256) with one-time reveal — the platform never stores raw keys
- Shared Postgres instance with per-service table separation

---

## Tech Stack

| Layer           | Technology                        |
| --------------- | --------------------------------- |
| Language        | Go                                |
| HTTP            | Gin                               |
| Database        | PostgreSQL + GORM                 |
| Event streaming | Apache Kafka (segmentio/kafka-go) |
| Work queue      | RabbitMQ (amqp091-go)             |
| AI inference    | Ollama — llama3.2:3b (local)      |
| Frontend        | React, Vite, Tailwind CSS         |

---

## Services

| Service              | Port | Responsibility                                                   |
| -------------------- | ---- | ---------------------------------------------------------------- |
| identity-service     | 8080 | JWT auth, organizations, projects, API key management            |
| ingestion-service    | 8081 | Validates API key, publishes events to Kafka                     |
| issue-service        | 8082 | Fingerprints events, deduplicates into issues, exposes issue API |
| intelligence-service | 8083 | Generates AI insights via Ollama for qualifying issues           |
| alert-service        | 8084 | Evaluates alert rules, stores and exposes alert logs             |

---

## Prerequisites

- Go 1.21+
- Node.js 18+
- Docker (Kafka + RabbitMQ)
- Ollama running locally with llama3.2:3b

```bash
ollama pull llama3.2:3b
```

---

## Running

**1. Start Kafka and RabbitMQ**

```bash
docker-compose up -d
```

**2. Configure environment**

Each service reads from its own `.env` file. Example:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=atlas_identity
DB_USERNAME=postgres
DB_PASSWORD=yourpassword

KAFKA_BROKERS=localhost:9092

JwtKey=yoursecretkey

RABBIT_USER=guest
RABBIT_PASSWORD=guest
RABBIT_HOST=localhost
RABBIT_PORT=5672

OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.2:3b
```

**3. Start all services**

```bash
make
```

Each service opens in a separate terminal window. To run individually:

```bash
make identity
make ingestion
make issue
make alert
make intelligence
```

**4. Start the frontend**

```bash
cd frontend
cp .env.example .env.development
npm install
npm run dev
```

---

## SDK

```go
import atlas "github.com/k1ngalph0x/atlas-go-sdk"

client := atlas.NewClient("atlas_yourprojectapikey",
    atlas.WithBaseURL("http://localhost:8081"),
)

client.CaptureError(err)
client.CaptureMessage("something happened", "warning")

// Gin middleware — captures panics automatically
router.Use(client.GinMiddleware())
```

API keys are shown once on project creation. The platform stores only a SHA-256 hash.

---

## Alert Rules

Configured per project via the dashboard or API (`POST /projects/:id/rules`).

| Condition         | Behaviour                                          |
| ----------------- | -------------------------------------------------- |
| `new_issue`       | Fires when a new unique fingerprint is detected    |
| `critical_error`  | Fires on any error or critical level event         |
| `count_threshold` | Fires when an issue exceeds a set occurrence count |

---

## Project Structure

```
atlas/
├── services/
│   ├── identity-service/
│   ├── ingestion-service/
│   ├── issue-service/
│   ├── alert-service/
│   └── intelligence-service/
├── shared/
│   └── models/          # shared Go types (Project, etc.)
├── sdk/                 # atlas-go-sdk
├── frontend/
├── test-app/            # sample app for local testing
└── Makefile
```

---

## Notes

- AI insights are generated asynchronously. The frontend polls until the insight is ready or times out after 60 seconds.
- Issue deduplication uses SHA-256 of the error message. The same error always maps to the same issue regardless of stack trace differences between occurrences.
- The intelligence service only processes issues with level `error` or `critical`. Low-volume `warning` level events are skipped.
- All services share one Postgres instance with separate databases per service.
