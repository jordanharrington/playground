from abc import abstractmethod, ABC
from typing import List, Dict


class BaseExtractor(ABC):
    @abstractmethod
    def extract(self, config: Dict[str, any]) -> List[str]:
        pass
