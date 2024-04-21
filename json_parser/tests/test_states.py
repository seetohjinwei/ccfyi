from src.str_slice import Str
from src.states import StateTransitionResult, StringState, WhitespaceState
import unittest


def assert_equal_result(t: unittest.TestCase, actual: StateTransitionResult, expected: StateTransitionResult) -> None:
    t.assertEqual(actual.is_success, expected.is_success)
    t.assertEqual(actual.new_txt, expected.new_txt)

    if isinstance(actual.json_struct, dict) and isinstance(expected.json_struct, dict):
        t.assertDictEqual(actual.json_struct, expected.json_struct)
    elif isinstance(actual.json_struct, list) and isinstance(expected.json_struct, list):
        t.assertListEqual(actual.json_struct, expected.json_struct)
    else:
        t.assertEqual(actual.json_struct, expected.json_struct)


class TestStringState(unittest.TestCase):
    def test_basic_string(self):
        actual = StringState.transition(Str("\"123\""))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct="123",
        )

        assert_equal_result(self, actual, expected)

    def test_complex_string(self):
        actual = StringState.transition(Str("\"\\\"escape?42\n\": {}"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(": {}"),
            json_struct="\"escape?42\n",
        )

        assert_equal_result(self, actual, expected)

    def test_empty_string(self):
        actual = StringState.transition(Str("\"\""))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct="",
        )

        assert_equal_result(self, actual, expected)

    def test_not_string1(self):
        actual = StringState.transition(Str(""))
        expected = StateTransitionResult(
            is_success=False,
            new_txt=Str(""),
            json_struct=None,
        )

        assert_equal_result(self, actual, expected)

    def test_not_string2(self):
        actual = StringState.transition(Str("{}"))
        expected = StateTransitionResult(
            is_success=False,
            new_txt=Str("{}"),
            json_struct=None,
        )

        assert_equal_result(self, actual, expected)


class TestWhitespaceState(unittest.TestCase):
    def test_leading_whitespace(self):
        s = Str("   123")

        res = WhitespaceState.transition(s)
        actual = res.new_txt
        expected = Str("123")

        self.assertTrue(res.is_success)
        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")

    def test_empty(self):
        s = Str("")

        res = WhitespaceState.transition(s)
        actual = res.new_txt
        expected = Str("")

        self.assertFalse(res.is_success)
        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")

    def test_no_whitespace(self):
        s = Str("123")

        res = WhitespaceState.transition(s)
        actual = res.new_txt
        expected = Str("123")

        self.assertFalse(res.is_success)
        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")

    def test_non_leading_whitespace(self):
        s = Str("123 456")

        res = WhitespaceState.transition(s)
        actual = res.new_txt
        expected = Str("123 456")

        self.assertFalse(res.is_success)
        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")
