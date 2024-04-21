from .states import State
from dataclasses import dataclass


@dataclass
class Transition:
    next_state: State | None
    can_transition: bool
