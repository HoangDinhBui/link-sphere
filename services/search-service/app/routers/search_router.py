from fastapi import APIRouter
from pydantic import BaseModel
from typing import Optional

from app.services.search_service import SearchService

router = APIRouter()
search_service = SearchService()


class SearchRequest(BaseModel):
    """Search request body."""
    keyword: Optional[str] = None
    hashtag: Optional[str] = None
    author: Optional[str] = None
    page: int = 1
    limit: int = 10


@router.post("/posts")
async def search_posts(request: SearchRequest):
    """Search posts by keyword, hashtag, or author using OpenSearch."""
    results = await search_service.search_posts(
        keyword=request.keyword,
        hashtag=request.hashtag,
        author=request.author,
        page=request.page,
        limit=request.limit,
    )

    return {
        "code_status": 200,
        "message": f"Found {len(results)} results",
        "result": True,
        "errors": {},
        "data": results,
    }
