from dataclasses import dataclass
from typing import Generic, TypeVar
import unittest

from src.json_struct import JSONStruct, parse


def get_test_case(path: str) -> str:
    with open("test_json/" + path, "r") as f:
        return f.read()


E = TypeVar("E")


@dataclass
class TestCase(Generic[E]):
    name: str
    path: str
    expected: E


class TestJSONStruct(unittest.TestCase):
    def test_unittest(self):
        self.assertEqual(1, 1)

    def test_invalid_json(self):
        tc: TestCase[JSONStruct] = TestCase(
                name="invalid_json",
                path="step1/invalid.json",
                expected=None,
            )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json(self):
        tc: TestCase[JSONStruct] = TestCase(
                name="valid_json",
                path="step1/valid.json",
                expected={},
            )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )
