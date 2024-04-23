from .state import State, StateTransitionResult
from ..str_slice import Str


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
