#!../venv/bin/python3

from src.json_struct.filters import apply_filters, get_filters
from src.json_struct.json_struct import JSONStruct, parse, pretty_print
import sys


def get_input_contents() -> str:
    # this is guaranteed to read until EOF
    return sys.stdin.read().strip()


def main() -> None:
    txt = get_input_contents()
    struct: JSONStruct = parse(txt)

    argument: str
    if len(sys.argv) < 2:
        argument = ""
    else:
        argument = sys.argv[1]

    filters, flags = get_filters(argument)
    struct = apply_filters(struct, filters, flags)

    result = pretty_print(struct)
    print(result)


if __name__ == "__main__":
    main()
