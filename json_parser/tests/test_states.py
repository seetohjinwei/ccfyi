from src.str_slice import Str
from src.states import WhitespaceState
import unittest


class TestWhitespaceState(unittest.TestCase):
    def test_leading_whitespace(self):
        s = Str("   123")

        res = WhitespaceState.transition(s)
        actual = res.txt
        expected = Str("123")

        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")

    def test_empty(self):
        s = Str("")

        res = WhitespaceState.transition(s)
        actual = res.txt
        expected = Str("")

        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")

    def test_no_whitespace(self):
        s = Str("123")

        res = WhitespaceState.transition(s)
        actual = res.txt
        expected = Str("123")

        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")

    def test_non_leading_whitespace(self):
        s = Str("123 456")

        res = WhitespaceState.transition(s)
        actual = res.txt
        expected = Str("123 456")

        self.assertEqual(actual, expected, f"expected {expected}, but got {actual}")
