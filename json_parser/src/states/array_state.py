from ..str_slice import Str
from .common import (
    AttemptResult,
    AttemptResultAction,
    attempt_states,
    get_failed_result,
)
from .state import State, StateTransitionResult
from .value_state import ValueState
from .whitespace_state import WhitespaceState


class _ArrayOpenState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != "[":
            # not array
            return get_failed_result(txt)

        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )


class _ArrayCloseState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != "]":
            # not array
            return get_failed_result(txt)

        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )


class _ArrayCommaState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != ",":
            return get_failed_result(txt)

        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=",",
        )


class ArrayState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        result = _ArrayOpenState.transition(txt)
        if not result.is_success:
            return failed_result
        txt = result.new_txt

        is_first_element: bool = True
        is_prev_comma: bool = False

        json_struct: list = []

        while len(txt) > 0:

            def attempt_whitespace(result: StateTransitionResult) -> AttemptResult:
                return AttemptResult(AttemptResultAction.PASS, result)

            def attempt_value(result: StateTransitionResult) -> AttemptResult:
                return AttemptResult(AttemptResultAction.VALUE, result)

            def attempt_close(result: StateTransitionResult) -> AttemptResult:
                if is_prev_comma:
                    return AttemptResult(AttemptResultAction.FAIL, result)

                result = StateTransitionResult(
                    is_success=True,
                    new_txt=result.new_txt,
                    json_struct=json_struct,
                )
                return AttemptResult(AttemptResultAction.RETURN, result)

            def attempt_comma(result: StateTransitionResult) -> AttemptResult:
                if is_first_element or is_prev_comma:
                    return AttemptResult(AttemptResultAction.FAIL, result)
                return AttemptResult(AttemptResultAction.VALUE, result)

            result = WhitespaceState.transition(txt)

            result = attempt_states(
                txt,
                [
                    (WhitespaceState, attempt_whitespace),
                    (ValueState, attempt_value),
                    (_ArrayCloseState, attempt_close),
                    (_ArrayCommaState, attempt_comma),
                ],
            )

            is_prev_comma = False
            txt = result.new_txt

            match result:
                case AttemptResult(action=AttemptResultAction.RETURN):
                    return result.get()
                case AttemptResult(action=AttemptResultAction.FAIL):
                    return failed_result
                case AttemptResult(action=AttemptResultAction.VALUE):
                    if result.get().json_struct == ",":
                        # handle comma separately
                        is_prev_comma = True
                        continue

                    is_first_element = False
                    json_struct.append(result.get().json_struct)
                case AttemptResult(action=AttemptResultAction.NO_MATCH):
                    return failed_result
                case AttemptResult(action=AttemptResultAction.PASS):
                    pass

        return failed_result
