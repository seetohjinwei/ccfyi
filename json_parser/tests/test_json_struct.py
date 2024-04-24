from dataclasses import dataclass
from typing import Generic, TypeVar
import unittest

from src.json_struct import InvalidJSONStruct, JSONStruct, parse


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
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

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
        tc: TestCase[JSONStruct] = TestCase(
            name="valid_json",
            path="step2/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.path)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_invalid_json_2_1(self):
        tc: TestCase[JSONStruct] = TestCase(
            name="valid_json",
            path="step2/invalid2.json",
            expected=None,
        )

        txt = get_test_case(tc.path)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_2_0(self):
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

    def test_invalid_json_3(self):
        tc: TestCase[JSONStruct] = TestCase(
            name="invalid_json",
            path="step3/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.path)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_3(self):
        tc: TestCase[JSONStruct] = TestCase(
            name="valid_json",
            path="step3/valid.json",
            expected={
                "key1": True,
                "key2": False,
                "key3": None,
                "key4": "value",
                "key5": 101,
            },
        )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_invalid_json_4(self):
        tc: TestCase[JSONStruct] = TestCase(
            name="invalid_json",
            path="step4/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.path)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_4_0(self):
        tc: TestCase[JSONStruct] = TestCase(
            name="valid_json",
            path="step4/valid.json",
            expected={"key": "value", "key-n": 101, "key-o": {}, "key-l": []},
        )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json_4_1(self):
        tc: TestCase[JSONStruct] = TestCase(
            name="valid_json",
            path="step4/valid2.json",
            expected={
                "key": "value",
                "key-n": 101,
                "key-o": {"inner key": "inner value"},
                "key-l": ["list value"],
            },
        )

        txt = get_test_case(tc.path)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )
