import time
from contextlib import contextmanager
from .logger import get_logger

logger = get_logger(__name__)


@contextmanager
def trace(name: str):
    start = time.perf_counter()
    try:
        logger.info(f"Starting: {name}...")
        yield
        duration_ms = (time.perf_counter() - start) * 1000
        logger.info(f"Completed: {name} in {duration_ms:.2f} ms")
    except Exception as e:
        duration_ms = (time.perf_counter() - start) * 1000
        logger.exception(f"Error during '{name}' after {duration_ms:.2f} ms: {e}")
        raise
