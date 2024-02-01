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
    data["id"] = reader.guid()
    data["name"] = reader.fstring()
    data["state"] = reader.byte()
    data["transform"] = reader.ftransform()
    data["area_range"] = reader.float()
    data["group_id_belong_to"] = reader.guid()
    data["fast_travel_local_transform"] = reader.ftransform()
    data["owner_map_object_instance_id"] = reader.guid()
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
    writer.guid(p["id"])
    writer.fstring(p["name"])
    writer.byte(p["state"])
    writer.ftransform(p["transform"])
    writer.float(p["area_range"])
    writer.guid(p["group_id_belong_to"])
    writer.ftransform(p["fast_travel_local_transform"])
    writer.guid(p["owner_map_object_instance_id"])
    encoded_bytes = writer.bytes()
    return encoded_bytes
