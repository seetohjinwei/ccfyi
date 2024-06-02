from ..str_slice import Str


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


INDENTATION = "  "


def _pretty_print_array(array: list["JSONStruct"], indent_level: int) -> str:
    length = len(array)

    if length == 0:
        return "[]"

    str_builder: list[str] = []

    for index, struct in enumerate(array):
        indentation = INDENTATION * (indent_level + 1)

        item: str = indentation + _pretty_print(struct, indent_level + 1)
        if index != length - 1:
            item += ","

        str_builder.append(item)

    indentation = INDENTATION * indent_level
    return f"""[
{'\n'.join(str_builder)}
{indentation}]"""


def _pretty_print_object(object: dict[str, "JSONStruct"], indent_level: int) -> str:
    length = len(object)

    if length == 0:
        return "{}"

    str_builder: list[str] = []

    index = 0
    for key, struct in object.items():
        item: str = _pretty_print(struct, indent_level + 1)

        if index != length - 1:
            item += ","
        index += 1

        indentation = INDENTATION * (indent_level + 1)
        str_builder.append(f'{indentation}"{key}": {item}')

    indentation = INDENTATION * indent_level
    return f"""{{
{'\n'.join(str_builder)}
{indentation}}}"""


def _pretty_print(struct: JSONStruct, indent_level: int) -> str:
    if isinstance(struct, bool):
        if struct is True:
            return "true"
        elif struct is False:
            return "false"
    elif struct is None:
        return "null"
    elif isinstance(struct, dict):
        return _pretty_print_object(struct, indent_level)
    elif isinstance(struct, list):
        return _pretty_print_array(struct, indent_level)
    elif isinstance(struct, float) or isinstance(struct, int):
        return str(struct)
    elif isinstance(struct, str):
        repr_ = repr(struct)
        if repr_[0] == "'":
            # swap '' -> ""
            return f'"{repr_[1:-1]}"'
        return repr_


def pretty_print(struct: JSONStruct) -> str:
    return _pretty_print(struct, 0)
