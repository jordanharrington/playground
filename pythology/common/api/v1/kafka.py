import json
from typing import Any, Dict

from confluent_kafka import Producer
from pydantic import BaseModel, field_validator

from common.utils import get_logger

logger = get_logger(__name__)


def delivery_report(err, msg):
    if err is not None:
        logger.error(f"Delivery failed for {msg.topic()}: {err}")
    else:
        logger.info(f"Delivered to {msg.topic()} [{msg.partition()}] at offset {msg.offset()}")


class KafkaProducer:
    def __init__(self, bootstrap_servers: str):
        self.producer = Producer({
            "bootstrap.servers": bootstrap_servers,
        })

    def send(self, topic: str, value: Any):
        self.producer.produce(
            topic=topic,
            value=json.dumps(value).encode("utf-8"),
            callback=delivery_report
        )
        self.producer.poll(0)  # Allow background delivery

    def flush(self):
        self.producer.flush()


class KafkaTopicConfig(BaseModel):
    bootstrap_servers: str
    group_id: str
    topics: Dict[str, str]

    @classmethod
    @field_validator('bootstrap_servers', 'group_id')
    def check_non_empty(cls, v, info):
        if not v or not isinstance(v, str):
            raise ValueError(f"{info.field_name} must be a non-empty string")
        return v

    @classmethod
    @field_validator('topics')
    def check_topics_not_empty(cls, v):
        if not v or not isinstance(v, dict):
            raise ValueError("topics must be a non-empty dictionary mapping SourceType to topic names")
        return v
