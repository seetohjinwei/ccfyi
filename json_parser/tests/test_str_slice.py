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
