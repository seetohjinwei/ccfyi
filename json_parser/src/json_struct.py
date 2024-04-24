from .str_slice import Str


# https://www.rfc-editor.org/rfc/rfc8259.html#section-3
JSONStruct = bool | None | dict[str, "JSONStruct"] | list["JSONStruct"] | float | str


class InvalidJSONStruct(Exception):
    pass


def parse(txt: str) -> JSONStruct:
    from src.states.value_state import ValueState

    result = ValueState.transition(Str(txt))
    if not result.is_success:
        raise InvalidJSONStruct

    return result.json_struct
