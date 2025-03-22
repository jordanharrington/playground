from abc import abstractmethod, ABC
from typing import List

from common.api.v1 import ExtractionConfig


class BaseExtractor(ABC):
    @abstractmethod
    def extract(self, config: ExtractionConfig) -> List[str]:
        pass
