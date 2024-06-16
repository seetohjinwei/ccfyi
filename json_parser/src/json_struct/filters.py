from abc import ABC, abstractmethod
import re
from typing import Optional

from src.json_struct.json_struct import JSONStruct


class InvalidFilterApplication(Exception):
    pass


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


array_index_pattern = re.compile(r"^\.\[(\d+)\]")
object_identifier_pattern_1 = re.compile(r'^\.\["([a-zA-Z]+)"\](\?)?')
object_identifier_pattern_2 = re.compile(r"^\.([a-zA-Z]+)(\?)?")


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

    return None, argument


def get_filters(argument: str) -> list[Filter]:
    filters: list[Filter] = []

    while argument:
        filter_, argument = get_filter(argument)
        if filter_ is None:
            break
        filters.append(filter_)

    return filters


def apply_filters(struct: JSONStruct, filters: list[Filter]) -> JSONStruct:
    for filter in filters:
        struct = filter.apply(struct)

    return struct
