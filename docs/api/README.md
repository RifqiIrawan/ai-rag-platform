# API Reference

All endpoints below can be called either directly against a service's own port, or through `api-gateway` (port 8080) under the path prefixes it proxies. Direct ports are useful for local debugging; the gateway is the intended entrypoint in normal use.

| Path prefix (via gateway :8080) | Proxied to |
|---|---|
| `/api/v1/auth/*` | auth-service |
| `/api/v1/documents/*` | document-service |
| `/api/v1/rag/*` | rag-service |
| `/api/v1/notifications/*` | notification-service |

`ocr-service` and `embedding-service` are not currently proxied through the gateway — call them directly on their own ports (8083, 8084). This will likely change once the document ingestion pipeline is wired end-to-end.

Every service also exposes `GET /health` (liveness) and, where it has a real dependency, `GET /health/ready` (readiness — returns `503` if the dependency is unreachable). These are omitted from the tables below; see [docs/architecture/README.md](../architecture/README.md#health-checks) for the full readiness matrix.

---

## auth-service (:8081)

### `POST /api/v1/auth/register`

Request:
```json
{ "email": "user@example.com", "password": "at-least-8-chars" }
```

Response `201`:
```json
{ "id": "6fb58030-...", "email": "user@example.com" }
```

Response `409` if the email is already registered. Response `400` if `email`/`password` fail validation (`password` requires 8+ characters).

### `POST /api/v1/auth/login`

Request:
```json
{ "email": "user@example.com", "password": "at-least-8-chars" }
```

Response `200`:
```json
{ "token": "eyJhbGciOiJIUzI1NiIs..." }
```

JWT is HS256-signed, contains `sub` (user id) and `exp`, expires after `JWT_EXPIRY_HOURS` (default 24). Response `401` on bad credentials.

---

## document-service (:8082)

Both endpoints currently require an `X-User-Id` header (the owning user's UUID) — there is no auth middleware wired yet to derive this from the JWT automatically; that integration is a follow-up sprint item.

### `POST /api/v1/documents`

`multipart/form-data` with a `file` field. Stores the file in MinIO and its metadata in Postgres.

Response `201`:
```json
{ "id": "...", "filename": "report.pdf", "object_key": "<uuid>-report.pdf" }
```

### `GET /api/v1/documents`

Lists documents owned by `X-User-Id`, most recent first.

Response `200`:
```json
{
  "documents": [
    {
      "id": "...",
      "filename": "report.pdf",
      "content_type": "application/pdf",
      "size_bytes": 20481,
      "status": "uploaded",
      "created_at": "2026-07-05T00:00:00Z"
    }
  ]
}
```

---

## ocr-service (:8083, direct only)

### `POST /api/v1/ocr/extract`

`multipart/form-data` with a `file` field (image). Runs PaddleOCR (lazy-loaded on first call — expect a one-time delay while weights load).

Response `200`:
```json
{ "lines": ["Invoice #1042", "Total: $84.20"] }
```

---

## embedding-service (:8084, direct only)

### `POST /api/v1/embeddings/generate`

Request:
```json
{ "collection": "documents", "text": "the text to embed" }
```

Encodes `text` with `BAAI/bge-m3` (lazy-loaded on first call) and upserts the resulting vector into the named Qdrant collection, creating the collection (1024-dim, cosine distance) if it doesn't already exist.

Response `200`:
```json
{ "id": "<point-uuid>", "collection": "documents", "dimensions": 1024 }
```

---

## rag-service (:8085)

### `POST /api/v1/rag/query`

Request:
```json
{ "query": "What does the Q3 report say about revenue?" }
```

**Not yet implemented.** Currently always returns `501`:
```json
{ "message": "rag pipeline not yet implemented", "query": "..." }
```

See [docs/diagrams/README.md](../diagrams/README.md#rag-query--planned-pipeline-not-yet-implemented) for the intended retrieval+generation flow (embed query → search Qdrant → prompt Ollama).

---

## notification-service (:8086)

### `POST /api/v1/notifications`

Request:
```json
{ "message": "your document finished processing" }
```

Publishes `message` to the shared Redis pub/sub channel `notifications`. Response `202`:
```json
{ "status": "published" }
```

Per-user targeting, delivery channels (WebSocket/email), and persistence are not implemented yet — this is a broadcast-only publish today.

---

## api-gateway (:8080)

Exposes only `GET /health` itself; all other traffic is proxied per the table at the top of this document via `httputil.ReverseProxy`, path and method preserved.
