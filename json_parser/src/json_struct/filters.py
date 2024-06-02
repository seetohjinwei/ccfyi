from abc import ABC, abstractmethod

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


def get_filters(argument: str) -> list[Filter]:
    if argument == "" or argument == ".":
        return [IdentityFilter()]

    return []


def apply_filters(struct: JSONStruct, filters: list[Filter]) -> JSONStruct:
    for filter in filters:
        struct = filter.apply(struct)

    return struct
