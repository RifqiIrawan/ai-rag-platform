import uuid

from fastapi import APIRouter
from pydantic import BaseModel
from qdrant_client.models import Distance, PointStruct, VectorParams

from app.services.embedding_engine import get_embedding_model
from app.services.qdrant_client import get_qdrant_client

router = APIRouter()


class GenerateRequest(BaseModel):
    collection: str
    text: str


@router.post("/generate")
def generate_embedding(req: GenerateRequest):
    model = get_embedding_model()
    vector = model.encode(req.text).tolist()

    client = get_qdrant_client()
    if not client.collection_exists(req.collection):
        client.create_collection(
            collection_name=req.collection,
            vectors_config=VectorParams(size=len(vector), distance=Distance.COSINE),
        )

    point_id = str(uuid.uuid4())
    client.upsert(
        collection_name=req.collection,
        points=[PointStruct(id=point_id, vector=vector, payload={"text": req.text})],
    )

    return {"id": point_id, "collection": req.collection, "dimensions": len(vector)}
