import threading

from app.core.config import settings

_lock = threading.Lock()
_model_instance = None


def get_embedding_model():
    """Lazily loads the bge-m3 sentence-transformer on first real use.

    Downloads/loads multi-GB weights from Hugging Face on first call, so
    this must never run at import time or the container would fail its
    startup healthcheck / take too long to become ready.
    """
    global _model_instance
    if _model_instance is None:
        with _lock:
            if _model_instance is None:
                from sentence_transformers import SentenceTransformer

                _model_instance = SentenceTransformer(settings.embedding_model)
    return _model_instance
