from datetime import datetime
from typing import List, Optional

from loguru import logger


class NotificationService:
    """Business logic for notification operations."""

    def __init__(self):
        # In-memory store for scaffold (replace with DB in production)
        self._notifications: List[dict] = []

    async def create_notification(
        self,
        user_id: str,
        notification_type: str,
        actor_id: str,
        post_id: Optional[str] = None,
        message: str = "",
    ) -> dict:
        """Create and store a new notification."""
        notification = {
            "id": str(len(self._notifications) + 1),
            "user_id": user_id,
            "type": notification_type,
            "actor_id": actor_id,
            "post_id": post_id,
            "message": message,
            "is_read": False,
            "created_at": datetime.utcnow().isoformat(),
        }
        self._notifications.append(notification)
        logger.info(f"Notification created: {notification_type} for user {user_id}")
        return notification

    async def get_notifications(
        self, user_id: str, page: int = 1, limit: int = 20
    ) -> List[dict]:
        """Get paginated notifications for a user."""
        user_notifications = [
            n for n in self._notifications if n["user_id"] == user_id
        ]
        # Sort by created_at descending
        user_notifications.sort(key=lambda x: x["created_at"], reverse=True)

        # Paginate
        start = (page - 1) * limit
        end = start + limit
        return user_notifications[start:end]

    async def mark_as_read(self, notification_id: str) -> bool:
        """Mark a notification as read."""
        for n in self._notifications:
            if n["id"] == notification_id:
                n["is_read"] = True
                return True
        return False
