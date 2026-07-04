import cv2
import numpy as np
from fastapi import APIRouter, File, UploadFile

from app.services.ocr_engine import get_ocr_engine

router = APIRouter()


@router.post("/extract")
async def extract_text(file: UploadFile = File(...)):
    contents = await file.read()
    img = cv2.imdecode(np.frombuffer(contents, np.uint8), cv2.IMREAD_COLOR)

    result = get_ocr_engine().ocr(img, cls=True)
    lines = [line[1][0] for block in (result or []) for line in block]

    return {"lines": lines}
