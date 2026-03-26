from typing import List, Optional

from loguru import logger
from app.repositories.notification_repository import NotificationRepository


class NotificationService:
    def __init__(self):
        self.repo: Optional[NotificationRepository] = None

    async def create_notification(self, user_id: str, notification_type: str, actor_id: str, post_id: Optional[str] = None, message: str = "") -> dict:
        if not self.repo:
            logger.error("NotificationRepository not initialized")
            return {}
        notification = await self.repo.create(
            user_id=user_id,
            actor_id=actor_id,
            type=notification_type,
            message=message,
            post_id=post_id
        )
        logger.info(f"Notification created: {notification_type} for user {user_id}")
        return notification

    async def get_notifications(self, user_id: str, page: int = 1, limit: int = 20) -> List[dict]:
        if not self.repo:
            return []
        
        return await self.repo.get_by_user(user_id=user_id, page=page, limit=limit)

    async def mark_as_read(self, notification_id: str) -> bool:
        if not self.repo:
            return False
            
        result = await self.repo.mark_as_read(notification_id=notification_id)
        if result:
            logger.info(f"Marked notification {notification_id} as read")
        return result
