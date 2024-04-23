from .states.state import State
from dataclasses import dataclass


@dataclass
class Transition:
    next_state: State | None
    can_transition: bool
