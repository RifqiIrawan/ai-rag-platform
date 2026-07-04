# Database

## PostgreSQL

Schema is applied automatically on first container start via `scripts/init-db/001_init.sql`, mounted into `/docker-entrypoint-initdb.d/` in the `postgres` container (see `docker-compose.yml`). It only runs against a fresh, empty data volume — if you need to re-apply it after changing the file, remove the volume first (`docker compose down -v`, or `docker volume rm ai-rag-platform_postgres_data`) since Postgres's official image only runs init scripts once.

### `users` (owned by auth-service)

| Column | Type | Notes |
|---|---|---|
| `id` | `UUID` | PK, `gen_random_uuid()` (built into Postgres 13+, no extension needed) |
| `email` | `VARCHAR(255)` | `UNIQUE`, `NOT NULL` |
| `password_hash` | `VARCHAR(255)` | bcrypt hash, `NOT NULL` |
| `created_at` | `TIMESTAMPTZ` | default `now()` |
| `updated_at` | `TIMESTAMPTZ` | default `now()` — not currently auto-updated on row change; no trigger exists yet |

### `documents` (owned by document-service)

| Column | Type | Notes |
|---|---|---|
| `id` | `UUID` | PK, `gen_random_uuid()` |
| `owner_id` | `UUID` | FK → `users.id`, `ON DELETE CASCADE` |
| `filename` | `VARCHAR(512)` | original uploaded filename |
| `object_key` | `VARCHAR(1024)` | MinIO object key (`<uuid>-<filename>`), not the filename itself |
| `content_type` | `VARCHAR(255)` | nullable, from the upload's `Content-Type` |
| `size_bytes` | `BIGINT` | nullable |
| `status` | `VARCHAR(50)` | default `'uploaded'`; no state machine enforced yet (e.g. `processing`/`indexed` states are planned once OCR/embedding orchestration lands) |
| `created_at` / `updated_at` | `TIMESTAMPTZ` | default `now()` |

Index: `idx_documents_owner_id` on `documents(owner_id)`, since every read query filters by owner.

### Connection

Both `auth-service` and `document-service` connect via `DATABASE_URL` (standard `postgres://user:pass@host:port/db?sslmode=disable`), using `jackc/pgx/v5` connection pools. Pools are created eagerly at startup but do **not** ping on creation — this lets a service start even if Postgres is briefly unavailable; use each service's `/health/ready` to check actual connectivity (see [docs/architecture/README.md](../architecture/README.md#health-checks)).

## Qdrant

No fixed schema — collections are created **dynamically** by `embedding-service` on first write to a given collection name:

- Vector size: 1024 (bge-m3's output dimensionality)
- Distance metric: Cosine
- Payload: `{ "text": "<original text that was embedded>" }` — currently the only payload field; richer payloads (source document id, chunk index, page number) will be needed once document ingestion is wired end-to-end.

`rag-service` queries the same collections via Qdrant's `/collections/{name}/points/search` REST endpoint (see `services/rag-service/internal/clients/qdrant.go`).

There is currently no fixed convention for collection *naming* — the caller of `/api/v1/embeddings/generate` chooses the collection name per request. A per-tenant or per-document-type naming convention should be defined before this goes further (e.g. one collection per user, or one shared collection with a `owner_id` payload filter).

## Redis

No persisted data model — used purely as an ephemeral pub/sub bus. `notification-service` publishes to a single fixed channel, `notifications` (see `services/notification-service/internal/handlers/notifications.go`). There are no Redis-backed caches or queues yet.

## MinIO

Single bucket, configured via `MINIO_BUCKET` (default `documents`). Objects are stored under a flat key `<uuid>-<original filename>`; there's no bucket-creation step yet — the bucket must exist before `document-service` can write to it (this is currently a manual setup step; see [docs/deployment/README.md](../deployment/README.md#first-run-checklist)).
