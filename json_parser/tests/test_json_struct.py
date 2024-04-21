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

    def test_invalid_json_1(self):
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

    def test_valid_json_1(self):
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

    def test_invalid_json_2_0(self):
        self.skipTest("TODO")
        tc: TestCase[JSONStruct] = TestCase(
                name="valid_json",
                path="step2/invalid.json",
                expected=None,
            )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_invalid_json_2_1(self):
        self.skipTest("TODO")
        tc: TestCase[JSONStruct] = TestCase(
                name="valid_json",
                path="step2/invalid2.json",
                expected=None,
            )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json_2_0(self):
        self.skipTest("TODO")
        tc: TestCase[JSONStruct] = TestCase(
                name="valid_json",
                path="step2/valid.json",
                expected={"key": "value"},
            )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json_2_1(self):
        self.skipTest("TODO")
        tc: TestCase[JSONStruct] = TestCase(
                name="valid_json",
                path="step2/valid2.json",
                expected={"key": "value", "key2": "value"},
            )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )
