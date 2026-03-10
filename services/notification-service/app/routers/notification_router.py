from typing import Dict, Set

from fastapi import APIRouter, WebSocket, WebSocketDisconnect
from loguru import logger

from app.models.notification import NotificationListRequest
from app.services.notification_service import NotificationService

router = APIRouter()
notification_service = NotificationService()

# Store active WebSocket connections: user_id -> set of WebSocket connections
active_connections: Dict[str, Set[WebSocket]] = {}


@router.post("/list")
async def list_notifications(request: NotificationListRequest):
    """List notifications for the authenticated user."""
    # TODO: Extract user_id from JWT token
    user_id = "placeholder"

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
