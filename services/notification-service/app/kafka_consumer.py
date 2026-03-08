import json
import os

from aiokafka import AIOKafkaConsumer
from loguru import logger

from app.services.notification_service import NotificationService
from app.routers.notification_router import broadcast_to_user

notification_service = NotificationService()


async def start_kafka_consumer():
    """Start consuming events from Kafka topics."""
    kafka_brokers = os.getenv("KAFKA_BROKERS", "localhost:9092")

    consumer = AIOKafkaConsumer(
        "post-events",
        "comment-events",
        "user-events",
        bootstrap_servers=kafka_brokers,
        group_id="notification-service",
        value_deserializer=lambda m: json.loads(m.decode("utf-8")),
    )

    try:
        await consumer.start()
        logger.info(f"Kafka consumer started, listening on brokers: {kafka_brokers}")

        async for msg in consumer:
            try:
                await process_event(msg.topic, msg.value)
            except Exception as e:
                logger.error(f"Error processing event: {e}")

    except Exception as e:
        logger.error(f"Kafka consumer error: {e}")
    finally:
        await consumer.stop()


async def process_event(topic: str, event: dict):
    """Process a Kafka event and create appropriate notification."""
    event_type = event.get("event", "")
    user_id = event.get("userId", "")
    post_id = event.get("postId", "")

    logger.info(f"Processing event: {event_type} from topic: {topic}")

    if event_type == "post.liked":
        # Notify the post owner that someone liked their post
        notification = await notification_service.create_notification(
            user_id=user_id,  # Should be post owner, simplified here
            notification_type="like",
            actor_id=user_id,
            post_id=post_id,
            message="Someone liked your post",
        )
        await broadcast_to_user(user_id, notification)

    elif event_type == "post.commented":

        notification = await notification_service.create_notification(
            user_id=user_id,
            notification_type="comment",
            actor_id=user_id,
            post_id=post_id,
            message="Someone commented on your post",
        )
        await broadcast_to_user(user_id, notification)

    elif event_type == "user.followed":
        target_user_id = event.get("targetUserId", "")
        notification = await notification_service.create_notification(
            user_id=target_user_id,
            notification_type="follow",
            actor_id=user_id,
            message="Someone started following you",
        )
        await broadcast_to_user(target_user_id, notification)

    else:
        logger.warning(f"Unknown event type: {event_type}")
