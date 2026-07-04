from fastapi import FastAPI

from app.api.routes import health, ocr

app = FastAPI(title="ocr-service", version="0.1.0")

app.include_router(health.router)
app.include_router(ocr.router, prefix="/api/v1/ocr", tags=["ocr"])
