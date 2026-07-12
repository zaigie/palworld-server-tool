# SPDX-License-Identifier: Apache-2.0
# Derived from zaigie/palworld-server-tool sav_cli @ fb45624 (Apache-2.0).
# Runtime deps (palsav-flex/palooz/ooz) are GPL-3.0-or-later, so a Docker image
# built from docker/Dockerfile.oss is a GPL-3.0-or-later combined work.
"""Structured player / pal / guild / base-camp views over a decoded Palworld
1.0 save.

These classes take the *fully decoded* palsav property trees (see
``structurer.py``) and flatten them into the JSON shape that
palworld-server-tool's backend consumes (``/player`` and ``/guild`` PUTs).

Field access paths were verified against a real v1.0.0.100427 save. Notable
1.0 changes from the old (closed-source) sav_cli:

* HP is stored under ``Hp`` (was ``HP``), as a ``FixedPoint64`` struct whose
  integer lives at ``Hp.value.Value.value``.
* ``Level`` / talent / rank values are ``ByteProperty`` with the number nested
  at ``.value.value``.
* ``MaxHP`` / ``ShieldMaxHP`` / ``MaxSP`` are not persisted in 1.0 player saves
  (they are recomputed at runtime), so they default to 0 here.
"""

import datetime


def hexuid_to_decimal(uid):
    """First 8 hex chars of a UID -> decimal string.

    palsav yields its own ``archive.UUID`` type (not ``uuid.UUID``); stringifying
    always gives the canonical ``xxxxxxxx-....`` form, so taking the first
    hyphen-separated segment works for both palsav UUIDs and plain strings.
    """
    if uid is None:
        return ""
    hex_part = str(uid).split("-")[0]
    return str(int(hex_part, 16))


def _fixed_point(struct):
    """Extract the integer out of a FixedPoint64 HP/shield struct."""
    if not struct:
        return 0
    try:
        return int(struct["value"]["Value"]["value"])
    except (KeyError, TypeError):
        return 0


def _byte_value(prop, default=0):
    """Extract the number from a ByteProperty (value nested one level)."""
    if not prop:
        return default
    v = prop["value"]
    if isinstance(v, dict):
        return int(v["value"])
    return int(v)


def tick2local(tick, real_date_time_ticks, filetime):
    ts = filetime + (tick - real_date_time_ticks) / 1e7
    t = datetime.datetime.fromtimestamp(ts, tz=datetime.timezone.utc)
    return t.strftime("%Y-%m-%dT%H:%M:%SZ%z").replace("+0000", "")


class Player:
    def __init__(self, uid, data):
        self.player_uid = hexuid_to_decimal(uid)
        self.nickname = data["NickName"]["value"] if data.get("NickName") else ""
        self.level = _byte_value(data.get("Level"), 1)
        self.exp = int(data["Exp"]["value"]) if data.get("Exp") else 0
        # HP was renamed HP -> Hp in 1.0; accept either.
        self.hp = _fixed_point(data.get("Hp") or data.get("HP"))
        self.max_hp = _fixed_point(data.get("MaxHP"))
        self.shield_hp = _fixed_point(data.get("ShieldHP"))
        self.shield_max_hp = _fixed_point(data.get("ShieldMaxHP"))
        self.max_status_point = _fixed_point(data.get("MaxSP"))
        self.status_point = {
            s["StatusName"]["value"]: s["StatusPoint"]["value"]
            for s in data["GotStatusPointList"]["value"]["values"]
        } if data.get("GotStatusPointList") else {}
        full_stomach = (
            float(data["FullStomach"]["value"]) if data.get("FullStomach") else 0
        )
        self.full_stomach = round(full_stomach, 2)
        self.pals = []
        self.items = (
            data["Items"]
            if data.get("Items") is not None
            else {
                "CommonContainerId": [],
                "DropSlotContainerId": [],
                "EssentialContainerId": [],
                "FoodEquipContainerId": [],
                "PlayerEquipArmorContainerId": [],
                "WeaponLoadOutContainerId": [],
            }
        )

        self.__order = [
            "player_uid",
            "nickname",
            "level",
            "exp",
            "hp",
            "max_hp",
            "shield_hp",
            "shield_max_hp",
            "max_status_point",
            "status_point",
            "full_stomach",
            "pals",
            "items",
        ]

    def to_dict(self):
        return {attr: getattr(self, attr) for attr in self.__order}


class Pal:
    def __init__(self, data, real_date_time_ticks, filetime):
        self.owner = hexuid_to_decimal(data["OwnerPlayerUId"]["value"])
        self.nickname = data["NickName"]["value"] if data.get("NickName") else ""
        self.level = _byte_value(data.get("Level"), 1)
        self.exp = int(data["Exp"]["value"]) if data.get("Exp") else 0
        self.hp = _fixed_point(data.get("Hp") or data.get("HP"))
        self.max_hp = _fixed_point(data.get("MaxHP"))
        self.gender = (
            data["Gender"]["value"]["value"].split("::")[-1]
            if data.get("Gender")
            else "Unknow"
        )
        self.is_lucky = data["IsRarePal"]["value"] if data.get("IsRarePal") else False
        self.is_boss = False

        if data.get("CharacterID"):
            typename = data["CharacterID"]["value"]
            typename_upper = typename.upper()
            if typename_upper[:5] == "BOSS_":
                typename_upper = typename_upper.replace("BOSS_", "")
                self.is_boss = not self.is_lucky
            self.is_tower = typename_upper.startswith("GYM_")
            self.type = typename
        else:
            self.is_tower = False
            self.type = "Unknow"

        self.workspeed = _byte_value(data.get("CraftSpeed"), 0)
        self.melee = _byte_value(data.get("Talent_Melee"), 0)
        self.ranged = _byte_value(data.get("Talent_Shot"), 0)
        self.defense = _byte_value(data.get("Talent_Defense"), 0)
        self.rank = _byte_value(data.get("Rank"), 1)
        self.rank_attack = _byte_value(data.get("Rank_Attack"), 0)
        self.rank_defence = _byte_value(data.get("Rank_Defence"), 0)
        self.rank_craftspeed = _byte_value(data.get("Rank_CraftSpeed"), 0)
        self.skills = (
            data["PassiveSkillList"]["value"]["values"]
            if data.get("PassiveSkillList")
            else []
        )

        self.__order = [
            "owner",
            "nickname",
            "level",
            "exp",
            "hp",
            "max_hp",
            "type",
            "gender",
            "is_lucky",
            "is_boss",
            "is_tower",
            "workspeed",
            "melee",
            "ranged",
            "defense",
            "rank",
            "rank_attack",
            "rank_defence",
            "rank_craftspeed",
            "skills",
        ]

    def to_dict(self):
        return {attr: getattr(self, attr) for attr in self.__order}


class Guild:
    def __init__(self, data, real_date_time_ticks, filetime):
        self.name = data["guild_name"]
        self.base_camp_level = data["base_camp_level"]
        self.admin_player_uid = hexuid_to_decimal(data["admin_player_uid"])
        self.players = [
            {
                "player_uid": hexuid_to_decimal(player["player_uid"]),
                "nickname": player["player_info"]["player_name"],
                "last_online": (
                    tick2local(
                        player["player_info"]["last_online_real_time"],
                        real_date_time_ticks,
                        filetime,
                    )
                    if player["player_info"].get("last_online_real_time")
                    else ""
                ),
            }
            for player in data.get("players", [])
        ]
        self.base_ids = [hexuid_to_decimal(x) for x in data.get("base_ids", [])]
        self.base_camp = []
        self.__order = [
            "name",
            "base_camp_level",
            "admin_player_uid",
            "players",
            "base_ids",
            "base_camp",
        ]

    def to_dict(self):
        return {attr: getattr(self, attr) for attr in self.__order}


class BaseCamp:
    def __init__(self, data):
        self.id = hexuid_to_decimal(data["id"])
        self.state = data["state"]
        self.transform = {
            "x": data["transform"]["translation"]["x"],
            "y": data["transform"]["translation"]["y"],
            "z": data["transform"]["translation"]["z"],
            "rotation": {
                "x": data["transform"]["rotation"]["x"],
                "y": data["transform"]["rotation"]["y"],
                "z": data["transform"]["rotation"]["z"],
                "w": data["transform"]["rotation"]["w"],
            },
        }
        self.area_range = data["area_range"]
        self.group_id_belong_to = hexuid_to_decimal(data["group_id_belong_to"])
        self.owner_map_object_instance_id = hexuid_to_decimal(
            data["owner_map_object_instance_id"]
        )
        self.__order = [
            "id",
            "state",
            "transform",
            "area_range",
            "group_id_belong_to",
            "owner_map_object_instance_id",
        ]

    def to_dict(self):
        return {attr: getattr(self, attr) for attr in self.__order}
