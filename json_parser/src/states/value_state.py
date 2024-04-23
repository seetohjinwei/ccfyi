from ..str_slice import Str
from .common import (
    AttemptResult,
    AttemptResultAction,
    attempt_states,
    get_failed_result,
)
from .state import State, StateTransitionResult
from .whitespace_state import WhitespaceState
from .string_state import StringState
from .number_state import NumberState

# let value_state.py be the one to break the circular imports
import src.states.array_state
import src.states.object_state


class _ValueOtherState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        # TODO: use txt.startswith

        t = str(txt)
        if t.startswith("true"):
            return StateTransitionResult(True, txt.substring(4), True)
        elif t.startswith("false"):
            return StateTransitionResult(True, txt.substring(5), False)
        elif t.startswith("null"):
            return StateTransitionResult(True, txt.substring(4), None)

        return get_failed_result(txt)


class ValueState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        def attempt(result: StateTransitionResult) -> AttemptResult:
            return AttemptResult(AttemptResultAction.VALUE, result)

        result = WhitespaceState.transition(txt)
        txt = result.new_txt
        json_struct = None

        result = attempt_states(
            txt,
            [
                (StringState, attempt),
                (NumberState, attempt),
                (src.states.object_state.ObjectState, attempt),
                (src.states.array_state.ArrayState, attempt),
                (_ValueOtherState, attempt),
            ],
        )
        match result:
            case AttemptResult(action=AttemptResultAction.VALUE):
                txt = result.new_txt
                json_struct = result.get().json_struct
            case AttemptResult(action=AttemptResultAction.NO_MATCH):
                return failed_result

        result = WhitespaceState.transition(txt)
        txt = result.new_txt

        return StateTransitionResult(True, txt, json_struct)
