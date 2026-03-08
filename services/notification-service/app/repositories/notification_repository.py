from typing import Optional

from loguru import logger


class NotificationRepository:
    """Database operations for notifications.
    
    Currently uses in-memory storage. Replace with asyncpg/SQLAlchemy
    for production use.
    """

    def __init__(self):
        self._store: list = []

    async def create(self, notification: dict) -> dict:
        """Insert a notification record."""
        self._store.append(notification)
        return notification

    async def get_by_user_id(
        self, user_id: str, limit: int = 20, offset: int = 0
    ) -> list:
        """Get notifications for a user with pagination."""
        user_notifs = [n for n in self._store if n["user_id"] == user_id]
        user_notifs.sort(key=lambda x: x.get("created_at", ""), reverse=True)
        return user_notifs[offset : offset + limit]

    async def mark_read(self, notification_id: str) -> bool:
        """Mark a notification as read."""
        for n in self._store:
            if n.get("id") == notification_id:
                n["is_read"] = True
                return True
        return False
