from .state import State, StateTransitionResult
from ..str_slice import Str


class InitialState(State):
    pass
    # @staticmethod
    # def transition(txt: Str) -> StateTransitionResult:
    #     states: list[type[State]] = [
    #         ObjectState,
    #         ArrayState,
    #         ValueState,
    #         StringState,
    #         NumberState,
    #         WhitespaceState,
    #     ]
    #
    #     for state in states:
    #         res = state.transition(txt)
    #         if res.is_success:
    #             return res
    #
    #     # does not match any of the states
    #     return StateTransitionResult(
    #         is_success=False,
    #         new_txt=txt,
    #         json_struct=None,
    #     )
