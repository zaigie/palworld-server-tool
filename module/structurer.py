import sys
import zlib
import json
import ijson

# https://github.com/cheahjs/palworld-save-tools/tree/main/lib
from lib.gvas import GvasFile
from lib.json_tools import CustomEncoder
from lib.palsav import decompress_sav_to_gvas
from lib.paltypes import PALWORLD_CUSTOM_PROPERTIES, PALWORLD_TYPE_HINTS

from world_types import Player, Pal, Guild
from logger import log


def convert_sav(file):
    if file.endswith(".sav.json"):
        log("Loading...")
        with open(file, "r", encoding="utf-8") as f:
            return f.read()
    log("Converting...")
    try:
        with open(file, "rb") as f:
            data = f.read()
            raw_gvas, _ = decompress_sav_to_gvas(data)
        gvas_file = GvasFile.read(
            raw_gvas, PALWORLD_TYPE_HINTS, PALWORLD_CUSTOM_PROPERTIES
        )
    except zlib.error:
        log("This .sav file is corrupted. :(", "ERROR")
        sys.exit(1)
    return json.dumps(gvas_file.dump(), cls=CustomEncoder)


def structure_player(converted):
    log("Structuring players...")
    uid_character = (
        (
            c["key"]["PlayerUId"]["value"],
            c["value"]["RawData"]["value"]["object"]["SaveParameter"]["value"],
        )
        for c in ijson.items(
            converted,
            "properties.worldSaveData.value.CharacterSaveParameterMap.value.item",
        )
    )
    players = []
    pals = []
    for uid, c in uid_character:
        if c.get("IsPlayer") and c["IsPlayer"]["value"]:
            players.append(Player(uid, c).to_dict())
        else:
            pals.append(Pal(c).to_dict())
    for pal in pals:
        for player in players:
            if player["player_uid"] == pal["owner"]:
                pal.pop("owner")
                player["pals"].append(pal)
                break
    sorted_players = sorted(players, key=lambda p: p["level"], reverse=True)
    return list(sorted_players)


def structure_guild(converted):
    log("Structuring guilds...")
    groups = (
        g["value"]["RawData"]["value"]
        for g in ijson.items(
            converted, "properties.worldSaveData.value.GroupSaveDataMap.value.item"
        )
        if g["value"]["GroupType"]["value"]["value"] == "EPalGroupType::Guild"
    )
    guilds_generator = (Guild(g).to_dict() for g in groups)
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
