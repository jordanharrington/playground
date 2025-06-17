import sys
import json
from common.api import ObjectStore, JobConfig, Stage
from common.utils import get_logger, load_yaml, trace

logger = get_logger(__name__)



def main():
    try:
        with trace("Load Stage"):
            config = load_yaml('job_config.yaml')
            job_config = JobConfig(**config)
            stage_config = job_config.stages[Stage.TRANSFORM]
            object_store = ObjectStore(job_config.object_storage)


    except Exception as e:
        logger.exception(f"Fatal error during extraction: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()