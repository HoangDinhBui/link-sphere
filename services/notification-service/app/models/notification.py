from datetime import datetime
from typing import Optional

from pydantic import BaseModel


class Notification(BaseModel):
    """Notification domain model."""

    id: Optional[str] = None
    user_id: str
    type: str  # "like", "comment", "follow"
    actor_id: str  # who triggered the notification
    post_id: Optional[str] = None
    message: str
    is_read: bool = False
    created_at: Optional[datetime] = None


class NotificationListRequest(BaseModel):
    """Request body for listing notifications."""

    page: int = 1
    limit: int = 20


class NotificationResponse(BaseModel):
    """Standard API response."""

    success: bool
    data: Optional[list] = None
    message: Optional[str] = None
    error: Optional[str] = None
