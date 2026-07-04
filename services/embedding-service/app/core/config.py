import os
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    port: int = int(os.getenv("PORT", "8084"))
    qdrant_url: str = os.getenv("QDRANT_URL", "http://qdrant:6333")
    embedding_model: str = os.getenv("EMBEDDING_MODEL", "BAAI/bge-m3")


settings = Settings()
