from enum import Enum
from typing import Optional, Dict, Any

from pydantic import BaseModel, field_validator

from .objectstorage import ObjectStorageConfig


class Stage(str, Enum):
    EXTRACT = "extract"


class SourceType(str, Enum):
    REDDIT = "reddit"


class Source(BaseModel):
    type: SourceType


class StageConfig(BaseModel):
    input: Optional[str] = None
    output: str
    spec: Dict[str, Any]

    @classmethod
    @field_validator('output', 'params')
    def check_non_empty(cls, v, info):
        if not v:
            raise ValueError(f"{info.field_name} cannot be empty")
        return v


class JobConfig(BaseModel):
    source: Source
    object_storage: ObjectStorageConfig
    stages: Dict[Stage, StageConfig]

    @classmethod
    @field_validator('source', 'object_storage', 'stages')
    def check_non_empty(cls, v, info):
        if not v:
            raise ValueError(f"{info.field_name} cannot be empty")
        return v
