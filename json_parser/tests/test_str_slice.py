from src.str_slice import Str
import unittest


class TestStrSlice(unittest.TestCase):
    def test_iter(self):
        expected = "ABCdef123"
        s = Str(expected)

        actual = ""
        for c in s:
            actual += c

        self.assertEqual(actual, expected)

    def test_enumerate(self):
        s = Str("ABCdef123")

        expected = 0
        for actual, _ in enumerate(s):
            self.assertEqual(actual, expected)
            expected += 1

    def test_substring(self):
        s = Str("abcde")

        expected = Str("de")
        actual = s.substring(1).substring(1).substring(1)
        self.assertEqual(actual, expected)

    def test_startswith1(self):
        s = Str("abcde")

        actual = s.starts_with("ab")
        self.assertTrue(actual)

    def test_startswith2(self):
        s = Str("bacde")

        actual = s.starts_with("ab")
        self.assertFalse(actual)

    def test_startswith3(self):
        s = Str("abcde")

        actual = s.starts_with(Str("ab"))
        self.assertTrue(actual)

    def test_startswith4(self):
        s = Str("bacde")

        actual = s.starts_with(Str("ab"))
        self.assertFalse(actual)
