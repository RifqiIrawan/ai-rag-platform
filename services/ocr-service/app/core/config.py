import os
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    port: int = int(os.getenv("PORT", "8083"))
    document_service_url: str = os.getenv("DOCUMENT_SERVICE_URL", "http://document-service:8082")
    ocr_lang: str = os.getenv("OCR_LANG", "en")


settings = Settings()
