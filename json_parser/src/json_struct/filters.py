from abc import ABC, abstractmethod
import re

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


array_index_pattern = re.compile(r"^\.\[(\d+)\]$")


def get_filters(argument: str) -> list[Filter]:
    if argument == "" or argument == ".":
        return [IdentityFilter()]

    if m := array_index_pattern.match(argument):
        index = int(m.group(1))
        return [ArrayIndexFilter(index)]

    return []


def apply_filters(struct: JSONStruct, filters: list[Filter]) -> JSONStruct:
    for filter in filters:
        struct = filter.apply(struct)

    return struct
