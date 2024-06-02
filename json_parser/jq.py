#!../venv/bin/python3

from src.json_struct.json_struct import JSONStruct, parse, pretty_print
import sys


def get_input_contents() -> str:
    # this is guaranteed to read until EOF
    return sys.stdin.read().strip()


def main() -> None:
    txt = get_input_contents()
    struct: JSONStruct = parse(txt)
    # TODO: parse arguments into functions
    # TODO: apply functions on struct
    result = pretty_print(struct)
    print(result)


if __name__ == "__main__":
    main()
