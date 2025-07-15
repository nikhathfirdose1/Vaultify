# 🔐 Vaultify — Secure Secrets Management in Go

Vaultify is a lightweight, secure secrets manager built in **Go**, designed to simulate a production-ready system that stores and retrieves encrypted credentials for services and users. It supports **token-based access control**, **file-backed encrypted storage**, **PostgreSQL metadata tracking**, and **Prometheus-based observability**.

This project was built as part of a deep dive into distributed systems, backend infrastructure, and systems engineering — with alignment to real-world expectations at companies like Meta, Google, and Stripe.

---

## ✨ Features

- 🔐 Encrypted secret storage (AES-GCM)
- 🪪 Token-based access authorization
- 🗂️ PostgreSQL-backed metadata (TTL, versioning, audit trails)
- ⚙️ Configurable via VIM-editable YAML (`vaultify.yml`)
- 📈 Prometheus `/metrics` and `/healthz` endpoints
- 📜 Unix-style access logs (tail, grep friendly)
- 📁 File-backed secret blobs with secure key rotation (in progress)
- 🐍 Python-based benchmarking scripts (planned)

---

## 📦 Tech Stack

| Layer       | Tech            |
|-------------|-----------------|
| Language    | Go (1.22+)      |
| Database    | PostgreSQL      |
| Auth        | Token-based     |
| Encryption  | AES-256-GCM     |
| Config      | YAML (`gopkg.in/yaml.v3`) |
| Observability | Prometheus, logs |
| Tooling     | Docker, Bash, VIM, Git |

---

## 🧱 Architecture

```
Client ↔ HTTP API (Go)
               │
               ├── PostgreSQL: stores encrypted secrets, TTL, versioning, and audit logs
               ├── Token Auth: access gated by scoped tokens
               └── Prometheus: exposed system metrics via /metrics
```

---

## 🔧 Getting Started

### 📁 1. Clone the repo

```bash
git clone git@github.com:nikhathfirdose1/Vaultify.git
cd Vaultify
```

### 🚀 2. Build and run

```bash
go build -o vaultify ./cmd/vaultify
./vaultify --config ./config/vaultify.yml
```

> `vaultify.yml` contains server, DB, and encryption settings.

---

### 🧪 3. Store a secret (HTTP Example)

```bash
curl -X POST http://localhost:8080/store \
  -H "Authorization: Bearer YOUR-TOKEN" \
  -d '{"name": "API_KEY", "value": "sk-xyz-abc", "ttl": 86400}'
```

### 🔐 4. Fetch a secret

```bash
curl -X GET http://localhost:8080/fetch/API_KEY \
  -H "Authorization: Bearer YOUR-TOKEN"
```

---

## ⚙️ Configuration (`vaultify.yml`)

```yaml
server:
  port: 8080
  log_path: ./logs/access.log

encryption:
  key_path: ./config/master.key
  rotate_days: 30

database:
  host: localhost
  port: 5432
  user: vaultadmin
  password: securepass
  name: vaultdb
```

---

## 📈 Observability

- 🔍 `GET /metrics` — Prometheus scrape endpoint
- 💚 `GET /healthz` — basic service health check
- 📁 `logs/access.log` — structured logs in Unix format

---

## 📊 Benchmarks (planned)

- 🔬 Python scripts to simulate concurrent usage
- ⏱️ Track:
  - request throughput (ops/sec)
  - encryption latency
  - disk I/O bottlenecks

---

## 🛠️ Roadmap

- [x] Project scaffolding and initial commits
- [x] Config parsing and logging setup
- [ ] Encrypted blob writing + retrieval
- [ ] PostgreSQL schema for metadata
- [ ] Token auth middleware
- [ ] TTL enforcement + background cleanup
- [ ] Health and metrics endpoints
- [ ] Python benchmark tooling
- [ ] Docker Compose setup for local infra

---

## 🧠 Motivation

Vaultify was created to showcase production-grade backend skills through a real-world, security-first system. It emphasizes:

- Observability
- Systems resilience
- Secure file storage
- Config-driven architecture
- CLI-centric workflows (VIM, bash, logs)

---

## 🧑‍💻 Author

**Nikhath Firdose**  
📍 San Jose, CA  
[LinkedIn](https://linkedin.com/in/nikhath-firdose) | [GitHub](https://github.com/nikhathfirdose1)

---

## 📜 License

MIT License (coming soon)
