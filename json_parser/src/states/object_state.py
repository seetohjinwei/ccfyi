from ..json_struct import JSONStruct
from ..str_slice import Str
from .common import (
    AttemptResult,
    AttemptResultAction,
    attempt_states,
    get_failed_result,
)
from .state import State, StateTransitionResult
from .string_state import StringState
from .whitespace_state import WhitespaceState

import src.states.value_state


class _ObjectOpenState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != "{":
            # not object
            return get_failed_result(txt)

        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )


class _ObjectCloseState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != "}":
            # not object
            return get_failed_result(txt)

        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )


class _ObjectCommaState(State):
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


class _ObjectPairState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        result = WhitespaceState.transition(txt)
        txt = result.new_txt

        result = StringState.transition(txt)
        if not result.is_success:
            return failed_result
        txt = result.new_txt
        key = result.json_struct
        assert isinstance(key, str)

        result = WhitespaceState.transition(txt)
        txt = result.new_txt

        # consume colon
        if len(txt) == 0 or txt.at(0) != ":":
            return failed_result
        txt = txt.substring(1)

        result = src.states.value_state.ValueState.transition(txt)
        txt = result.new_txt
        value = result.json_struct

        return StateTransitionResult(True, txt, {key: value})


class ObjectState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        result = _ObjectOpenState.transition(txt)
        if not result.is_success:
            return failed_result
        txt = result.new_txt

        is_first_element: bool = True
        is_prev_comma: bool = False

        json_struct: dict[str, JSONStruct] = {}

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
                    (_ObjectPairState, attempt_value),
                    (_ObjectCloseState, attempt_close),
                    (_ObjectCommaState, attempt_comma),
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
                    d = result.get().json_struct
                    assert isinstance(d, dict)
                    json_struct.update(d)
                case AttemptResult(action=AttemptResultAction.NO_MATCH):
                    return failed_result
                case AttemptResult(action=AttemptResultAction.PASS):
                    pass

        return failed_result
