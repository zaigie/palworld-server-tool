import copy
import datetime
import os
import sys
import zlib
import json
import time
from typing import Any, Callable

from palworld_save_tools.gvas import GvasFile
from palworld_save_tools.palsav import decompress_sav_to_gvas
from palworld_save_tools.paltypes import PALWORLD_CUSTOM_PROPERTIES, PALWORLD_TYPE_HINTS
from palworld_save_tools.archive import FArchiveReader, FArchiveWriter
from palworld_save_tools.rawdata import (
    character,
    group,
    item_container,
    item_container_slots,
)

from world_types import Player, Pal, Guild
from logger import log, redirect_stdout_stderr

# PALWORLD_CUSTOM_PROPERTIES: dict[
#     str,
#     tuple[
#         Callable[[FArchiveReader, str, int, str], dict[str, Any]],
#         Callable[[FArchiveWriter, str, dict[str, Any]], int],
#     ],
# ] = {
#     ".worldSaveData.GroupSaveDataMap": (group.decode, group.encode),
#     ".worldSaveData.CharacterSaveParameterMap.Value.RawData": (
#         character.decode,
#         character.encode,
#     ),
#     ".worldSaveData.ItemContainerSaveData.Value.RawData": (
#         item_container.decode,
#         item_container.encode,
#     ),
# }
wsd = None
gvas_file = None


def skip_decode(
    reader: FArchiveReader, type_name: str, size: int, path: str
) -> dict[str, Any]:
    if type_name == "ArrayProperty":
        array_type = reader.fstring()
        value = {
            "skip_type": type_name,
            "array_type": array_type,
            "id": reader.optional_guid(),
            "value": reader.read(size),
        }
    elif type_name == "MapProperty":
        key_type = reader.fstring()
        value_type = reader.fstring()
        _id = reader.optional_guid()
        value = {
            "skip_type": type_name,
            "key_type": key_type,
            "value_type": value_type,
            "id": _id,
            "value": reader.read(size),
        }
    elif type_name == "StructProperty":
        value = {
            "skip_type": type_name,
            "struct_type": reader.fstring(),
            "struct_id": reader.guid(),
            "id": reader.optional_guid(),
            "value": reader.read(size),
        }
    else:
        raise Exception(
            f"Expected ArrayProperty or MapProperty or StructProperty, got {type_name} in {path}"
        )
    return value


def skip_encode(
    writer: FArchiveWriter, property_type: str, properties: dict[str, Any]
) -> int:
    if "skip_type" not in properties:
        if properties["custom_type"] in PALWORLD_CUSTOM_PROPERTIES is not None:
            return PALWORLD_CUSTOM_PROPERTIES[properties["custom_type"]][1](
                writer, property_type, properties
            )
    if property_type == "ArrayProperty":
        del properties["custom_type"]
        del properties["skip_type"]
        writer.fstring(properties["array_type"])
        writer.optional_guid(properties.get("id", None))
        writer.write(properties["value"])
        return len(properties["value"])
    elif property_type == "MapProperty":
        del properties["custom_type"]
        del properties["skip_type"]
        writer.fstring(properties["key_type"])
        writer.fstring(properties["value_type"])
        writer.optional_guid(properties.get("id", None))
        writer.write(properties["value"])
        return len(properties["value"])
    elif property_type == "StructProperty":
        del properties["custom_type"]
        del properties["skip_type"]
        writer.fstring(properties["struct_type"])
        writer.guid(properties["struct_id"])
        writer.optional_guid(properties.get("id", None))
        writer.write(properties["value"])
        return len(properties["value"])
    else:
        raise Exception(
            f"Expected ArrayProperty or MapProperty or StructProperty, got {property_type}"
        )


def load_skiped_decode(wsd, skip_paths, recursive=True):
    if isinstance(skip_paths, str):
        skip_paths = [skip_paths]
    for skip_path in skip_paths:
        properties = wsd[skip_path]

        if "skip_type" not in properties:
            continue
        parse_skiped_item(properties, skip_path, recursive)
        if ".worldSaveData.%s" % skip_path in SKP_PALWORLD_CUSTOM_PROPERTIES:
            del SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.%s" % skip_path]


SKP_PALWORLD_CUSTOM_PROPERTIES = copy.deepcopy(PALWORLD_CUSTOM_PROPERTIES)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.MapObjectSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.FoliageGridSaveDataMap"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.MapObjectSpawnerInStageSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.DynamicItemSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.CharacterContainerSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.CharacterContainerSaveData.Value.Slots"
] = (skip_decode, skip_encode)
SKP_PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.CharacterContainerSaveData.Value.RawData"
] = (skip_decode, skip_encode)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.ItemContainerSaveData"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[
    ".worldSaveData.ItemContainerSaveData.Value.BelongInfo"
] = (skip_decode, skip_encode)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.ItemContainerSaveData.Value.Slots"] = (
    skip_decode,
    skip_encode,
)
SKP_PALWORLD_CUSTOM_PROPERTIES[".worldSaveData.ItemContainerSaveData.Value.RawData"] = (
    skip_decode,
    skip_encode,
)


def convert_sav(file):
    global gvas_file, wsd
    if file.endswith(".sav.json"):
        log("Loading...")
        with open(file, "r", encoding="utf-8") as f:
            return f.read()
    log("Converting...")
    with redirect_stdout_stderr():
        try:
            with open(file, "rb") as f:
                data = f.read()
                raw_gvas, _ = decompress_sav_to_gvas(data)
            gvas_file = GvasFile.read(
                raw_gvas, PALWORLD_TYPE_HINTS, SKP_PALWORLD_CUSTOM_PROPERTIES
            )
        except zlib.error:
            log("This .sav file is corrupted. :(", "ERROR")
            sys.exit(1)
    # return json.dumps(gvas_file.dump(), cls=CustomEncoder)
    wsd = gvas_file.properties["worldSaveData"]["value"]


def structure_player(dir_path, data_source=None):
    log("Structuring players...")
    global wsd
    if data_source is None:
        data_source = wsd
    if not data_source.get("CharacterSaveParameterMap"):
        return []
    uid_character = (
        (
            c["key"]["PlayerUId"]["value"],
            c["value"]["RawData"]["value"]["object"]["SaveParameter"]["value"],
        )
        for c in wsd["CharacterSaveParameterMap"]["value"]
    )

    players = []
    pals = []
    for uid, c in uid_character:
        if c.get("IsPlayer") and c["IsPlayer"]["value"]:
            c["Items"] = getPlayerItems(uid, dir_path)
            players.append(Player(uid, c).to_dict())
        else:
            if not c.get("OwnerPlayerUId"):
                continue
            pals.append(Pal(c).to_dict())

    unique_players_dict = {}
    for player in players:
        player_uid = player["player_uid"]
        if player_uid in unique_players_dict:
            existing_player = unique_players_dict[player_uid]
            if player["level"] > existing_player["level"]:
                unique_players_dict[player_uid] = player
        else:
            unique_players_dict[player_uid] = player

    unique_players = list(unique_players_dict.values())
    for pal in pals:
        for player in unique_players:
            if player["player_uid"] == pal["owner"]:
                pal.pop("owner")
                player["pals"].append(pal)
                break

    sorted_players = sorted(unique_players, key=lambda p: p["level"], reverse=True)

    return sorted_players


def parse_skiped_item(properties, skip_path, recursive=True):
    if "skip_type" not in properties:
        return properties

    with FArchiveReader(
        properties["value"],
        PALWORLD_TYPE_HINTS,
        (
            SKP_PALWORLD_CUSTOM_PROPERTIES
            if recursive == False
            else PALWORLD_CUSTOM_PROPERTIES
        ),
    ) as reader:
        if properties["skip_type"] == "ArrayProperty":
            properties["value"] = reader.array_property(
                properties["array_type"],
                len(properties["value"]) - 4,
                ".worldSaveData.%s" % skip_path,
            )
        elif properties["skip_type"] == "StructProperty":
            properties["value"] = reader.struct_value(
                properties["struct_type"], ".worldSaveData.%s" % skip_path
            )
        elif properties["skip_type"] == "MapProperty":
            reader.u32()
            count = reader.u32()
            path = ".worldSaveData.%s" % skip_path
            key_path = path + ".Key"
            key_type = properties["key_type"]
            value_type = properties["value_type"]
            if key_type == "StructProperty":
                key_struct_type = reader.get_type_or(key_path, "Guid")
            else:
                key_struct_type = None
            value_path = path + ".Value"
            if value_type == "StructProperty":
                value_struct_type = reader.get_type_or(value_path, "StructProperty")
            else:
                value_struct_type = None
            values: list[dict[str, Any]] = []
            for _ in range(count):
                key = reader.prop_value(key_type, key_struct_type, key_path)
                value = reader.prop_value(value_type, value_struct_type, value_path)
                values.append(
                    {
                        "key": key,
                        "value": value,
                    }
                )
            properties["key_struct_type"] = key_struct_type
            properties["value_struct_type"] = value_struct_type
            properties["value"] = values
        del properties["custom_type"]
        del properties["skip_type"]
    return properties


def parse_item(properties, skip_path):
    if isinstance(properties, dict):
        for key in properties:
            call_skip_path = skip_path + "." + key[0].upper() + key[1:]
            if (
                isinstance(properties[key], dict)
                and "type" in properties[key]
                and properties[key]["type"]
                in ["StructProperty", "ArrayProperty", "MapProperty"]
            ):
                if "skip_type" in properties[key]:
                    # print("Parsing worldSaveData.%s..." % call_skip_path, end="", flush=True)
                    properties[key] = parse_skiped_item(
                        properties[key], call_skip_path, True
                    )
                    # print("Done")
                else:
                    properties[key]["value"] = parse_item(
                        properties[key]["value"], call_skip_path
                    )
            else:
                properties[key] = parse_item(properties[key], call_skip_path)
    elif isinstance(properties, list):
        top_skip_path = ".".join(skip_path.split(".")[:-1])
        for idx, item in enumerate(properties):
            properties[idx] = parse_item(item, top_skip_path)
    return properties


def getPlayerItems(player_uid, dir_path):
    load_skiped_decode(wsd, ["ItemContainerSaveData"], False)
    item_containers = {}
    for item_container in wsd["ItemContainerSaveData"]["value"]:
        item_containers[str(item_container["key"]["ID"]["value"])] = item_container

    player_sav_file = os.path.join(
        dir_path, str(player_uid).upper().replace("-", "") + ".sav"
    )
    if not os.path.exists(player_sav_file):
        # log("Player Sav file Not exists: %s" % player_sav_file)
        return
    else:
        with redirect_stdout_stderr():
            try:
                with open(player_sav_file, "rb") as f:
                    raw_gvas, _ = decompress_sav_to_gvas(f.read())
                    player_gvas_file = GvasFile.read(
                        raw_gvas, PALWORLD_TYPE_HINTS, PALWORLD_CUSTOM_PROPERTIES
                    )
                player_gvas = player_gvas_file.properties["SaveData"]["value"]
            except Exception as e:
                log(
                    f"Player Sav file is corrupted: {os.path.basename(player_sav_file)}: {str(e)}",
                    "ERROR",
                )
                return
    containers_data = {
        "CommonContainerId": [],
        "DropSlotContainerId": [],
        "EssentialContainerId": [],
        "FoodEquipContainerId": [],
        "PlayerEquipArmorContainerId": [],
        "WeaponLoadOutContainerId": [],
    }
    for idx_key in containers_data.keys():
        container_id = str(
            player_gvas["inventoryInfo"]["value"][idx_key]["value"]["ID"]["value"]
        )
        if container_id in item_containers:
            # 解析对应的物品容器数据
            item_container = parse_item(
                item_containers[container_id], "ItemContainerSaveData"
            )

            # 提取每个物品的相关数据并保存到字典中
            containers_data[idx_key] = [
                {
                    "SlotIndex": item["SlotIndex"]["value"],
                    "ItemId": item["ItemId"]["value"]["StaticId"]["value"].lower(),
                    "StackCount": item["StackCount"]["value"],
                }
                for item in item_container["value"]["Slots"]["value"]["values"]
                if item["ItemId"]["value"]["StaticId"]["value"].lower() != "none"
            ]
    return containers_data


def structure_guild(filetime: int = -1):
    log("Structuring guilds...")
    if not wsd.get("GroupSaveDataMap"):
        return []
    groups = (
        g["value"]["RawData"]["value"]
        for g in wsd["GroupSaveDataMap"]["value"]
        if g["value"]["GroupType"]["value"]["value"] == "EPalGroupType::Guild"
    )
    Ticks = wsd["GameTimeSaveData"]["value"]["RealDateTimeTicks"]["value"]
    guilds_generator = (Guild(g, Ticks, filetime).to_dict() for g in groups)
    sorted_guilds = sorted(
        guilds_generator, key=lambda g: g["base_camp_level"], reverse=True
    )
    return list(sorted_guilds)


if __name__ == "__main__":
    import time

    start = time.time()
    file = "./Level.sav"
    converted = convert_sav(file)
    players = structure_player(converted)
    log("Saving players...")
    with open("players.json", "w", encoding="utf-8") as f:
        json.dump(players, f, indent=4, ensure_ascii=False)
    guilds = structure_guild(converted)
    log("Saving guilds...")
    with open("guilds.json", "w", encoding="utf-8") as f:
        json.dump(guilds, f, indent=4, ensure_ascii=False)
    log(f"Done in {time.time() - start}s")
