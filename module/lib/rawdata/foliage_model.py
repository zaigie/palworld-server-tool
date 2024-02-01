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


def decode_bytes(b_bytes: Sequence[int]) -> dict[str, Any]:
    reader = FArchiveReader(bytes(b_bytes), debug=False)
    data = {}
    data["model_id"] = reader.fstring()
    data["foliage_preset_type"] = reader.byte()
    data["cell_coord"] = {
        "x": reader.i64(),
        "y": reader.i64(),
        "z": reader.i64(),
    }
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
    writer = FArchiveWriter()

    writer.fstring(p["model_id"])
    writer.byte(p["foliage_preset_type"])
    writer.i64(p["cell_coord"]["x"])
    writer.i64(p["cell_coord"]["y"])
    writer.i64(p["cell_coord"]["z"])

    encoded_bytes = writer.bytes()
    return encoded_bytes
