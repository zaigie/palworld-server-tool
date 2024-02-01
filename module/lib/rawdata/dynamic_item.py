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
    buf = bytes(c_bytes)
    reader = FArchiveReader(buf, debug=False)
    data = {}
    data["id"] = {
        "created_world_id": reader.guid(),
        "local_id_in_created_world": reader.guid(),
        "static_id": reader.fstring(),
    }
    data["type"] = "unknown"
    egg_data = try_read_egg(reader)
    if egg_data != None:
        data |= egg_data
    elif (reader.size - reader.data.tell()) == 4:
        data["type"] = "armor"
        data["durability"] = reader.float()
        if not reader.eof():
            raise Exception("Warning: EOF not reached")
    else:
        cur_pos = reader.data.tell()
        temp_data = {"type": "weapon"}
        try:
            temp_data["durability"] = reader.float()
            temp_data["remaining_bullets"] = reader.i32()
            temp_data["passive_skill_list"] = reader.tarray(lambda r: r.fstring())
            if not reader.eof():
                raise Exception("Warning: EOF not reached")
            data |= temp_data
        except Exception as e:
            print(
                f"Warning: Failed to parse weapon data, continuing as raw data {buf}: {e}"
            )
            reader.data.seek(cur_pos)
            data["trailer"] = [int(b) for b in reader.read_to_end()]
    return data


def try_read_egg(reader: FArchiveReader) -> Optional[dict[str, Any]]:
    cur_pos = reader.data.tell()
    try:
        data = {"type": "egg"}
        data["character_id"] = reader.fstring()
        data["object"] = reader.properties_until_end()
        data["unknown_bytes"] = reader.byte_list(4)
        data["unknown_id"] = reader.guid()
        if not reader.eof():
            raise Exception("Warning: EOF not reached")
        return data
    except Exception as e:
        if e.args[0] == "Warning: EOF not reached":
            raise e
        reader.data.seek(cur_pos)
        return None


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
    writer.guid(p["id"]["created_world_id"])
    writer.guid(p["id"]["local_id_in_created_world"])
    writer.fstring(p["id"]["static_id"])
    if p["type"] == "unknown":
        writer.write(bytes(p["trailer"]))
    elif p["type"] == "egg":
        writer.fstring(p["character_id"])
        writer.properties(p["object"])
        writer.write(bytes(p["unknown_bytes"]))
        writer.guid(p["unknown_id"])
    elif p["type"] == "armor":
        writer.float(p["durability"])
    elif p["type"] == "weapon":
        writer.float(p["durability"])
        writer.i32(p["remaining_bullets"])
        writer.tarray(lambda w, d: w.fstring(d), p["passive_skill_list"])
    encoded_bytes = writer.bytes()
    return encoded_bytes
