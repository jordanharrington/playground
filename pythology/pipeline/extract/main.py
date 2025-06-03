import sys
import json
from common.api import ObjectStore, SourceType, JobConfig, Stage
from common.utils import get_logger, load_yaml, trace
from pipeline.extract.extractor import BaseExtractor
from pipeline.extract.reddit import RedditPostExtractor

logger = get_logger(__name__)


def get_extractor(source_type: SourceType) -> BaseExtractor:
    if source_type == SourceType.REDDIT:
        return RedditPostExtractor()
    raise ValueError(f"Unsupported source type: {source_type}")


def main():
    try:
        with trace("Extract Stage"):
            config = load_yaml('job_config.yaml')
            job_config = JobConfig(**config)

            extractor_type = job_config.source.type
            extractor = get_extractor(extractor_type)

            stage_config = job_config.stages[Stage.EXTRACT]
            with trace(f"Extraction with config {stage_config}"):
                result = extractor.extract(stage_config.spec)

            object_store = ObjectStore(job_config.object_storage)
            object_store.put(key=stage_config.output, data=json.dumps(result).encode('utf-8'))
    except Exception as e:
        logger.exception(f"Fatal error during extraction: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()
