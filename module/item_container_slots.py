import base64

from palworld_save_tools.archive import *


def decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    if type_name != "ArrayProperty":
        raise Exception(f"Expected ArrayProperty, got {type_name}")
    value = reader.property(type_name, size, path, nested_caller_path=path)
    data_bytes = value["value"]["values"]
    value["value"] = decode_bytes(reader, data_bytes)
    return value


def decode_bytes(
    parent_reader: FArchiveReader, c_bytes: Sequence[int]
) -> Optional[dict[str, Any]]:
    if len(c_bytes) == 0:
        return None
    reader = parent_reader.internal_copy(bytes(c_bytes), debug=False)
    data: dict[str, Any] = {
        "permission": {
            "type_a": reader.u32(),
            "type_b": reader.u32(),
            "item_static_id": reader.fstring(),
        },
        "corruption_progress_value": reader.float(),
    }
    unknown_bytes = reader.read_to_end()
    try:
        uuid_bytes = unknown_bytes[12:28]
        local_id = UUID(uuid_bytes)
        data["local_id"] = local_id
    except ValueError:
        data["unknown_padding"] = base64.b64encode(unknown_bytes).decode()

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
    writer.u32(p["permission"]["type_a"])
    writer.u32(p["permission"]["type_b"])
    writer.fstring(p["permission"]["item_static_id"])
    writer.float(p["corruption_progress_value"])
    writer.write(base64.b64decode(p["unknown_padding"]))
    encoded_bytes = writer.bytes()
    return encoded_bytes
