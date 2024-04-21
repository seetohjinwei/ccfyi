from .states import InitialState, State
from .str_slice import Str


class StateMachine:
    txt: Str
    state: State

    def __init__(self, txt: str):
        self.txt = Str(txt)
        self.state = InitialState()

    def transition(self):
        pass
