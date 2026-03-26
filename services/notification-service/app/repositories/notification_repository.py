from typing import List, Optional
import asyncpg

class NotificationRepository:
    def __init__(self, db_pool: asyncpg.Pool):
        self.pool = db_pool

    async def create(self, user_id: str, actor_id: str, type: str, message: str, post_id: Optional[str] = None) -> dict:
        query = """ 
        INSERT INTO notifications (user_id, actor_id, type, message, post_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, user_id, actor_id, type, post_id, message, is_read, created_at
        """
        # Returning will help us get the notification id and send it to the user, no need to query again
        async with self.pool.acquire() as conn:
            record = await conn.fetchrow(query, user_id, actor_id, type, message, post_id)
            return dict(record) if record else {}

    
    async def get_by_user(self, user_id: str, page: int = 1, limit: int = 20) -> List[dict]:
        offset = (page - 1) * limit
        query = """ 
        SELECT id, user_id, actor_id, type, post_id, message, is_read, created_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
        """

        async with self.pool.acquire() as conn:
            records = await conn.fetch(query, user_id, limit, offset)
            return [dict(r) for r in records]

        
    async def mark_as_read(self, notification_id: str) -> bool:
        query = "UPDATE notifications SET is_read = TRUE WHERE id = $1"

        async with self.pool.acquire() as conn:
            result = await conn.execute(query, notification_id)
            return result == "UPDATE 1"

    


