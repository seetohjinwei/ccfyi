from src.states.array_state import ArrayState
from src.states.object_state import ObjectState
from src.states.state import StateTransitionResult
from src.states.string_state import StringState
from src.states.value_state import ValueState
from src.states.whitespace_state import WhitespaceState
from src.str_slice import Str
import unittest


# Useful site: https://www.pythonescaper.com


def assert_equal_result(
    t: unittest.TestCase, actual: StateTransitionResult, expected: StateTransitionResult
) -> None:
    message = f"expected {expected}, but got {actual}"

    t.assertEqual(actual.is_success, expected.is_success, message)
    t.assertEqual(actual.new_txt, expected.new_txt, message)

    if isinstance(actual.json_struct, dict) and isinstance(expected.json_struct, dict):
        t.assertDictEqual(actual.json_struct, expected.json_struct, message)
    elif isinstance(actual.json_struct, list) and isinstance(
        expected.json_struct, list
    ):
        t.assertListEqual(actual.json_struct, expected.json_struct, message)
    elif isinstance(actual.json_struct, float) and isinstance(
        expected.json_struct, float
    ):
        t.assertAlmostEqual(actual.json_struct, expected.json_struct, msg=message)
        pass
    else:
        t.assertEqual(actual.json_struct, expected.json_struct, message)


class TestArrayState(unittest.TestCase):
    def test_basic_array(self):
        actual = ArrayState.transition(Str('[  "1","2",   "3"\n]'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=["1", "2", "3"],
        )

        assert_equal_result(self, actual, expected)

    def test_mixed_array(self):
        actual = ArrayState.transition(Str('[  "1","2",   "3",\n true, null, false]'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=["1", "2", "3", True, None, False],
        )

        assert_equal_result(self, actual, expected)

    def test_empty_array(self):
        actual = ArrayState.transition(Str("[]"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=[],
        )

        assert_equal_result(self, actual, expected)


class TestValueState(unittest.TestCase):
    def test_string(self):
        actual = ValueState.transition(Str('  "str"'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct="str",
        )

        assert_equal_result(self, actual, expected)

    def test_no_value(self):
        actual = ValueState.transition(Str("   "))
        expected = StateTransitionResult(
            is_success=False,
            new_txt=Str("   "),
            json_struct=None,
        )

        assert_equal_result(self, actual, expected)


class TestObjectState(unittest.TestCase):
    def test_empty_object(self):
        actual = ObjectState.transition(Str("{}"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct={},
        )

        assert_equal_result(self, actual, expected)

    def test_basic_object(self):
        actual = ObjectState.transition(Str('{"a": "b"}'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct={"a": "b"},
        )

        assert_equal_result(self, actual, expected)

    def test_complex_object(self):
        actual = ObjectState.transition(
            Str('{"key": \t\t"value", "false": false, "null": \n    null, "num": 123}')
        )
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct={"key": "value", "false": False, "null": None, "num": 123},
        )

        assert_equal_result(self, actual, expected)

    def test_trailing_comma(self):
        actual = ObjectState.transition(Str('{"key": "value",}'))
        expected = StateTransitionResult(
            is_success=False,
            new_txt=Str('{"key": "value",}'),
            json_struct=None,
        )

        assert_equal_result(self, actual, expected)

    def test_slash(self):
        actual = ObjectState.transition(Str('{"slash": "/ & \\/"}'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct={"slash": "/ & /"},
        )

        assert_equal_result(self, actual, expected)


class TestStringState(unittest.TestCase):
    def test_basic_string(self):
        actual = StringState.transition(Str('"123"'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct="123",
        )

        assert_equal_result(self, actual, expected)

    def test_complex_string(self):
        r"\\"

        actual = StringState.transition(Str('"escape?42\\n": {}'))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(": {}"),
            json_struct="escape?42\n",
        )

        assert_equal_result(self, actual, expected)

    def test_empty_string(self):
        actual = StringState.transition(Str('""'))
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

    def test_control(self):
        actual = StringState.transition(Str('"linebreak\n"'))
        expected = StateTransitionResult(
            is_success=False,
            new_txt=Str('"linebreak\n"'),
            json_struct=None,
        )

        assert_equal_result(self, actual, expected)

    def test_control2(self):
        from src.states.string_state import _StringControlState

        actual = _StringControlState.transition(Str("\\x15"))

        self.assertTrue(actual.is_success)

    def test_unicode(self):
        from src.states.string_state import _StringEscapeUnicodeState

        actual = _StringEscapeUnicodeState.transition(Str("uFCDE"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=chr(int("FCDE", 16)),
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


class TestNumberState(unittest.TestCase):
    def test_positive_integer(self):
        actual = ValueState.transition(Str("1234789"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=1234789,
        )

        assert_equal_result(self, actual, expected)

    def test_negative_integer(self):
        actual = ValueState.transition(Str("-1234789"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=-1234789,
        )

        assert_equal_result(self, actual, expected)

    def test_positive_zero(self):
        actual = ValueState.transition(Str("0"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=0,
        )

        assert_equal_result(self, actual, expected)

    def test_negative_zero(self):
        actual = ValueState.transition(Str("-0"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=0,
        )

        assert_equal_result(self, actual, expected)

    def test_fraction(self):
        actual = ValueState.transition(Str("4.321"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=4.321,
        )

        assert_equal_result(self, actual, expected)

    def test_exponent1(self):
        actual = ValueState.transition(Str("32E-10"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=32e-10,
        )

        assert_equal_result(self, actual, expected)

    def test_exponent2(self):
        actual = ValueState.transition(Str("32e4"))
        expected = StateTransitionResult(
            is_success=True,
            new_txt=Str(""),
            json_struct=32e4,
        )

        assert_equal_result(self, actual, expected)
