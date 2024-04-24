from ..str_slice import Str
from .state import State, StateTransitionResult


# TODO:
class NumberState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        return StateTransitionResult(False, Str(""), None)
