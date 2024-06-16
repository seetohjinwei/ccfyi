from dataclasses import dataclass
from src.json_struct.filters import (
    ArrayIndexFilter,
    Filter,
    IdentityFilter,
    InvalidFilterApplication,
    apply_filters,
    get_filters,
)
from src.json_struct.json_struct import (
    InvalidJSONStruct,
    JSONStruct,
    parse,
    pretty_print,
)
from typing import Generic, Optional, TypeVar
import unittest
import os


def get_test_case(path: str) -> str:
    with open("test_json/" + path, "r") as f:
        return f.read().strip()


V = TypeVar("V")
E = TypeVar("E")


@dataclass
class TestCase(Generic[V, E]):
    input: V
    expected: E
    name: Optional[str] = None


class TestJSONStruct(unittest.TestCase):
    def test_unittest(self):
        self.assertEqual(1, 1)

    def test_invalid_json_1(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="invalid_json",
            input="step1/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.input)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_1(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step1/valid.json",
            expected={},
        )

        txt = get_test_case(tc.input)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_invalid_json_2_0(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step2/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.input)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_invalid_json_2_1(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step2/invalid2.json",
            expected=None,
        )

        txt = get_test_case(tc.input)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_2_0(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step2/valid.json",
            expected={"key": "value"},
        )

        txt = get_test_case(tc.input)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json_2_1(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step2/valid2.json",
            expected={"key": "value", "key2": "value"},
        )

        txt = get_test_case(tc.input)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_invalid_json_3(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="invalid_json",
            input="step3/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.input)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_3(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step3/valid.json",
            expected={
                "key1": True,
                "key2": False,
                "key3": None,
                "key4": "value",
                "key5": 101,
            },
        )

        txt = get_test_case(tc.input)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_invalid_json_4(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="invalid_json",
            input="step4/invalid.json",
            expected=None,
        )

        txt = get_test_case(tc.input)
        self.assertRaises(InvalidJSONStruct, lambda: parse(txt))

    def test_valid_json_4_0(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step4/valid.json",
            expected={"key": "value", "key-n": 101, "key-o": {}, "key-l": []},
        )

        txt = get_test_case(tc.input)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json_4_1(self):
        tc: TestCase[str, JSONStruct] = TestCase(
            name="valid_json",
            input="step4/valid2.json",
            expected={
                "key": "value",
                "key-n": 101,
                "key-o": {"inner key": "inner value"},
                "key-l": ["list value"],
            },
        )

        txt = get_test_case(tc.input)
        result = parse(txt)
        self.assertEqual(
            result, tc.expected, f"expected {tc.expected}, but got {result}"
        )

    def test_valid_json_others_1(self):
        # from `curl -sL 'https://api.github.com/repos/CodingChallegesFYI/SharedSolutions/commits?per_page=3'`

        txt = get_test_case("others/valid1.json")
        try:
            parse(txt)
        except InvalidJSONStruct:
            self.fail("expected to parse without exception")

    def test_full_suite(self):
        # from http://www.json.org/JSON_checker/test.zip
        # modified cases: fail18 -> pass18

        directory = "full_suite/"
        filepaths = os.listdir("test_json/" + directory)

        failed_cases: list[str] = []

        for f in filepaths:
            txt = get_test_case(directory + f)

            if f.startswith("fail"):
                try:
                    result = parse(txt)
                    failed_cases.append(
                        f"expected {f} to fail: {txt}, but got {result}"
                    )
                except InvalidJSONStruct:
                    pass
            else:
                try:
                    parse(txt)
                except InvalidJSONStruct as e:
                    failed_cases.append(f"expected {f} to succeed: {txt}, but got {e}")

        if failed_cases:
            failed_cases.append(f"failures={len(failed_cases)}, total={len(filepaths)}")
            message = "\n".join(failed_cases)
            self.fail(message)


class TestJSONStruct_pretty_print(unittest.TestCase):
    def test_bool(self):
        test_cases: list[TestCase[JSONStruct, str]] = [
            TestCase[JSONStruct, str](
                input=False,
                expected="false",
            ),
            TestCase[JSONStruct, str](
                input=True,
                expected="true",
            ),
        ]

        for tc in test_cases:
            actual = pretty_print(tc.input)
            self.assertEqual(tc.expected, actual)

    def test_none(self):
        test_cases: list[TestCase[JSONStruct, str]] = [
            TestCase[JSONStruct, str](
                input=None,
                expected="null",
            ),
        ]

        for tc in test_cases:
            actual = pretty_print(tc.input)
            self.assertEqual(tc.expected, actual)

    def test_float(self):
        test_cases: list[TestCase[JSONStruct, str]] = [
            TestCase[JSONStruct, str](
                input=3.54,
                expected="3.54",
            ),
            TestCase[JSONStruct, str](
                input=-241.2,
                expected="-241.2",
            ),
        ]

        for tc in test_cases:
            actual = pretty_print(tc.input)
            self.assertEqual(tc.expected, actual)

    def test_str(self):
        test_cases: list[TestCase[JSONStruct, str]] = [
            TestCase[JSONStruct, str](
                input="",
                expected='""',
            ),
            TestCase[JSONStruct, str](
                input="json",
                expected='"json"',
            ),
        ]

        for tc in test_cases:
            actual = pretty_print(tc.input)
            self.assertEqual(tc.expected, actual)

    def test_array(self):
        test_cases: list[TestCase[JSONStruct, str]] = [
            TestCase[JSONStruct, str](
                input=[],
                expected="[]",
            ),
            TestCase[JSONStruct, str](
                input=[1],
                expected="""[
  1
]""",
            ),
            TestCase[JSONStruct, str](
                input=[1, 2, 3],
                expected="""[
  1,
  2,
  3
]""",
            ),
            TestCase[JSONStruct, str](
                input=[1, [2, 3], [4, 5, [6, 7]]],
                expected="""[
  1,
  [
    2,
    3
  ],
  [
    4,
    5,
    [
      6,
      7
    ]
  ]
]""",
            ),
            TestCase[JSONStruct, str](
                input=[{"a": "b", "c": "d"}, 2, "xd"],
                expected="""[
  {
    "a": "b",
    "c": "d"
  },
  2,
  "xd"
]""",
            ),
        ]

        for tc in test_cases:
            actual = pretty_print(tc.input)
            self.assertEqual(tc.expected, actual)

    def test_dict(self):
        test_cases: list[TestCase[JSONStruct, str]] = [
            TestCase[JSONStruct, str](
                input={},
                expected="{}",
            ),
            TestCase[JSONStruct, str](
                input={"2": 4},
                expected="""{
  "2": 4
}""",
            ),
            TestCase[JSONStruct, str](
                input={"a": 1, "b": 2, "c": 3},
                expected="""{
  "a": 1,
  "b": 2,
  "c": 3
}""",
            ),
            TestCase[JSONStruct, str](
                input={"a": 1, "b": 2, "c": [3, 4, 5, {"key": "value"}]},
                expected="""{
  "a": 1,
  "b": 2,
  "c": [
    3,
    4,
    5,
    {
      "key": "value"
    }
  ]
}""",
            ),
        ]

        for tc in test_cases:
            actual = pretty_print(tc.input)
            self.assertEqual(tc.expected, actual)


class TestJSONStruct_get_filters(unittest.TestCase):
    def test_filters(self):
        test_cases: list[TestCase[str, list[Filter]]] = [
            TestCase(
                input=".",
                expected=[IdentityFilter()],
            ),
            TestCase(
                input=".[0]",
                expected=[ArrayIndexFilter(0)],
            ),
            TestCase(
                input=".[1000]",
                expected=[ArrayIndexFilter(1000)],
            ),
        ]

        for tc in test_cases:
            actual = get_filters(tc.input)
            self.assertEqual(tc.expected, actual)


class TestJSONStruct_apply_filters(unittest.TestCase):
    def test_filters(self):
        test_cases: list[TestCase[tuple[JSONStruct, list[Filter]], JSONStruct]] = [
            TestCase(
                input=({}, [IdentityFilter()]),
                expected={},
            ),
            TestCase(
                input=([1], [ArrayIndexFilter(0)]),
                expected=1,
            ),
        ]

        for tc in test_cases:
            actual = apply_filters(tc.input[0], tc.input[1])
            self.assertEqual(tc.expected, actual)

    def test_exceptions(self):
        test_cases: list[tuple[JSONStruct, list[Filter]]] = [
            ({}, [ArrayIndexFilter(0)]),
            ([1, 2, 3], [ArrayIndexFilter(10)]),
        ]

        for tc in test_cases:
            self.assertRaises(
                InvalidFilterApplication, lambda: apply_filters(tc[0], tc[1])
            )
