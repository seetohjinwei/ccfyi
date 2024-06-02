#!../venv/bin/python3

from src.json_struct import parse
from dataclasses import dataclass
import sys


@dataclass
class Args:
    input_file: str

    @staticmethod
    def get() -> "Args":
        if len(sys.argv) < 2:
            raise ValueError(f"""
filepath is missing!
Usage: `{sys.argv[0]} <filepath>`
or `make run input=<filepath>`
""")

        input_file = sys.argv[1]
        return Args(input_file=input_file)


def get_file_contents(path: str) -> str:
    with open(path, "r") as f:
        return f.read().strip()


def main(args: Args) -> None:
    txt = get_file_contents(args.input_file)
    result = parse(txt)
    print(result)


if __name__ == "__main__":
    args = Args.get()
    main(args)
