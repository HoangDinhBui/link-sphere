import os

import uvicorn
from fastapi import FastAPI
from loguru import logger

from app.routers import search_router

app = FastAPI(
    title="LinkSphere Search Service",
    description="Full-text search for posts using OpenSearch",
    version="1.0.0",
)

# Include routers
app.include_router(search_router.router, prefix="/api/v1/search")


@app.on_event("startup")
async def startup_event():
    logger.info("Search Service starting up...")


@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Search Service shutting down...")


@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "search-service"}


if __name__ == "__main__":
    port = int(os.getenv("SERVER_PORT", "8007"))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
