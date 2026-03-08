import asyncio
import os

import uvicorn
from fastapi import FastAPI
from loguru import logger

from app.routers import notification_router
from app.kafka_consumer import start_kafka_consumer

app = FastAPI(
    title="LinkSphere Notification Service",
    description="Handles realtime notifications via WebSocket and Kafka events",
    version="1.0.0",
)

# Include routers
app.include_router(notification_router.router, prefix="/api/v1/notifications")


@app.on_event("startup")
async def startup_event():
    """Start Kafka consumer on application startup."""
    logger.info("Notification Service starting up...")
    asyncio.create_task(start_kafka_consumer())


@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Notification Service shutting down...")


@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "notification-service"}


if __name__ == "__main__":
    port = int(os.getenv("SERVER_PORT", "8006"))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
