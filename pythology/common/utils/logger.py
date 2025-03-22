import logging.config

from .parser import load_yaml

_logger_initialized = False


def get_logger(name: str = __name__, config_path: str = "logging.yaml") -> logging.Logger:
    global _logger_initialized
    if not _logger_initialized:
        config = load_yaml(config_path)
        logging.config.dictConfig(config)
        _logger_initialized = True
    return logging.getLogger(name)
