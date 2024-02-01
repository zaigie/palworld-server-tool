from typing import Any, Sequence

from lib.archive import *


def decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    if type_name != "ArrayProperty":
        raise Exception(f"Expected ArrayProperty, got {type_name}")
    value = reader.property(type_name, size, path, nested_caller_path=path)
    char_bytes = value["value"]["values"]
    value["value"] = decode_bytes(char_bytes)
    return value


def decode_bytes(char_bytes: Sequence[int]) -> dict[str, Any]:
    reader = FArchiveReader(bytes(char_bytes), debug=False)
    char_data = {}
    char_data["object"] = reader.properties_until_end()
    char_data["unknown_bytes"] = reader.byte_list(4)
    char_data["group_id"] = reader.guid()
    if not reader.eof():
        raise Exception("Warning: EOF not reached")
    return char_data


def encode(
    writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
) -> int:
    if property_type != "ArrayProperty":
        raise Exception(f"Expected ArrayProperty, got {property_type}")
    del properties["custom_type"]
    encoded_bytes = encode_bytes(properties["value"])
    properties["value"] = {"values": [b for b in encoded_bytes]}
    return writer.property_inner(property_type, properties)


def encode_bytes(p: dict[str, Any]) -> bytes:
    writer = FArchiveWriter()
    writer.properties(p["object"])
    writer.write(bytes(p["unknown_bytes"]))
    writer.guid(p["group_id"])
    encoded_bytes = writer.bytes()
    return encoded_bytes
