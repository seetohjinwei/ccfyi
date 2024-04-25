from .str_slice import Str


# https://www.rfc-editor.org/rfc/rfc8259.html#section-3
JSONStruct = bool | None | dict[str, "JSONStruct"] | list["JSONStruct"] | float | str


class InvalidJSONStruct(Exception):
    pass


def parse(txt: str) -> JSONStruct:
    from src.states.array_state import ArrayState
    from src.states.object_state import ObjectState
    from src.states.state import StateTransitionResult

    def _is_success(result: StateTransitionResult) -> bool:
        ret = result.is_success and len(result.new_txt) == 0
        return ret

    result = ArrayState.transition(Str(txt))
    if _is_success(result):
        return result.json_struct
    result = ObjectState.transition(Str(txt))
    if _is_success(result):
        return result.json_struct

    raise InvalidJSONStruct(result)
