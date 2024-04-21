# https://www.rfc-editor.org/rfc/rfc8259.html#section-3
JSONStruct = bool | None | dict[str, "JSONStruct"] | list["JSONStruct"] | float | str


def parse(text: str) -> JSONStruct:
    if not text:
        return None

    return {}
