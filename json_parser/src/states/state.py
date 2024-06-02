from src.json_struct.json_struct import JSONStruct
from ..str_slice import Str
from abc import ABC, abstractmethod
from dataclasses import dataclass


@dataclass
class StateTransitionResult:
    is_success: bool
    new_txt: Str
    json_struct: JSONStruct


class State(ABC):
    @staticmethod
    @abstractmethod
    def transition(txt: Str) -> StateTransitionResult:
        """
        Starting from current State, consume `txt` to transition to a next state.
        If no next state is possible, None is returned.
        """

        pass
