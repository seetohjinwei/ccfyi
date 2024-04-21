class Str:
    """
    Str is a wrapper around a `str` object, allowing cheap substring operations.
    Similar idea to string slices in Go and Rust.
    """

    s: str
    length: int
    lo: int  # inclusive
    hi: int  # exclusive

    def __init__(self, s: str):
        self.s = s
        self.length = len(s)
        self.lo = 0
        self.hi = self.length

    def clone(self) -> "Str":
        ret = Str(self.s)

        ret.lo = self.lo
        ret.hi = self.hi

        return ret

    def substring(self, lo: int = 0, hi: int | None = None) -> "Str":
        if hi is None:
            hi = len(self)

        ret = self.clone()
        ret.lo = self.lo + lo
        ret.hi = self.lo + hi
        return ret

    def at(self, index: int) -> str:
        if index < 0:
            return self.s[self.hi + index]

        return self.s[self.lo + index]

    def __len__(self) -> int:
        return self.hi - self.lo

    def __iter__(self) -> "StrIter":
        return StrIter(self)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, Str):
            return False

        if len(self) != len(other):
            return False

        for x, y in zip(self, other):
            if x != y:
                return False

        return True

    def __str__(self) -> str:
        s = self.s[self.lo:self.hi]
        return f"\"{s}\""


class StrIter:
    def __init__(self, s: Str):
        self.index = 0
        self.s = s

    def __iter__(self):
        return self

    def __next__(self) -> str:
        if self.index >= len(self.s):
            raise StopIteration

        ret = self.s.at(self.index)
        self.index += 1
        return ret
