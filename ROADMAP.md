# Roadmap AI RAG Enterprise Platform (Tanpa API Key AI)

## Repository

**GitHub:** https://github.com/RifqiIrawan

### Nama Repository

`ai-rag-platform`

## Struktur Repository (Monorepo)

``` text
ai-rag-platform/
├── .github/workflows/
├── docs/
│   ├── architecture/
│   ├── api/
│   ├── database/
│   ├── deployment/
│   ├── diagrams/
│   └── roadmap/
├── services/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── document-service/
│   ├── ocr-service/
│   ├── embedding-service/
│   ├── rag-service/
│   └── notification-service/
├── frontend/
├── deployment/
├── scripts/
├── docker-compose.yml
├── Makefile
├── README.md
├── ROADMAP.md
├── DEVELOPMENT.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── LICENSE
├── .env.example
└── .gitignore
```

## Sprint 1 (Minggu 1--2)

### Inisialisasi Repository

-   Buat repository `ai-rag-platform`
-   Branch `main`
-   Branch `develop`
-   Tambahkan LICENSE
-   Tambahkan README.md
-   Tambahkan CONTRIBUTING.md
-   Tambahkan CODE_OF_CONDUCT.md
-   Tambahkan .gitignore
-   Tambahkan .env.example

### Setup Infrastruktur Lokal

-   Docker Desktop
-   Docker Compose
-   PostgreSQL
-   Redis
-   MinIO
-   Qdrant
-   Ollama

### Setup Service

-   api-gateway (Go)
-   auth-service (Go)
-   document-service (Go)
-   ocr-service (Python + PaddleOCR)
-   embedding-service (Python + bge-m3)
-   rag-service (Go)
-   notification-service (Go)

### CI/CD

-   GitHub Actions Go
-   GitHub Actions Python
-   Docker Build Workflow

## Teknologi

-   Go (Gin/Fiber)
-   Python (FastAPI)
-   PostgreSQL
-   Redis
-   MinIO
-   Qdrant
-   PaddleOCR
-   Ollama
-   bge-m3
-   Docker Compose

## Target Akhir

Platform RAG enterprise yang berjalan sepenuhnya secara lokal tanpa API
key AI dan siap dikembangkan ke production.
