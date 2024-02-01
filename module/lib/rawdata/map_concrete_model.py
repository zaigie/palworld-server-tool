from typing import Any, Sequence

from lib.archive import *

# def decode_map_concrete_model(
#     reader: FArchiveReader, type_name: str, size: int, path: str
# ) -> dict[str, Any]:
#     if type_name != "StructProperty":
#         raise Exception(f"Expected StructProperty, got {type_name}")
#     value = reader.property(type_name, size, path, nested_caller_path=path)
#     # Decode the raw bytes for the map object and replace the raw data
#     raw_bytes = value["value"]["RawData"]["value"]["values"]
#     print("".join(f"{b:02x}" for b in raw_bytes))
#     # value["value"]["RawData"]["value"] = decode_map_concrete_model_bytes(raw_bytes)
#     # Decode the raw bytes for the module map and replace the raw data
#     # group_map = value["value"]
#     # for group in group_map:
#     #     group_type = group["value"]["GroupType"]["value"]["value"]
#     #     group_bytes = group["value"]["RawData"]["value"]["values"]
#     #     group["value"]["RawData"]["value"] = decode_map_concrete_model_bytes(
#     #         group_bytes, group_type
#     #     )
#     # EPalMapObjectConcreteModelModuleType::None = 0,
#     # EPalMapObjectConcreteModelModuleType::ItemContainer = 1,
#     # EPalMapObjectConcreteModelModuleType::CharacterContainer = 2,
#     # EPalMapObjectConcreteModelModuleType::Workee = 3,
#     # EPalMapObjectConcreteModelModuleType::Energy = 4,
#     # EPalMapObjectConcreteModelModuleType::StatusObserver = 5,
#     # EPalMapObjectConcreteModelModuleType::ItemStack = 6,
#     # EPalMapObjectConcreteModelModuleType::Switch = 7,
#     # EPalMapObjectConcreteModelModuleType::PlayerRecord = 8,
#     # EPalMapObjectConcreteModelModuleType::BaseCampPassiveEffect = 9,
#     # EPalMapObjectConcreteModelModuleType::PasswordLock = 10,
#     return value


# def decode_map_concrete_model_bytes(m_bytes: Sequence[int]) -> dict[str, Any]:
#     if len(m_bytes) == 0:
#         return None
#     reader = FArchiveReader(bytes(m_bytes), debug=False)
#     map_concrete_model = {}

#     if not reader.eof():
#         raise Exception("Warning: EOF not reached")
#     return map_concrete_model


# def encode_map_concrete_model(
#     writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
# ) -> int:
#     if property_type != "MapProperty":
#         raise Exception(f"Expected MapProperty, got {property_type}")
#     del properties["custom_type"]
#     # encoded_bytes = encode_map_concrete_model_bytes(properties["value"]["RawData"]["value"])
#     # properties["value"]["RawData"]["value"] = {"values": [b for b in encoded_bytes]}
#     return writer.property_inner(property_type, properties)


# def encode_map_concrete_model_bytes(p: dict[str, Any]) -> bytes:
#     writer = FArchiveWriter()

#     encoded_bytes = writer.bytes()
#     return encoded_bytes
