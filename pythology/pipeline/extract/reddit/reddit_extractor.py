from abc import ABC

import praw
from praw.models import Submission

from common.api.v1 import ExtractionConfig, RedditPost
from common.utils import load_ini, get_logger
from pipeline.extract import BaseExtractor

# Setup pyraw Reddit client
config = load_ini('reddit_config.ini')
reddit = praw.Reddit(
    client_id=config['reddit']['client_id'],
    client_secret=config['reddit']['client_secret'],
    user_agent=config['reddit']['user_agent'],
)
# Setup logging
logger = get_logger(__name__)


def from_raw_post(post: Submission) -> RedditPost:
    if not all([
        post.id,
        post.title,
        post.selftext,
        post.author and post.author.name,
        post.subreddit and post.subreddit.display_name,
        post.num_comments is not None,
        post.created_utc is not None,
        post.over_18 is not None,
        post.upvote_ratio is not None,
        post.score is not None,
    ]):
        raise ValueError(f"Incomplete post: {post}")

    return RedditPost(
        id=post.id,
        author=post.author.name,
        title=post.title,
        content=post.selftext,
        num_comments=post.num_comments,
        subreddit=post.subreddit.display_name,
        created_utc=post.created_utc,
        over_18=post.over_18,
        upvote_ratio=post.upvote_ratio,
        score=post.score
    )


def extract_posts(query: str, subreddit: str, limit: int, sort: str) -> list[str]:
    logger.info(f"Extracting {limit} posts from {subreddit} with {query} sorted by {sort}")
    _subreddit = reddit.subreddit(subreddit)
    results = _subreddit.search(query=query, sort=sort, limit=limit)
    num_results = 0
    posts = []
    for post in results:
        try:
            num_results += 1
            reddit_post = from_raw_post(post)
            posts.append(reddit_post.model_dump_json())
        except ValueError as e:
            logger.warning(f"Skipping invalid post: {e}")
            continue
    logger.info(f"Finished extracting {len(posts)} posts from {num_results} responses and subreddit [{subreddit}]")
    return posts


class RedditPostExtractor(BaseExtractor, ABC):
    def __init__(self):
        super().__init__()

    def extract(self, extraction_config: ExtractionConfig) -> list[str]:
        job_config = extraction_config.reddit
        if job_config is None:
            raise ValueError("Reddit config is empty")

        reddit_posts = []
        for subreddit in job_config.subreddits:
            request_limit = subreddit.limit
            extracted = extract_posts(query=subreddit.query, subreddit=subreddit.name, limit=request_limit,
                                      sort=subreddit.sort)
            num_extracted = len(extracted)
            if num_extracted < request_limit:
                logger.warning(f"Requested {request_limit} posts but only got {num_extracted} from '{subreddit}''")
            reddit_posts.extend(extracted)

        return reddit_posts
