# Development Guide

## Prerequisites

- Docker Desktop (Compose v2+)
- Go 1.26+
- Python 3.12+
- Git
- (Windows) GNU Make via Git Bash or WSL — native `cmd`/PowerShell does not ship `make`

## Running the full stack

```bash
cp .env.example .env
docker compose up -d --build
docker compose ps        # all services should show "healthy"
```

Check each service directly:

```bash
curl http://localhost:8080/health   # api-gateway
curl http://localhost:8087/health   # auth-service (host port 8087; container port stays 8081)
curl http://localhost:8082/health   # document-service
curl http://localhost:8083/health   # ocr-service
curl http://localhost:8084/health   # embedding-service
curl http://localhost:8085/health   # rag-service
curl http://localhost:8086/health   # notification-service
```

Or via the gateway proxy: `curl http://localhost:8080/api/v1/auth/...` etc.

## Pulling a local LLM model

The Ollama container does **not** auto-pull a model (multi-GB download, decoupled from `docker compose up`):

```bash
make ollama-pull            # pulls $OLLAMA_MODEL (default: llama3.2)
docker compose exec ollama ollama list
```

## Running a single Go service locally (outside Docker)

```bash
cd services/api-gateway
go run ./cmd/server
```

## Running a single Python service locally (outside Docker)

```bash
cd services/ocr-service
python -m venv .venv && .venv\Scripts\activate   # Windows
pip install -r requirements.txt -r requirements-dev.txt
uvicorn app.main:app --reload --port 8083
```

## Tests / lint

```bash
make test-go        # go vet + go test across all 5 Go services
make test-python     # pytest across both Python services
make lint-go         # gofmt check
make lint-python     # ruff check
```

## Troubleshooting

- **Healthcheck failing on a service**: `docker compose logs <service>` — most common causes are a downstream dependency (Postgres/Redis/MinIO/Qdrant) not yet ready, or a missing env var.
- **Port already in use**: another local process is bound to one of 5433/6379/9000/9001/6333/6334/11434/8080-8086. Stop it or override the `*_PORT` var in `.env`. Postgres defaults to host port 5433 (not 5432) specifically to avoid clashing with a locally installed PostgreSQL instance.
- **`exec format error` in a container**: usually CRLF line endings from a Windows checkout of a shell script/Dockerfile. `.gitattributes` forces LF for these; re-clone or `git add --renormalize .` if you hit this.
- **PaddleOCR/bge-m3 first-call latency**: both engines lazy-load on first real request (not at container startup), so the *first* `/extract` or `/generate` call after a fresh container start will be slow while weights download/load. `/health` itself is always instant.
