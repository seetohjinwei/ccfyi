from .json_struct import JSONStruct
from .str_slice import Str
from abc import ABC, abstractmethod
from dataclasses import dataclass


@dataclass
class StateTransitionResult:
    next_state: "State | None"
    txt: Str
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
            if res.next_state:
                return res

        # does not match any of the states
        return StateTransitionResult(
            next_state=None,
            txt=txt,
            json_struct=None,
        )


class ObjectState(State):
    pass


class ArrayState(State):
    pass


class ValueState(State):
    pass


class StringState(State):
    pass


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

        for i, c in enumerate(txt):
            if c in whitespace_chars:
                continue

            # lo is the first non-whitespace character
            lo = i
            txt = txt.substring(lo)

            return StateTransitionResult(
                next_state=None,
                txt=txt,
                json_struct=None,
            )

        return StateTransitionResult(
            next_state=None,
            txt=txt,
            json_struct=None,
        )
