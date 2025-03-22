from enum import Enum
from typing import Optional

from pydantic import BaseModel, field_validator

from .kafka import KafkaTopicConfig
from .reddit import RedditExtractorConfig


class SourceType(str, Enum):
    REDDIT = "reddit"


class Source(BaseModel):
    type: SourceType


class ExtractionConfig(BaseModel):
    source: Source
    kafka: KafkaTopicConfig
    reddit: Optional[RedditExtractorConfig]

    @classmethod
    @field_validator('extractor_type', 'kafka_topics')
    def check_non_empty(cls, v, info):
        if not v:
            raise ValueError(f"{info.field_name} cannot be empty")
        return v
