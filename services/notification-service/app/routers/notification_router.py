import os
from typing import Dict, Set, Optional

from fastapi import APIRouter, WebSocket, WebSocketDisconnect, Depends, HTTPException, Header
from jose import jwt, JWTError
from loguru import logger

from app.models.notification import NotificationListRequest, NotificationResponse, MarkAsReadRequest
from app.services.notification_service import NotificationService

router = APIRouter()
notification_service = NotificationService()

# Store active WebSocket connections: user_id -> set of WebSocket connections
active_connections: Dict[str, Set[WebSocket]] = {}

def get_current_user(authorization: Optional[str] = Header(None)) -> str:
    if not authorization or not authorization.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="Missing or invalid token")
    
    token = authorization.split(" ")[1]
    secret = os.getenv("JWT_SECRET", "super-secret-jwt-key")
    
    try:
        payload = jwt.decode(token, secret, algorithms=["HS256"])
        user_id = payload.get("user_id") or payload.get("sub")
        if not user_id:
            raise HTTPException(status_code=401, detail="Invalid token payload")
        return user_id
    except JWTError as e:
        logger.error(f"JWT decode error: {e}")
        raise HTTPException(status_code=401, detail="Could not validate credentials")


@router.post("/list", response_model=NotificationResponse)
async def list_notifications(
    request: NotificationListRequest,
    user_id: str = Depends(get_current_user) # Inject JWT middleware into this function
):
    """List notifications for the authenticated user."""
    notifications = await notification_service.get_notifications(
        user_id=user_id,
        page=request.page,
        limit=request.limit,
    )
    return {
        "code_status": 200,
        "message": "notifications retrieved successfully",
        "result": True,
        "errors": {},
        "data": notifications,
    }


@router.post("/read", response_model=NotificationResponse)
async def mark_as_read(
    request: MarkAsReadRequest,
    user_id: str = Depends(get_current_user)
):
    """Mark a notification as read for the authenticated user."""
    success = await notification_service.mark_as_read(request.notification_id)
    
    if not success:
        raise HTTPException(status_code=404, detail="Notification not found")

    return {
        "code_status": 200,
        "message": "notification marked as read",
        "result": True,
        "errors": {},
        "data": None,
    }


@router.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket):
    """WebSocket endpoint for realtime notifications."""
    await websocket.accept()

    # TODO: Extract user_id from query param or JWT
    user_id = websocket.query_params.get("user_id", "anonymous")

    # Register connection
    if user_id not in active_connections:
        active_connections[user_id] = set()
    active_connections[user_id].add(websocket)

    logger.info(f"WebSocket connected: user_id={user_id}")

    try:
        while True:
            # Keep connection alive, listen for client messages
            data = await websocket.receive_text()
            logger.debug(f"Received from {user_id}: {data}")
    except WebSocketDisconnect:
        active_connections[user_id].discard(websocket)
        if not active_connections[user_id]:
            del active_connections[user_id]
        logger.info(f"WebSocket disconnected: user_id={user_id}")

# If user open 5 tabs at the same time, all of them will display the notification
async def broadcast_to_user(user_id: str, message: dict):
    """Send a notification to all active WebSocket connections of a user."""
    if user_id in active_connections:
        disconnected = set()
        for ws in active_connections[user_id]:
            try:
                await ws.send_json(message)
            except Exception:
                disconnected.add(ws)
        active_connections[user_id] -= disconnected
