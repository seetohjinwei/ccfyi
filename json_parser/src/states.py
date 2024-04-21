from .json_struct import JSONStruct
from .str_slice import Str
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

    # @abstractmethod
    # def consume(self, txt: Str) -> Str:
    #     """
    #     Consume the string length that is required for the state transition.
    #     """
    #     pass

    # @abstractmethod
    # def next_states(self) -> list[type["State"]]:
    #     """
    #     Returns all next states reachable from this state.
    #     """
    #     pass


class InitialState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        states: list[type[State]] = [
            ObjectState,
            ArrayState,
            ValueState,
            StringState,
            NumberState,
            WhitespaceState,
        ]

        for state in states:
            res = state.transition(txt)
            if res.is_success:
                return res

        # does not match any of the states
        return StateTransitionResult(
            is_success=False,
            new_txt=txt,
            json_struct=None,
        )


class ObjectState(State):
    pass


class ArrayState(State):
    pass


class ValueState(State):
    pass


class StringState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        control_characters: list[str] = []  # ignore control characters
        escape_characters: dict[str, str] = {
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

        failed_result = StateTransitionResult(
            is_success=False,
            new_txt=txt,
            json_struct=None,
        )

        if len(txt) == 0 or txt.at(0) != '"':
            # not string
            return failed_result

        str_builder: list[str] = []

        # do custom iterator
        it = 1
        end = len(txt)
        while it < end:
            c = txt.at(it)

            if c == '"':
                # reached the end

                new_txt = txt.substring(it + 1)
                json_struct = "".join(str_builder)

                return StateTransitionResult(
                    is_success=True,
                    new_txt=new_txt,
                    json_struct=json_struct,
                )
            elif c == "\\":
                # escape character

                # increment iterator
                it += 1
                if not (it < end):
                    # unexpected end
                    return failed_result
                c = txt.at(it)
                if c == "u":
                    # special escape character (hex digits for unicode)
                    digits = txt.substring(it, 4)
                    ch = chr(int(str(digits), 16))
                    str_builder.append(ch)
                    it += 3
                    continue
                if c not in escape_characters:
                    # invalid escape character
                    return failed_result
                ch = escape_characters[c]
                str_builder.append(ch)

            elif c in control_characters:
                pass
            else:
                str_builder.append(c)

            it += 1

        # did not find an end
        return failed_result


class NumberState(State):
    pass


class WhitespaceState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        whitespace_chars: list[str] = [
            " ",  # space
            "\n",  # linefeed
            "\r",  # carriage return
            "\t",  # horizontal tab
        ]

        if len(txt) == 0 or txt.at(0) not in whitespace_chars:
            # not whitespace
            return StateTransitionResult(
                is_success=False,
                new_txt=txt,
                json_struct=None,
            )

        lo = 0
        for i, c in enumerate(txt):
            lo = i
            if c in whitespace_chars:
                continue
            break

        # lo is the first non-whitespace character
        new_txt = txt.substring(lo)

        return StateTransitionResult(
            is_success=True,
            new_txt=new_txt,
            json_struct=None,
        )
