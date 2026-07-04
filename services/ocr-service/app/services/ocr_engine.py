import threading

from app.core.config import settings

_lock = threading.Lock()
_ocr_instance = None


def get_ocr_engine():
    """Lazily initializes PaddleOCR on first real use.

    PaddleOCR downloads/loads model weights on init (can take 10-30s+),
    so this must never run at import time or the container would take
    too long to become healthy / fail its startup healthcheck.
    """
    global _ocr_instance
    if _ocr_instance is None:
        with _lock:
            if _ocr_instance is None:
                from paddleocr import PaddleOCR

                _ocr_instance = PaddleOCR(use_angle_cls=True, lang=settings.ocr_lang)
    return _ocr_instance
