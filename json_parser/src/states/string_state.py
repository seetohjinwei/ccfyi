from .state import State, StateTransitionResult
from .common import AttemptResult, AttemptResultAction, attempt_states, get_failed_result
from ..str_slice import Str


control_characters: list[str] = []  # ignore control characters
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
            # not escape character
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
        if len(txt) == 0 or txt.at(0) != "\\":
            # not escape character
            return get_failed_result(txt)

        def attempt(result: StateTransitionResult) -> AttemptResult:
            return AttemptResult(AttemptResultAction.RETURN, result)

        # consume the escape character
        txt = txt.substring(1)

        result = attempt_states(txt, [
            (_StringEscapeUnicodeState, attempt),
            (_StringEscapeRegularState, attempt),
        ])

        match result:
            case AttemptResult(action=AttemptResultAction.RETURN):
                return result.get()
        return get_failed_result(txt)


class _StringControlState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        # not implemented
        return get_failed_result(txt)


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
                new_txt = txt.substring(1)
                json_struct = "".join(str_builder)

                result = StateTransitionResult(
                    is_success=True,
                    new_txt=new_txt,
                    json_struct=json_struct,
                )
                return AttemptResult(AttemptResultAction.RETURN, result)

            def attempt(result: StateTransitionResult) -> AttemptResult:
                assert isinstance(result.json_struct, str)
                str_builder.append(result.json_struct)
                return AttemptResult(AttemptResultAction.PASS, result)

            result = attempt_states(txt, [
                (_StringCloseQuotesState, attempt_close_quotes),
                (_StringEscapeState, attempt),
                (_StringConsumeState, attempt),
                (_StringControlState, attempt),
            ])
            match result:
                case AttemptResult(action=AttemptResultAction.RETURN):
                    return result.get()
                case AttemptResult(action=AttemptResultAction.PASS):
                    pass
            txt = result.new_txt

        return failed_result
