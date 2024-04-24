from ..str_slice import Str
from .common import get_failed_result
from .state import State, StateTransitionResult


def _loop_number(txt: Str) -> tuple[Str, int]:
    value = 0
    while len(txt) > 0:
        c = txt.at(0)
        if c in [str(x) for x in range(0, 9 + 1)]:
            txt = txt.substring(1)
            value *= 10
            value += ord(c) - ord("0")
        else:
            break

    return txt, value


def _loop_fraction(txt: Str) -> tuple[Str, float]:
    divisor = 10
    value = 0
    while len(txt) > 0:
        c = txt.at(0)
        if c in [str(x) for x in range(0, 9 + 1)]:
            txt = txt.substring(1)
            value += (ord(c) - ord("0")) / divisor
            divisor *= 10
        else:
            break

    return txt, value


class _NumberFractionState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        if not txt.startswith("."):
            return failed_result

        txt = txt.substring(1)

        txt, value = _loop_fraction(txt)

        return StateTransitionResult(True, txt, value)


class _NumberExponentState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        if len(txt) == 0 or txt.at(0) not in ["E", "e"]:
            return failed_result

        txt = txt.substring(1)

        is_negative: bool = False

        if txt.at(0) == "-":
            is_negative = True
            txt = txt.substring(1)
        elif txt.at(0) == "+":
            txt = txt.substring(1)

        txt, value = _loop_number(txt)

        if is_negative:
            value *= -1

        return StateTransitionResult(True, txt, value)


class NumberState(State):
    @staticmethod
    def transition(txt: Str) -> StateTransitionResult:
        failed_result = get_failed_result(txt)

        is_negative: bool = False

        if txt.startswith("-"):
            is_negative = True
            txt = txt.substring(1)

        if txt.startswith("0"):
            value = 0
            txt = txt.substring(1)
        else:
            if len(txt) == 0 or txt.at(0) not in [str(x) for x in range(0, 9 + 1)]:
                return failed_result

            txt, value = _loop_number(txt)

        result = _NumberFractionState.transition(txt)
        if result.is_success:
            assert isinstance(result.json_struct, float)
            fractional = result.json_struct
            value += fractional
            txt = result.new_txt

        if is_negative:
            value *= -1

        exponent: int | None = None
        result = _NumberExponentState.transition(txt)
        if result.is_success:
            assert isinstance(result.json_struct, int) and not isinstance(
                result.json_struct, bool
            )
            exponent = result.json_struct
            while exponent > 0:
                exponent -= 1
                value *= 10
            while exponent < 0:
                exponent += 1
                value /= 10
            txt = result.new_txt

        return StateTransitionResult(True, txt, value)
