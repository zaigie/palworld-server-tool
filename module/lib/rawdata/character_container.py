from typing import Any, Sequence

from lib.archive import *


def decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    if type_name != "ArrayProperty":
        raise Exception(f"Expected ArrayProperty, got {type_name}")
    value = reader.property(type_name, size, path, nested_caller_path=path)
    data_bytes = value["value"]["values"]
    value["value"] = decode_bytes(data_bytes)
    return value


def decode_bytes(c_bytes: Sequence[int]) -> dict[str, Any]:
    if len(c_bytes) == 0:
        return None
    reader = FArchiveReader(bytes(c_bytes), debug=False)
    data = {}
    data["player_uid"] = reader.guid()
    data["instance_id"] = reader.guid()
    data["permission_tribe_id"] = reader.byte()
    if not reader.eof():
        raise Exception("Warning: EOF not reached")
    return data


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
    if p is None:
        return bytes()
    writer = FArchiveWriter()
    writer.guid(p["player_uid"])
    writer.guid(p["instance_id"])
    writer.byte(p["permission_tribe_id"])
    encoded_bytes = writer.bytes()
    return encoded_bytes
