# SPDX-License-Identifier: Apache-2.0
# Derived from zaigie/palworld-server-tool sav_cli @ fb45624 (Apache-2.0).
# Runtime deps (palsav-flex/palooz/ooz) are GPL-3.0-or-later, so a Docker image
# built from the root Dockerfile includes these runtime components.
"""Decode a Palworld 1.0 save and structure it into player / guild JSON.

It uses the ``palsav`` parser from PalworldSaveTools, which ships Palworld 1.0
mappings (GroupSaveDataMap / character / item-container decoders) plus Oodle
(``PlM1``) decompression via the native ``palooz`` module.

The parser performs a full decode. A ~260KB compressed / ~4MB decompressed
Level.sav completes in a couple of seconds on the validated fixtures.
"""

import os

from palsav.core import decompress_sav_to_gvas
from palsav.gvas import GvasFile
from palsav.paltypes import PALWORLD_TYPE_HINTS, PALWORLD_CUSTOM_PROPERTIES

from world_types import Player, Pal, Guild, BaseCamp
from logger import log

# Global state shared by the current decode helpers.
wsd = None
gvas_file = None

PLAYER_CONTAINER_KEYS = [
    "CommonContainerId",
    "DropSlotContainerId",
    "EssentialContainerId",
    "FoodEquipContainerId",
    "PlayerEquipArmorContainerId",
    "WeaponLoadOutContainerId",
]


def _read_gvas(path):
    with open(path, "rb") as f:
        raw_gvas, _ = decompress_sav_to_gvas(f.read())
    return GvasFile.read(raw_gvas, PALWORLD_TYPE_HINTS, PALWORLD_CUSTOM_PROPERTIES)


def convert_sav(file):
    """Decode Level.sav into the module-global ``wsd`` (worldSaveData)."""
    global gvas_file, wsd
    gvas_file = _read_gvas(file)
    wsd = gvas_file.properties["worldSaveData"]["value"]
    return wsd


def _save_parameter(character_entry):
    return character_entry["value"]["RawData"]["value"]["object"]["SaveParameter"][
        "value"
    ]


def structure_player(dir_path, filetime: int = -1):
    if not wsd.get("CharacterSaveParameterMap"):
        return []

    ticks = wsd["GameTimeSaveData"]["value"]["RealDateTimeTicks"]["value"]
    item_containers = _index_item_containers()

    players = []
    pals = []
    for c in wsd["CharacterSaveParameterMap"]["value"]:
        uid = c["key"]["PlayerUId"]["value"]
        sp = _save_parameter(c)
        if sp.get("IsPlayer") and sp["IsPlayer"]["value"]:
            sp["Items"] = getPlayerItems(uid, dir_path, item_containers)
            players.append(Player(uid, sp).to_dict())
        else:
            if not sp.get("OwnerPlayerUId"):
                continue
            pals.append(Pal(sp, ticks, filetime).to_dict())

    # De-dup players by uid, keeping the highest-level record.
    unique_players_dict = {}
    for player in players:
        pid = player["player_uid"]
        if pid not in unique_players_dict or player["level"] > unique_players_dict[pid]["level"]:
            unique_players_dict[pid] = player
    unique_players = list(unique_players_dict.values())

    for pal in pals:
        for player in unique_players:
            if player["player_uid"] == pal["owner"]:
                pal.pop("owner")
                player["pals"].append(pal)
                break

    return sorted(unique_players, key=lambda p: p["level"], reverse=True)


def _index_item_containers():
    """Map container-UUID string -> decoded slots list."""
    index = {}
    if not wsd.get("ItemContainerSaveData"):
        return index
    for container in wsd["ItemContainerSaveData"]["value"]:
        cid = str(container["key"]["ID"]["value"])
        index[cid] = container["value"]["Slots"]["value"]["values"]
    return index


def getPlayerItems(player_uid, dir_path, item_containers):
    containers_data = {k: [] for k in PLAYER_CONTAINER_KEYS}

    player_sav_file = os.path.join(
        dir_path, str(player_uid).upper().replace("-", "") + ".sav"
    )
    if not os.path.exists(player_sav_file):
        return containers_data

    try:
        player_gvas = _read_gvas(player_sav_file).properties["SaveData"]["value"]
    except Exception as e:
        log(
            f"Player Sav file is corrupted: {os.path.basename(player_sav_file)}: {e}",
            "ERROR",
        )
        return containers_data

    inv = player_gvas.get("InventoryInfo")
    if inv is None:
        return containers_data

    for key in PLAYER_CONTAINER_KEYS:
        ref = inv["value"].get(key)
        if ref is None:
            continue
        container_id = str(ref["value"]["ID"]["value"])
        slots = item_containers.get(container_id)
        if slots is None:
            continue
        items = []
        for slot in slots:
            raw = slot["RawData"]["value"]
            if not raw:  # empty slot decodes to None
                continue
            static_id = raw["item"]["static_id"]
            if not static_id or static_id.lower() == "none":
                continue
            items.append(
                {
                    "SlotIndex": raw["slot_index"],
                    "ItemId": static_id.lower(),
                    "StackCount": raw["count"],
                }
            )
        containers_data[key] = items
    return containers_data


def structure_base_camp():
    if not wsd.get("BaseCampSaveData"):
        return []
    return [
        BaseCamp(b["value"]["RawData"]["value"]).to_dict()
        for b in wsd["BaseCampSaveData"]["value"]
    ]


def structure_guild(filetime: int = -1):
    if not wsd.get("GroupSaveDataMap"):
        return []
    base_camps = structure_base_camp()
    ticks = wsd["GameTimeSaveData"]["value"]["RealDateTimeTicks"]["value"]
    groups = (
        g["value"]["RawData"]["value"]
        for g in wsd["GroupSaveDataMap"]["value"]
        if g["value"]["GroupType"]["value"]["value"] == "EPalGroupType::Guild"
    )
    sorted_guilds = sorted(
        (Guild(g, ticks, filetime).to_dict() for g in groups),
        key=lambda g: g["base_camp_level"],
        reverse=True,
    )
    for guild in sorted_guilds:
        for camp in base_camps:
            if camp["id"] in guild["base_ids"]:
                guild["base_camp"].append(
                    {
                        "id": camp["id"],
                        "area": camp["area_range"],
                        "location_x": camp["transform"]["x"],
                        "location_y": camp["transform"]["y"],
                    }
                )
    return list(sorted_guilds)
