from ..str_slice import Str
from .state import State, StateTransitionResult
from dataclasses import dataclass
from enum import Enum, auto
from typing import Callable


def get_failed_result(txt: Str) -> StateTransitionResult:
    failed_result = StateTransitionResult(
        is_success=False,
        new_txt=txt,
        json_struct=None,
    )
    return failed_result


def get_no_result() -> StateTransitionResult:
    return get_failed_result(Str(""))


class AttemptResultAction(Enum):
    UNDEFINED = auto()
    NO_MATCH = auto()
    FAIL = auto()
    RETURN = auto()
    BREAK = auto()
    CONTINUE = auto()
    PASS = auto()


@dataclass
class AttemptResult:
    action: AttemptResultAction
    new_txt: Str
    result: StateTransitionResult  # has a result <=> action == RETURN

    def __init__(self, action: AttemptResultAction, result: StateTransitionResult):
        self.action = action
        self.new_txt = result.new_txt
        self.result = result

    def get(self) -> StateTransitionResult:
        # for the type checker
        assert self.action == AttemptResultAction.RETURN
        assert self.result is not None

        return self.result


# Would be great if this is a macro :/
def attempt_states(
    txt: Str,
    states: list[tuple[type[State], Callable[[StateTransitionResult], AttemptResult]]],
) -> AttemptResult:
    for s, f in states:
        if (state := s.transition(txt)).is_success:
            return f(state)
    return AttemptResult(AttemptResultAction.NO_MATCH, get_no_result())
