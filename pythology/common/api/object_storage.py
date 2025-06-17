from typing import Optional

import boto3
from botocore.exceptions import ClientError
from pydantic import BaseModel, field_validator

from common.utils import get_logger

logger = get_logger(__name__)


class ObjectStorageCredentials(BaseModel):
    profile_name: Optional[str] = None
    access_key: Optional[str] = None
    secret_key: Optional[str] = None
    session_token: Optional[str] = None


class ObjectStorageConfig(BaseModel):
    bucket: str
    region: str
    credentials: ObjectStorageCredentials
    endpoint_url: Optional[str] = None

    @classmethod
    @field_validator('bucket', "region", 'credentials')
    def check_non_empty(cls, v, info):
        if not v or not isinstance(v, str):
            raise ValueError(f"{info.field_name} must be a non-empty string")
        return v


class ObjectStore:
    def __init__(self, config: ObjectStorageConfig):
        session = boto3.Session(region_name=config.region,
                                aws_access_key_id=config.credentials.access_key,
                                aws_secret_access_key=config.credentials.secret_key,
                                aws_session_token=config.credentials.session_token,
                                profile_name=config.credentials.profile_name)
        self.s3 = session.client('s3', endpoint_url=config.endpoint_url)
        self.bucket = config.bucket

    def put(self, key, data):
        """Upload data (bytes) to the bucket with the given key."""
        try:
            self.s3.put_object(Bucket=self.bucket, Key=key, Body=data)
            logger.info(f"Uploaded {key} to {self.bucket}")
        except ClientError as e:
            logger.error(f"Error Uploading {key}: {e}")
            raise ValueError(e)

    def get(self, key):
        """Download data (bytes) from the bucket with the given key."""
        try:
            response = self.s3.get_object(Bucket=self.bucket, Key=key)
            return response['Body'].read()
        except ClientError as e:
            logger.error(f"Error fetching {key}: {e}")
            raise ValueError(e)
