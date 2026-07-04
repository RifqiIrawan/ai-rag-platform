from fastapi import FastAPI

from app.api.routes import embeddings, health

app = FastAPI(title="embedding-service", version="0.1.0")

app.include_router(health.router)
app.include_router(embeddings.router, prefix="/api/v1/embeddings", tags=["embeddings"])
