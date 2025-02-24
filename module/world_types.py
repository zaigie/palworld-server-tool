import datetime
from uuid import UUID


def hexuid_to_decimal(uuid):
    if not isinstance(uuid, str) and not isinstance(uuid, UUID):
        uuid = str(uuid)
    if isinstance(uuid, str):
        hex_part = uuid.split("-")[0]
        decimal_number = int(hex_part, 16)
        return str(decimal_number)
    elif isinstance(uuid, UUID):
        return str(uuid.int)


def tick2local(tick, real_date_time_ticks, filetime):
    ts = filetime + (tick - real_date_time_ticks) / 1e7
    # to RFC3339 like 2006-01-02T15:04:05Z07:00
    t = datetime.datetime.fromtimestamp(ts, tz=datetime.timezone.utc)
    return t.strftime("%Y-%m-%dT%H:%M:%SZ%z").replace("+0000", "")


class Player:
    def __init__(self, uid, data):
        self.player_uid = hexuid_to_decimal(uid)
        self.nickname = data["NickName"]["value"] if data.get("NickName") else ""
        self.level = int(data["Level"]["value"]["value"]) if data.get("Level") else 1
        self.exp = int(data["Exp"]["value"]) if data.get("Exp") else 0
        self.hp = int(data["HP"]["value"]["Value"]["value"]) if data.get("HP") else 0
        self.max_hp = (
            int(data["MaxHP"]["value"]["Value"]["value"]) if data.get("MaxHP") else 0
        )
        self.shield_hp = (
            int(data["ShieldHP"]["value"]["Value"]["value"])
            if data.get("ShieldHP")
            else 0
        )
        self.shield_max_hp = (
            int(data["ShieldMaxHP"]["value"]["Value"]["value"])
            if data.get("ShieldMaxHP")
            else 0
        )
        self.max_status_point = (
            int(data["MaxSP"]["value"]["Value"]["value"]) if data.get("MaxSP") else 0
        )
        self.status_point = {
            s["StatusName"]["value"]: s["StatusPoint"]["value"]
            for s in data["GotStatusPointList"]["value"]["values"]
        }
        full_stomach = (
            float(data["FullStomach"]["value"]) if data.get("FullStomach") else 0
        )
        self.full_stomach = round(full_stomach, 2)
        self.pals = []
        self.items = (
            data["Items"]
            if data["Items"] is not None
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
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }


class Pal:
    def __init__(self, data, real_date_time_ticks, filetime):
        self.owner = hexuid_to_decimal(data["OwnerPlayerUId"]["value"])
        self.nickname = data["NickName"]["value"] if data.get("NickName") else ""
        self.level = int(data["Level"]["value"]["value"]) if data.get("Level") else 1
        self.exp = int(data["Exp"]["value"]) if data.get("Exp") else 0
        self.hp = int(data["HP"]["value"]["Value"]["value"]) if data.get("HP") else 0
        self.max_hp = (
            int(data["MaxHP"]["value"]["Value"]["value"]) if data.get("MaxHP") else 0
        )
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
        self.workspeed = data["CraftSpeed"]["value"] if data.get("CraftSpeed") else 0
        self.melee = (
            int(data["Talent_Melee"]["value"]) if data.get("Talent_Melee") else 0
        )
        self.ranged = (
            int(data["Talent_Shot"]["value"]["value"]) if data.get("Talent_Shot") else 0
        )
        self.defense = (
            int(data["Talent_Defense"]["value"]["value"])
            if data.get("Talent_Defense")
            else 0
        )
        self.rank = int(data["Rank"]["value"]["value"]) if data.get("Rank") else 1
        self.rank_attack = (
            int(data["Rank_Attack"]["value"]["value"]) if data.get("Rank_Attack") else 0
        )
        self.rank_defence = (
            int(data["Rank_Defence"]["value"]["value"])
            if data.get("Rank_Defence")
            else 0
        )
        self.rank_craftspeed = (
            int(data["Rank_CraftSpeed"]["value"]["value"])
            if data.get("Rank_CraftSpeed")
            else 0
        )

        # self.owned_time = (
        #     tick2local(
        #         data["OwnedTime"]["value"],
        #         real_date_time_ticks,
        #         filetime,
        #     )
        #     if data.get("OwnedTime")
        #     else ""
        # )
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
            # "owned_time",
            "skills",
        ]

    def to_dict(self):
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }


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
            for player in data["players"]
        ]
        self.base_ids = [hexuid_to_decimal(x) for x in data["base_ids"]]
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
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }


class BaseCamp:
    def __init__(self, data):
        self.id = hexuid_to_decimal(data["id"])
        # self.name = data["name"]
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
        # self.fast_travel_local_transform = {
        #     "x": data["fast_travel_local_transform"]["translation"]["x"],
        #     "y": data["fast_travel_local_transform"]["translation"]["y"],
        #     "z": data["fast_travel_local_transform"]["translation"]["z"],
        #     "rotation": {
        #         "x": data["fast_travel_local_transform"]["rotation"]["x"],
        #         "y": data["fast_travel_local_transform"]["rotation"]["y"],
        #         "z": data["fast_travel_local_transform"]["rotation"]["z"],
        #         "w": data["fast_travel_local_transform"]["rotation"]["w"],
        #     },
        # }
        self.owner_map_object_instance_id = hexuid_to_decimal(
            data["owner_map_object_instance_id"]
        )

        self.__order = [
            "id",
            # "name",
            "state",
            "transform",
            "area_range",
            "group_id_belong_to",
            # "fast_travel_local_transform",
            "owner_map_object_instance_id",
        ]

    def to_dict(self):
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }
