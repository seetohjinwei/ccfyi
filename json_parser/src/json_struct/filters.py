from abc import ABC, abstractmethod
from dataclasses import dataclass
import re
from typing import Optional

from src.json_struct.json_struct import JSONStruct


class InvalidFilterApplication(Exception):
    pass


class InvalidFilterArgument(Exception):
    pass


@dataclass
class Flags:
    should_collect_into_array: bool = False


class Filter(ABC):
    @abstractmethod
    def apply(self, struct: JSONStruct) -> JSONStruct: ...


class IdentityFilter(Filter):
    def apply(self, struct: JSONStruct) -> JSONStruct:
        return struct

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, IdentityFilter):
            return False

        return True


class ArrayIndexFilter(Filter):
    def __init__(self, index: int):
        self.index = index

    def apply(self, struct: JSONStruct) -> JSONStruct:
        if not isinstance(struct, list):
            raise InvalidFilterApplication("tried to index a non-array")

        if self.index < 0 or self.index >= len(struct):
            raise InvalidFilterApplication(
                f"tried to get index={self.index} from an array of length {len(struct)}"
            )

        return struct[self.index]

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, ArrayIndexFilter):
            return False

        return self.index == other.index


class ObjectIdentifierFilter(Filter):
    def __init__(self, identifier: str, is_optional: bool):
        self.identifier = identifier
        self.is_optional = is_optional

    def apply(self, struct: JSONStruct) -> JSONStruct:
        if not isinstance(struct, dict):
            if self.is_optional:
                return None
            raise InvalidFilterApplication(
                "tried to get an identifier from a non-object"
            )

        return struct.get(self.identifier, None)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, ObjectIdentifierFilter):
            return False

        return (
            self.identifier == other.identifier
            and self.is_optional == other.is_optional
        )


class PipeFilter(Filter):
    def apply(self, struct: JSONStruct) -> JSONStruct:
        return struct

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, PipeFilter):
            return False

        return True


array_collect_pattern = re.compile(r"^\[(.+)\]$")
array_index_pattern = re.compile(r"^\.\[(\d+)\]")
object_identifier_pattern_1 = re.compile(r'^\.\["([a-zA-Z]+)"\](\?)?')
object_identifier_pattern_2 = re.compile(r"^\.([a-zA-Z]+)(\?)?")
pipe_pattern = re.compile(r"\s+\|\s+")


def get_filter(argument: str) -> tuple[Optional[Filter], str]:
    if argument == "" or argument == ".":
        return IdentityFilter(), ""

    if m := array_index_pattern.match(argument):
        index = int(m.group(1))
        matched_length = len(m.group(0))
        remaining = argument[matched_length:]
        return ArrayIndexFilter(index), remaining

    if (m := object_identifier_pattern_1.match(argument)) or (
        m := object_identifier_pattern_2.match(argument)
    ):
        identifier = m.group(1)
        is_optional = m.group(2) is not None
        matched_length = len(m.group(0))
        remaining = argument[matched_length:]
        return ObjectIdentifierFilter(identifier, is_optional), remaining

    if m := pipe_pattern.match(argument):
        matched_length = len(m.group(0))
        remaining = argument[matched_length:]
        return PipeFilter(), remaining

    raise InvalidFilterArgument()


def get_filters(argument: str) -> tuple[list[Filter], Flags]:
    flags = Flags()

    if m := array_collect_pattern.match(argument):
        argument = m.group(1)
        flags.should_collect_into_array = True

    filters: list[Filter] = []

    while argument:
        filter_, argument = get_filter(argument)
        if filter_ is None:
            break
        filters.append(filter_)

    return filters, flags


def apply_filters(
    struct: JSONStruct, filters: list[Filter], flags: Flags
) -> JSONStruct:
    for filter_ in filters:
        struct = filter_.apply(struct)

    if flags.should_collect_into_array:
        pass

    return struct


# TODO: this whole array spreading and collecting thing breaks my whole system here :/
# i dont like it, maybe should think of a more elegant solution
# using Flags seems like a dirt patch :/
