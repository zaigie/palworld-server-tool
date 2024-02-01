from typing import Any, Sequence

from lib.archive import *

WORK_BASE_TYPES = set(
    [
        # "EPalWorkableType::Illegal",
        "EPalWorkableType::Progress",
        # "EPalWorkableType::CollectItem",
        # "EPalWorkableType::TransportItem",
        "EPalWorkableType::TransportItemInBaseCamp",
        "EPalWorkableType::ReviveCharacter",
        # "EPalWorkableType::CollectResource",
        "EPalWorkableType::LevelObject",
        "EPalWorkableType::Repair",
        "EPalWorkableType::Defense",
        "EPalWorkableType::BootUp",
        "EPalWorkableType::OnlyJoin",
        "EPalWorkableType::OnlyJoinAndWalkAround",
        "EPalWorkableType::RemoveMapObjectEffect",
        "EPalWorkableType::MonsterFarm",
    ]
)


def decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    if type_name != "ArrayProperty":
        raise Exception(f"Expected ArrayProperty, got {type_name}")
    value = reader.property(type_name, size, path, nested_caller_path=path)
    for work_element in value["value"]["values"]:
        work_bytes = work_element["RawData"]["value"]["values"]
        work_type = work_element["WorkableType"]["value"]["value"]
        work_element["RawData"]["value"] = decode_bytes(work_bytes, work_type)
        for work_assign in work_element["WorkAssignMap"]["value"]:
            work_assign_bytes = work_assign["value"]["RawData"]["value"]["values"]
            work_assign["value"]["RawData"]["value"] = decode_work_assign_bytes(
                work_assign_bytes
            )
    return value


def decode_bytes(b_bytes: Sequence[int], work_type: str) -> dict[str, Any]:
    reader = FArchiveReader(bytes(b_bytes), debug=False)
    data = {}
    # Handle base serialization
    if work_type in WORK_BASE_TYPES:
        data["id"] = reader.guid()
        data["workable_bounds"] = {
            "location": {
                "x": reader.double(),
                "y": reader.double(),
                "z": reader.double(),
            },
            "rotation": {
                "x": reader.double(),
                "y": reader.double(),
                "z": reader.double(),
                "w": reader.double(),
            },
            "box_sphere_bounds": {
                "origin": {
                    "x": reader.double(),
                    "y": reader.double(),
                    "z": reader.double(),
                },
                "box_extent": {
                    "x": reader.double(),
                    "y": reader.double(),
                    "z": reader.double(),
                },
                "sphere_radius": reader.double(),
            },
        }
        data["base_camp_id_belong_to"] = reader.guid()
        data["owner_map_object_model_id"] = reader.guid()
        data["owner_map_object_concrete_model_id"] = reader.guid()
        data["current_state"] = reader.byte()
        data["assign_locations"] = reader.tarray(
            lambda r: {
                "location": {
                    "x": r.double(),
                    "y": r.double(),
                    "z": r.double(),
                },
                "facing_direction": {
                    "x": r.double(),
                    "y": r.double(),
                    "z": r.double(),
                },
            }
        )
        data["behaviour_type"] = reader.byte()
        data["assign_define_data_id"] = reader.fstring()
        data["override_work_type"] = reader.byte()
        data["assignable_fixed_type"] = reader.byte()
        data["assignable_otomo"] = reader.u32() > 0
        data["can_trigger_worker_event"] = reader.u32() > 0
        data["can_steal_assign"] = reader.u32() > 0
        if work_type == "EPalWorkableType::Defense":
            data["defense_combat_type"] = reader.byte()
        elif work_type == "EPalWorkableType::Progress":
            data["required_work_amount"] = reader.float()
            data["work_exp"] = reader.i32()
            data["current_work_amount"] = reader.float()
            data["auto_work_self_amount_by_sec"] = reader.float()
        elif work_type == "EPalWorkableType::ReviveCharacter":
            data["target_individual_id"] = {
                "player_uid": reader.guid(),
                "instance_id": reader.guid(),
            }
    # These two do not serialize base data
    elif work_type in ["EPalWorkableType::Assign", "EPalWorkableType::LevelObject"]:
        data["handle_id"] = reader.guid()
        data["location_index"] = reader.i32()
        data["assign_type"] = reader.byte()
        data["assigned_individual_id"] = {
            "player_uid": reader.guid(),
            "instance_id": reader.guid(),
        }
        data["state"] = reader.byte()
        data["fixed"] = reader.u32()
        if work_type == "EPalWorkableType::LevelObject":
            data["target_map_object_model_id"] = reader.guid()

    if len(data.keys()) == 0:
        print(f"Warning, unable to parse {work_type}, falling back to raw bytes")
        return {"values": b_bytes}
    # UPalWorkProgressTransformBase->SerializeProperties
    transform_type = reader.byte()
    data["transform"] = {"type": transform_type}
    if transform_type == 1:
        data["transform"]["location"] = {
            "x": reader.double(),
            "y": reader.double(),
            "z": reader.double(),
        }
        data["transform"]["rotation"] = {
            "x": reader.double(),
            "y": reader.double(),
            "z": reader.double(),
            "w": reader.double(),
        }
        data["transform"]["scale"] = {
            "x": reader.double(),
            "y": reader.double(),
            "z": reader.double(),
        }
    elif transform_type == 2:
        data["transform"]["map_object_instance_id"] = reader.guid()
    elif transform_type == 3:
        data["transform"]["guid"] = reader.guid()
        data["transform"]["instance_id"] = reader.guid()
    else:
        remaining_data = reader.read_to_end()
        print(
            f"Unknown EPalWorkTransformType, please report this: {transform_type}: {work_type}: {''.join(f'{b:02x}' for b in remaining_data)}"
        )
        data["transform"]["raw_data"] = [b for b in remaining_data]

    if not reader.eof():
        raise Exception(
            f"Warning: EOF not reached for {work_type}, remaining bytes: {reader.read_to_end()}"
        )

    return data


def decode_work_assign_bytes(b_bytes: Sequence[int]) -> dict[str, Any]:
    reader = FArchiveReader(bytes(b_bytes), debug=False)
    data = {}

    data["id"] = reader.guid()
    data["location_index"] = reader.i32()
    data["assign_type"] = reader.byte()
    data["assigned_individual_id"] = {
        "player_uid": reader.guid(),
        "instance_id": reader.guid(),
    }
    data["state"] = reader.byte()
    data["fixed"] = reader.u32() > 0

    if not reader.eof():
        raise Exception("Warning: EOF not reached")

    return data


def encode(
    writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
) -> int:
    if property_type != "ArrayProperty":
        raise Exception(f"Expected ArrayProperty, got {property_type}")
    del properties["custom_type"]
    for work_element in properties["value"]["values"]:
        work_type = work_element["WorkableType"]["value"]["value"]
        work_element["RawData"]["value"] = {
            "values": [
                b for b in encode_bytes(work_element["RawData"]["value"], work_type)
            ]
        }
        for work_assign in work_element["WorkAssignMap"]["value"]:
            work_assign["value"]["RawData"]["value"] = {
                "values": [
                    b
                    for b in encode_work_assign_bytes(
                        work_assign["value"]["RawData"]["value"]
                    )
                ]
            }
    return writer.property_inner(property_type, properties)


def encode_bytes(p: dict[str, Any], work_type: str) -> bytes:
    writer = FArchiveWriter()

    # Handle base serialization
    if work_type in WORK_BASE_TYPES:
        writer.guid(p["id"])
        writer.double(p["workable_bounds"]["location"]["x"])
        writer.double(p["workable_bounds"]["location"]["y"])
        writer.double(p["workable_bounds"]["location"]["z"])
        writer.double(p["workable_bounds"]["rotation"]["x"])
        writer.double(p["workable_bounds"]["rotation"]["y"])
        writer.double(p["workable_bounds"]["rotation"]["z"])
        writer.double(p["workable_bounds"]["rotation"]["w"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["origin"]["x"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["origin"]["y"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["origin"]["z"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["box_extent"]["x"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["box_extent"]["y"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["box_extent"]["z"])
        writer.double(p["workable_bounds"]["box_sphere_bounds"]["sphere_radius"])
        writer.guid(p["base_camp_id_belong_to"])
        writer.guid(p["owner_map_object_model_id"])
        writer.guid(p["owner_map_object_concrete_model_id"])
        writer.byte(p["current_state"])
        writer.tarray(
            lambda w, l: (
                w.double(l["location"]["x"]),
                w.double(l["location"]["y"]),
                w.double(l["location"]["z"]),
                w.double(l["facing_direction"]["x"]),
                w.double(l["facing_direction"]["y"]),
                w.double(l["facing_direction"]["z"]),
            ),
            p["assign_locations"],
        )
        writer.byte(p["behaviour_type"])
        writer.fstring(p["assign_define_data_id"])
        writer.byte(p["override_work_type"])
        writer.byte(p["assignable_fixed_type"])
        writer.u32(1 if p["assignable_otomo"] else 0)
        writer.u32(1 if p["can_trigger_worker_event"] else 0)
        writer.u32(1 if p["can_steal_assign"] else 0)
        if work_type == "EPalWorkableType::Defense":
            writer.byte(p["defense_combat_type"])
        elif work_type == "EPalWorkableType::Progress":
            writer.float(p["required_work_amount"])
            writer.i32(p["work_exp"])
            writer.float(p["current_work_amount"])
            writer.float(p["auto_work_self_amount_by_sec"])
        elif work_type == "EPalWorkableType::ReviveCharacter":
            writer.guid(p["target_individual_id"]["player_uid"])
            writer.guid(p["target_individual_id"]["instance_id"])
    # These two do not serialize base data
    elif work_type in ["EPalWorkableType::Assign", "EPalWorkableType::LevelObject"]:
        writer.guid(p["handle_id"])
        writer.i32(p["location_index"])
        writer.byte(p["assign_type"])
        writer.guid(p["assigned_individual_id"]["player_uid"])
        writer.guid(p["assigned_individual_id"]["instance_id"])
        writer.byte(p["state"])
        writer.u32(p["fixed"])
        if work_type == "EPalWorkableType::LevelObject":
            writer.guid(p["target_map_object_model_id"])

    # UPalWorkProgressTransformBase->SerializeProperties
    transform_type = p["transform"]["type"]
    writer.byte(transform_type)
    if transform_type == 1:
        writer.double(p["transform"]["location"]["x"])
        writer.double(p["transform"]["location"]["y"])
        writer.double(p["transform"]["location"]["z"])
        writer.double(p["transform"]["rotation"]["x"])
        writer.double(p["transform"]["rotation"]["y"])
        writer.double(p["transform"]["rotation"]["z"])
        writer.double(p["transform"]["rotation"]["w"])
        writer.double(p["transform"]["scale"]["x"])
        writer.double(p["transform"]["scale"]["y"])
        writer.double(p["transform"]["scale"]["z"])
    elif transform_type == 2:
        writer.guid(p["transform"]["map_object_instance_id"])
    elif transform_type == 3:
        writer.guid(p["transform"]["guid"])
        writer.guid(p["transform"]["instance_id"])
    else:
        print(
            f"Unknown EPalWorkTransformType, please report this: {transform_type}: {work_type}"
        )
        writer.bytes(p["transform"]["raw_data"])

    encoded_bytes = writer.bytes()
    return encoded_bytes


def encode_work_assign_bytes(p: dict[str, Any]) -> bytes:
    writer = FArchiveWriter()

    writer.guid(p["id"])
    writer.i32(p["location_index"])
    writer.byte(p["assign_type"])
    writer.guid(p["assigned_individual_id"]["player_uid"])
    writer.guid(p["assigned_individual_id"]["instance_id"])
    writer.byte(p["state"])
    writer.u32(1 if p["fixed"] else 0)

    encoded_bytes = writer.bytes()
    return encoded_bytes
