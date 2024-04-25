from ..str_slice import Str
from .common import (
    AttemptResult,
    AttemptResultAction,
    attempt_states,
    get_failed_result,
)
from .state import State, StateTransitionResult


escaped_characters: dict[str, str] = {
    '"': '"',  # quotation mark
    "\\": "\\",  # reverse solidus
    "/": "/",  # solidus
    "b": "\b",  # backspace
    "f": "\f",  # formfeed
    "n": "\n",  # linefeed
    "r": "\r",  # carriage return
    "t": "\t",  # horizontal tab
    # \u hex digits is handled separately
}
# control_characters: list[str] = [chr(x) for x in range(0, 0x1F + 1)]
control_characters: list[str] = [
    "\\x00",
    "\\x01",
    "\\x02",
    "\\x03",
    "\\x04",
    "\\x05",
    "\\x06",
    "\\x07",
    "\\x08",
    "\t",
    "\n",
    "\\x0b",
    "\\x0c",
    "\r",
    "\\x0e",
    "\\x0f",
    "\\x10",
    "\\x11",
    "\\x12",
    "\\x13",
    "\\x14",
    "\\x15",
    "\\x16",
    "\\x17",
    "\\x18",
    "\\x19",
    "\\x1a",
    "\\x1b",
    "\\x1c",
    "\\x1d",
    "\\x1e",
    "\\x1f",
]


class _StringOpenQuotesState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != '"':
            # not string
            return get_failed_result(txt)

        # consume the open quotes
        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )


class _StringCloseQuotesState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != '"':
            # not string
            return get_failed_result(txt)

        # consume the close quotes
        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )


class _StringConsumeState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0:
            return get_failed_result(txt)

        json_struct = txt.at(0)
        new_txt = txt.substring(1)

        return StateTransitionResult(True, new_txt, json_struct)


class _StringEscapeUnicodeState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) != "u":
            # not u
            return get_failed_result(txt)

        txt = txt.substring(1)

        # special escape character (hex digits for unicode)
        digits = txt.substring(0, 4)
        try:
            json_struct = chr(int(str(digits), 16))
        except ValueError:
            return get_failed_result(txt)

        new_txt = txt.substring(4)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=json_struct,
        )


class _StringEscapeRegularState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        if len(txt) == 0 or txt.at(0) not in escaped_characters:
            # not an escaped character
            return get_failed_result(txt)

        json_struct = escaped_characters[txt.at(0)]

        new_txt = txt.substring(1)
        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=json_struct,
        )


class _StringEscapeState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        if len(txt) == 0 or txt.at(0) != "\\":
            # not escape character (OK!)
            return failed_result

        def attempt(result: StateTransitionResult) -> AttemptResult:
            return AttemptResult(AttemptResultAction.RETURN, result)

        # consume the escape character
        txt = txt.substring(1)

        result = attempt_states(
            txt,
            [
                (_StringEscapeUnicodeState, attempt),
                (_StringEscapeRegularState, attempt),
            ],
        )

        match result:
            case AttemptResult(action=AttemptResultAction.RETURN):
                return result.get()
            case AttemptResult(action=AttemptResultAction.NO_MATCH):
                # should fail the entire string
                from src.json_struct import InvalidJSONStruct
                raise InvalidJSONStruct

        return failed_result


class _StringControlState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        cc = txt.starts_with_any(control_characters)
        if cc is None:
            return get_failed_result(txt)

        json_struct = str(cc)
        length = len(cc)
        txt = txt.substring(length)

        return StateTransitionResult(True, txt, json_struct)


class StringState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        result = _StringOpenQuotesState.transition(txt)
        if not result.is_success:
            return failed_result
        txt = result.new_txt

        str_builder: list[str] = []

        while len(txt) > 0:

            def attempt_close_quotes(result: StateTransitionResult) -> AttemptResult:
                json_struct = "".join(str_builder)

                result = StateTransitionResult(
                    is_success=True,
                    new_txt=result.new_txt,
                    json_struct=json_struct,
                )
                return AttemptResult(AttemptResultAction.RETURN, result)

            def attempt(result: StateTransitionResult) -> AttemptResult:
                assert isinstance(result.json_struct, str)
                str_builder.append(result.json_struct)
                return AttemptResult(AttemptResultAction.PASS, result)

            def attempt_control(result: StateTransitionResult) -> AttemptResult:
                return AttemptResult(AttemptResultAction.FAIL, result)

            result = attempt_states(
                txt,
                [
                    (_StringCloseQuotesState, attempt_close_quotes),
                    (_StringEscapeState, attempt),
                    (_StringControlState, attempt_control),
                    (_StringConsumeState, attempt),
                ],
            )
            match result:
                case AttemptResult(action=AttemptResultAction.RETURN):
                    return result.get()
                case AttemptResult(action=AttemptResultAction.NO_MATCH):
                    return failed_result
                case AttemptResult(action=AttemptResultAction.FAIL):
                    return failed_result
                case AttemptResult(action=AttemptResultAction.PASS):
                    pass
            txt = result.new_txt

        return failed_result
