import os
from typing import List, Optional

from loguru import logger
from opensearchpy import OpenSearch


class SearchService:
    """Business logic for search operations using OpenSearch."""

    INDEX_NAME = "posts"

    def __init__(self):
        opensearch_url = os.getenv("OPENSEARCH_URL", "http://localhost:9200")
        try:
            self.client = OpenSearch(
                hosts=[opensearch_url],
                use_ssl=False,
                verify_certs=False,
            )
            logger.info(f"Connected to OpenSearch at {opensearch_url}")
        except Exception as e:
            logger.warning(f"Failed to connect to OpenSearch: {e}")
            self.client = None

    async def search_posts(
        self,
        keyword: Optional[str] = None,
        hashtag: Optional[str] = None,
        author: Optional[str] = None,
        page: int = 1,
        limit: int = 10,
    ) -> List[dict]:
        """Search posts in OpenSearch."""
        if self.client is None:
            logger.warning("OpenSearch client not available")
            return []

        # Build the query
        must_clauses = []

        if keyword:
            must_clauses.append(
                {
                    "multi_match": {
                        "query": keyword,
                        "fields": ["content", "hashtags"],
                    }
                }
            )

        if hashtag:
            must_clauses.append({"match": {"hashtags": hashtag}})

        if author:
            must_clauses.append({"match": {"user_id": author}})

        if not must_clauses:
            query = {"match_all": {}}
        else:
            query = {"bool": {"must": must_clauses}}

        body = {
            "query": query,
            "from": (page - 1) * limit,
            "size": limit,
            "sort": [{"created_at": {"order": "desc"}}],
        }

        try:
            response = self.client.search(index=self.INDEX_NAME, body=body)
            hits = response.get("hits", {}).get("hits", [])
            return [hit["_source"] for hit in hits]
        except Exception as e:
            logger.error(f"Search error: {e}")
            return []

    async def index_post(self, post: dict) -> bool:
        """Index a post in OpenSearch."""
        if self.client is None:
            return False

        try:
            self.client.index(
                index=self.INDEX_NAME,
                id=post.get("id"),
                body=post,
            )
            logger.info(f"Indexed post: {post.get('id')}")
            return True
        except Exception as e:
            logger.error(f"Index error: {e}")
            return False

    def ensure_index(self):
        """Create the posts index with proper mapping if it doesn't exist."""
        if self.client is None:
            return

        mapping = {
            "mappings": {
                "properties": {
                    "id": {"type": "keyword"},
                    "user_id": {"type": "keyword"},
                    "content": {"type": "text", "analyzer": "standard"},
                    "hashtags": {"type": "keyword"},
                    "images": {"type": "keyword"},
                    "like_count": {"type": "integer"},
                    "created_at": {"type": "date"},
                    "updated_at": {"type": "date"},
                }
            }
        }

        try:
            if not self.client.indices.exists(index=self.INDEX_NAME):
                self.client.indices.create(index=self.INDEX_NAME, body=mapping)
                logger.info(f"Created index: {self.INDEX_NAME}")
        except Exception as e:
            logger.error(f"Failed to create index: {e}")
