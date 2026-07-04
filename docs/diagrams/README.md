# Diagrams

Diagrams are written as [Mermaid](https://mermaid.js.org/) code blocks so they render directly on GitHub/GitLab without external tooling. See also the service dependency graph in [docs/architecture/README.md](../architecture/README.md).

## Auth: register + login

```mermaid
sequenceDiagram
    participant C as Client
    participant GW as api-gateway
    participant A as auth-service
    participant PG as PostgreSQL

    C->>GW: POST /api/v1/auth/register {email, password}
    GW->>A: proxy request
    A->>A: bcrypt.hash(password)
    A->>PG: INSERT INTO users (...)
    PG-->>A: id
    A-->>GW: 201 {id, email}
    GW-->>C: 201 {id, email}

    C->>GW: POST /api/v1/auth/login {email, password}
    GW->>A: proxy request
    A->>PG: SELECT id, password_hash WHERE email = ...
    PG-->>A: row
    A->>A: bcrypt.compare(password, hash)
    A->>A: sign JWT (HS256, exp)
    A-->>GW: 200 {token}
    GW-->>C: 200 {token}
```

## Document upload

```mermaid
sequenceDiagram
    participant C as Client
    participant GW as api-gateway
    participant D as document-service
    participant M as MinIO
    participant PG as PostgreSQL

    C->>GW: POST /api/v1/documents (multipart file, X-User-Id)
    GW->>D: proxy request
    D->>M: PutObject(bucket, objectKey, file)
    M-->>D: ok
    D->>PG: INSERT INTO documents (...)
    PG-->>D: id
    D-->>GW: 201 {id, filename, object_key}
    GW-->>C: 201 {id, filename, object_key}
```

## Embedding generation

```mermaid
sequenceDiagram
    participant C as Client
    participant E as embedding-service
    participant Q as Qdrant

    C->>E: POST /api/v1/embeddings/generate {collection, text}
    E->>E: bge-m3.encode(text) (lazy-loaded on first call)
    alt collection does not exist
        E->>Q: create_collection(size=1024, distance=Cosine)
    end
    E->>Q: upsert(point_id, vector, payload={text})
    Q-->>E: ok
    E-->>C: 200 {id, collection, dimensions}
```

## RAG query — planned pipeline (not yet implemented)

`rag-service`'s `/api/v1/rag/query` currently returns `501 Not Implemented`. This is the intended flow once the full pipeline lands:

```mermaid
sequenceDiagram
    participant C as Client
    participant GW as api-gateway
    participant R as rag-service
    participant E as embedding-service
    participant Q as Qdrant
    participant O as Ollama

    C->>GW: POST /api/v1/rag/query {query}
    GW->>R: proxy request
    R->>E: embed(query)
    E-->>R: query vector
    R->>Q: search(collection, query vector, limit)
    Q-->>R: top-k matches (payload + score)
    R->>R: assemble prompt (query + retrieved context)
    R->>O: /api/generate {model, prompt}
    O-->>R: generated response
    R-->>GW: 200 {answer, sources}
    GW-->>C: 200 {answer, sources}
```

## Infra healthcheck dependency order

```mermaid
flowchart TB
    subgraph Infra["Infra (no app dependencies)"]
        PG[(PostgreSQL)]
        REDIS[(Redis)]
        MINIO[(MinIO)]
        QDRANT[(Qdrant)]
        OLLAMA[(Ollama)]
    end

    AUTH[auth-service] -->|depends_on: healthy| PG
    DOC[document-service] -->|depends_on: healthy| PG
    DOC -->|depends_on: healthy| MINIO
    NOTIF[notification-service] -->|depends_on: healthy| REDIS
    EMB[embedding-service] -->|depends_on: healthy| QDRANT
    RAG[rag-service] -->|depends_on: healthy| QDRANT
    RAG -->|depends_on: healthy| OLLAMA
    RAG -.->|depends_on: started| EMB
    OCR[ocr-service] -.->|depends_on: started| DOC

    GW[api-gateway] -.->|depends_on: started| AUTH
    GW -.->|depends_on: started| DOC
    GW -.->|depends_on: started| RAG
    GW -.->|depends_on: started| NOTIF
```

Solid edges are `condition: service_healthy` (waits for the dependency's healthcheck to pass); dashed edges are `condition: service_started` (waits only for the process to start, since those dependencies expose liveness-only health today).
