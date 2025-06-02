import configparser
import os
from typing import Any, Optional

import yaml

CONF_DIR = os.getenv("CONF_DIR", "/conf")


def load_yaml(config_path: str) -> Optional[dict[str, Any]]:
    abs_path = get_abs_conf_path(config_path)
    if not os.path.exists(abs_path):
        raise FileNotFoundError(f"YAML config not found: {config_path} at {abs_path}")

    with open(abs_path, 'r') as f:
        config = yaml.safe_load(f)

    return config


def load_ini(config_path: str) -> configparser.ConfigParser:
    config = configparser.ConfigParser()
    read = config.read(get_abs_conf_path(config_path))
    if len(read) == 0:
        raise FileNotFoundError(f"INI config not found: {config_path}")
    return config


def get_abs_conf_path(path: str) -> str:
    global CONF_DIR
    return os.path.join(CONF_DIR, path)
