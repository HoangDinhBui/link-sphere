import asyncio
import asyncpg
import os

import uvicorn
from fastapi import FastAPI
from loguru import logger

from app.routers import notification_router
from app.kafka_consumer import start_kafka_consumer

from app.repositories.notification_repository import NotificationRepository
from app.routers.notification_router import notification_service as router_service
from app.kafka_consumer import notification_service as consumer_service

app = FastAPI(
    title="LinkSphere Notification Service",
    description="Handles realtime notifications via WebSocket and Kafka events",
    version="1.0.0",
)

# Include routers
app.include_router(notification_router.router, prefix="/api/v1/notifications")

db_pool = None

@app.on_event("startup")
async def startup_event():
    global db_pool
    db_url = os.getenv("DATABASE_URL", "postgresql://linksphere:linksphere@localhost:5432/linksphere")
    
    try:
        db_pool = await asyncpg.create_pool(db_url)
        logger.info("Connected to PostgreSQL successfully")

        repo = NotificationRepository(db_pool)
        router_service.repo = repo
        consumer_service.repo = repo
    except Exception as e:
        logger.error(f"Failed to connect to PostgreSQL: {e}")

    asyncio.create_task(start_kafka_consumer())


@app.on_event("shutdown")
async def shutdown_event():
    global db_pool
    if db_pool:
        await db_pool.close()
        logger.info("Closed PostgreSQL connection pool")


@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "notification-service"}


if __name__ == "__main__":
    port = int(os.getenv("SERVER_PORT", "8006"))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
