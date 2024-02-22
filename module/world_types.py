from uuid import UUID
from logger import log


def hexuid_to_decimal(uuid):
    if not isinstance(uuid, str) and not isinstance(uuid, UUID):
        uuid = str(uuid)
    if isinstance(uuid, str):
        hex_part = uuid.split("-")[0]
        decimal_number = int(hex_part, 16)
        return str(decimal_number)
    elif isinstance(uuid, UUID):
        return str(uuid.int)


class Player:
    def __init__(self, uid, data):
        self.player_uid = hexuid_to_decimal(uid)
        self.nickname = data["NickName"]["value"] if data.get("NickName") else ""
        self.level = int(data["Level"]["value"]) if data.get("Level") else 1
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
        ]

    def to_dict(self):
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }


class Pal:
    def __init__(self, data):
        self.owner = hexuid_to_decimal(data["OwnerPlayerUId"]["value"])
        # self.nickname = data["Nickname"]["value"] if data.get("Nicknme") else ""
        self.level = int(data["Level"]["value"]) if data.get("Level") else 1
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
            int(data["Talent_Shot"]["value"]) if data.get("Talent_Shot") else 0
        )
        self.defense = (
            int(data["Talent_Defense"]["value"]) if data.get("Talent_Defense") else 0
        )
        self.rank = int(data["Rank"]["value"]) if data.get("Rank") else 1
        self.skills = (
            data["PassiveSkillList"]["value"]["values"]
            if data.get("PassiveSkillList")
            else []
        )

        self.__order = [
            "owner",
            # "nickname",
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
            "skills",
        ]

    def to_dict(self):
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }


class Guild:
    def __init__(self, data):
        self.name = data["guild_name"]
        self.base_camp_level = data["base_camp_level"]
        self.admin_player_uid = hexuid_to_decimal(data["admin_player_uid"])
        self.players = [
            {
                "player_uid": hexuid_to_decimal(player["player_uid"]),
                "nickname": player["player_info"]["player_name"],
            }
            for player in data["players"]
        ]
        self.base_ids = [str(x) for x in data["base_ids"]]
        self.__order = [
            "name",
            "base_camp_level",
            "admin_player_uid",
            "players",
            "base_ids",
        ]

    def to_dict(self):
        return {
            attr: getattr(self, attr)
            for attr in self.__order
            if not attr.startswith("_") and not callable(getattr(self, attr))
        }
