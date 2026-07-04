# Deployment

This document covers running the stack locally via Docker Compose (the only supported deployment target today). Production deployment (Kubernetes manifests, secrets management, horizontal scaling, managed Postgres/Redis/object storage) is out of scope for Sprint 1 — see [ROADMAP.md](../../ROADMAP.md).

## Prerequisites

- Docker Desktop with Compose v2+ (`docker compose version`)
- ~10GB free disk for images/models (Ollama model + PyTorch/PaddleOCR wheels + base images)
- On Windows: give Docker Desktop enough resources (Settings → Resources) — the `embedding-service` and `ocr-service` images pull PyTorch/PaddlePaddle, which are memory-heavy to build

## First run checklist

1. `cp .env.example .env` — review and adjust values, especially the `*_PORT` variables if you already run PostgreSQL or other services locally (see [Port conflicts](#port-conflicts) below).
2. `docker compose up -d` — starts all 12 containers (5 infra + 7 app services).
3. **Create the MinIO bucket manually** (document-service does not auto-create it): open the MinIO console at `http://localhost:${MINIO_CONSOLE_PORT}` (default 9001, credentials from `.env`) and create a bucket matching `MINIO_BUCKET` (default `documents`). This is a one-time manual step until bucket auto-provisioning is added.
4. `make ollama-pull` — pulls the local LLM model (`OLLAMA_MODEL`, default `llama3.2`). This is a multi-GB download and intentionally decoupled from `docker compose up` so it doesn't block the rest of the stack from coming up.
5. Verify: `curl http://localhost:8080/health` should return `200`.

## Port conflicts

Two defaults were deliberately chosen to avoid colliding with commonly pre-installed local software:

- `POSTGRES_PORT` defaults to **5433** (not 5432) — many developer machines already run a local PostgreSQL service on 5432.
- `AUTH_SERVICE_PORT` defaults to **8087** (not 8081) — 8081 is a very common port for other local dev backends.

If you hit `Bind for 0.0.0.0:<port> failed: port is already allocated` on any other port, override the corresponding `*_PORT` variable in `.env` — the container's internal port (used for service-to-service calls on the Docker network) is unaffected by the host-side port you choose.

## Recreating containers after a config change

`docker-compose.yml` or `.env` changes affecting a single service only require recreating that service:

```bash
docker compose up -d <service-name>
```

Changes to shared infra (e.g. `POSTGRES_PORT`) may require recreating dependent services too, since Docker Compose doesn't always detect that a peer's connection string changed.

## Resetting state

```bash
make clean   # docker compose down -v --remove-orphans — wipes all volumes (Postgres, Redis, MinIO, Qdrant, Ollama model cache)
```

Use this when you want a truly clean slate (e.g. schema changes to `scripts/init-db/001_init.sql`, which only applies to a fresh Postgres volume).

## CI/CD

Three GitHub Actions workflows (`.github/workflows/`) gate merges — see [CONTRIBUTING.md](../../CONTRIBUTING.md):

| Workflow | Scope |
|---|---|
| `go-ci.yml` | `gofmt`, `go vet`, `go test`, `go build` across all 5 Go services (matrix) |
| `python-ci.yml` | `ruff check`, `pytest` across both Python services (matrix) — tests are health-endpoint-only, so CI never triggers a real PaddleOCR/bge-m3 model download |
| `docker-build.yml` | Build-only (no push/registry yet) for all 7 service images, using `docker/build-push-action` with GitHub Actions cache |

None of the three currently push images to a registry or deploy anywhere — that's a follow-up once a target environment (and its secrets) is decided.

## Known gaps (tracked for later sprints)

- No registry/image publishing step.
- No Kubernetes/production manifests (`deployment/` is still a placeholder).
- No automated MinIO bucket provisioning.
- No TLS termination — `api-gateway` serves plain HTTP; a reverse proxy (nginx/Traefik/cloud LB) would sit in front of it in production.
- No secrets management — `.env` holds plaintext credentials, fine for local dev, not for production.
