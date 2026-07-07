# ai-rag-platform

An enterprise RAG (Retrieval-Augmented Generation) platform that runs **entirely locally, with no AI API keys** — local LLM inference via [Ollama](https://ollama.com/), local embeddings via `bge-m3`, local OCR via `PaddleOCR`, and local vector search via [Qdrant](https://qdrant.tech/).

See [ROADMAP.md](ROADMAP.md) for the full build plan.

## Architecture

A Go/Python microservices monorepo fronted by a single API gateway:

| Service | Language | Port | Responsibility |
|---|---|---|---|
| api-gateway | Go | 8080 | Public entrypoint, reverse-proxies to downstream services |
| auth-service | Go | 8081 | User registration/login, JWT issuance (Postgres) |
| document-service | Go | 8082 | Document upload/metadata (Postgres + MinIO) |
| ocr-service | Python | 8083 | Text extraction from documents (PaddleOCR) |
| embedding-service | Python | 8084 | Vector embeddings (bge-m3) stored in Qdrant |
| rag-service | Go | 8085 | Retrieval (Qdrant) + generation (Ollama) orchestration |
| notification-service | Go | 8086 | Async notifications (Redis pub/sub + WebSocket fan-out) |
| frontend | React/Vite | 3000 | Web client (login, document upload, RAG chat, live notifications) |

Ports above are each service's internal container port (`PORT` env var, used for inter-service calls). Host-published ports come from `.env`'s `*_PORT` vars and may differ if a port is already taken locally (e.g. `AUTH_SERVICE_PORT` defaults to 8087, and `POSTGRES_PORT` to 5433, to avoid clashing with commonly pre-installed local services).

Infra: PostgreSQL, Redis, MinIO, Qdrant, Ollama — all run via Docker Compose.

Details: [docs/architecture/](docs/architecture/).

## Quick Start

```bash
cp .env.example .env
docker compose up -d
curl http://localhost:8080/health
open http://localhost:3000   # web client
```

Pull a local LLM model for rag-service (large download, run separately):

```bash
make ollama-pull
```

See [DEVELOPMENT.md](DEVELOPMENT.md) for local (non-Docker) development, and [CONTRIBUTING.md](CONTRIBUTING.md) for branch/PR conventions.

## Tech Stack

Go (Gin) · Python (FastAPI) · React (Vite, TypeScript, Tailwind) · PostgreSQL · Redis · MinIO · Qdrant · PaddleOCR · Ollama · bge-m3 · Docker Compose

## License

[MIT](LICENSE)
