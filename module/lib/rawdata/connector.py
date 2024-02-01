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


def connect_info_item_reader(reader: FArchiveReader) -> dict[str, Any]:
    return {
        "connect_to_model_instance_id": reader.guid(),
        "index": reader.byte(),
    }


def connect_info_item_writer(writer: FArchiveWriter, properties: dict[str, Any]):
    writer.guid(properties["connect_to_model_instance_id"])
    writer.byte(properties["index"])


def decode_bytes(c_bytes: Sequence[int]) -> dict[str, Any]:
    if len(c_bytes) == 0:
        return None
    reader = FArchiveReader(bytes(c_bytes), debug=False)
    data = {}
    data["supported_level"] = reader.i32()
    data["connect"] = {
        "index": reader.byte(),
        "any_place": reader.tarray(connect_info_item_reader),
    }
    # We are guessing here, we don't have information about the type without mapping object names -> types
    # Stairs have 2 connectors (up and down),
    # Roofs have 4 connectors (front, back, right, left)
    if not reader.eof():
        data["other_connectors"] = []
        while not reader.eof():
            data["other_connectors"].append(
                {
                    "index": reader.byte(),
                    "connect": reader.tarray(connect_info_item_reader),
                }
            )
        if len(data["other_connectors"]) not in [2, 4]:
            print(
                f"Warning: unknown connector type with {len(data['other_connectors'])} connectors"
            )
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
    writer.i32(p["supported_level"])
    writer.byte(p["connect"]["index"])
    writer.tarray(connect_info_item_writer, p["connect"]["any_place"])
    if "other_connectors" in p:
        for other in p["other_connectors"]:
            writer.byte(other["index"])
            writer.tarray(connect_info_item_writer, other["connect"])
    encoded_bytes = writer.bytes()
    return encoded_bytes
