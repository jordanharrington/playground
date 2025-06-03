from typing import List

from pydantic import BaseModel, field_validator


class RedditPost(BaseModel):
    id: str
    author: str
    title: str
    content: str
    num_comments: int
    subreddit: str
    created_utc: float
    over_18: bool
    upvote_ratio: float
    score: int


class SubredditQuery(BaseModel):
    name: str
    query: str
    limit: int
    sort: str
    time_filter: str

    @classmethod
    @field_validator('limit')
    def check_limit(cls, v):
        if v <= 0:
            raise ValueError("limit must be greater than 0")
        return v

    @classmethod
    @field_validator('query', 'name', 'sort', 'time_filter')
    def check_non_empty_strings(cls, v, info):
        if not v or not isinstance(v, str):
            raise ValueError(f"{info.field_name} must be a non-empty string")
        return v


class RedditExtractorConfig(BaseModel):
    subreddits: List[SubredditQuery]

    @classmethod
    @field_validator('subreddits')
    def check_subreddits(cls, v):
        if not v:
            raise ValueError("subreddits cannot be empty")
        if not all(isinstance(s, str) for s in v):
            raise TypeError("each subreddit must be a string")
        return v
