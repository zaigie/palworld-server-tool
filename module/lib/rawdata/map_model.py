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


def decode_bytes(m_bytes: Sequence[int]) -> dict[str, Any]:
    reader = FArchiveReader(bytes(m_bytes), debug=False)
    data = {}
    data["instance_id"] = reader.guid()
    data["concrete_model_instance_id"] = reader.guid()
    data["base_camp_id_belong_to"] = reader.guid()
    data["group_id_belong_to"] = reader.guid()
    data["hp"] = {
        "current": reader.i32(),
        "max": reader.i32(),
    }
    data["initital_transform_cache"] = reader.ftransform()
    data["repair_work_id"] = reader.guid()
    data["owner_spawner_level_object_instance_id"] = reader.guid()
    data["owner_instance_id"] = reader.guid()
    data["build_player_uid"] = reader.guid()
    data["interact_restrict_type"] = reader.byte()
    data["stage_instance_id_belong_to"] = {
        "id": reader.guid(),
        "valid": reader.u32() > 0,
    }
    data["created_at"] = reader.i64()
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

    writer.guid(p["instance_id"])
    writer.guid(p["concrete_model_instance_id"])
    writer.guid(p["base_camp_id_belong_to"])
    writer.guid(p["group_id_belong_to"])

    writer.i32(p["hp"]["current"])
    writer.i32(p["hp"]["max"])

    writer.ftransform(p["initital_transform_cache"])

    writer.guid(p["repair_work_id"])
    writer.guid(p["owner_spawner_level_object_instance_id"])
    writer.guid(p["owner_instance_id"])
    writer.guid(p["build_player_uid"])

    writer.byte(p["interact_restrict_type"])

    writer.guid(p["stage_instance_id_belong_to"]["id"])
    writer.u32(1 if p["stage_instance_id_belong_to"]["valid"] else 0)

    writer.i64(p["created_at"])

    encoded_bytes = writer.bytes()
    return encoded_bytes
