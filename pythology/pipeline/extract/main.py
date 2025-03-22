import sys

from common.api.v1 import ExtractionConfig, KafkaProducer, SourceType
from common.utils import get_logger, load_yaml, trace
from pipeline.extract import BaseExtractor
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
            extraction_config = ExtractionConfig(**config)
            logger.info(f"Extracting data with config {extraction_config}")

            extractor_type = extraction_config.source.type
            extractor = get_extractor(extractor_type)

            with trace("Extraction"):
                msgs = extractor.extract(extraction_config)

            topic_config = extraction_config.kafka
            topic = topic_config.topics[extractor_type]
            bootstrap_servers = topic_config.bootstrap_servers
            producer = KafkaProducer(bootstrap_servers=bootstrap_servers)

            with trace(f"Producing {len(msgs)} messages to Kafka topic '{topic}'"):
                for msg in msgs:
                    producer.send(topic, msg)
                producer.flush()
    except Exception as e:
        logger.exception(f"Fatal error during extraction: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()
