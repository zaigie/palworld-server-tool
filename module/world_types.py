from enum import Enum
from uuid import UUID
from logger import log


def hexuid_to_decimal(uuid):
    if isinstance(uuid, str):
        hex_part = uuid.split("-")[0]
        decimal_number = int(hex_part, 16)
        return str(decimal_number)
    elif isinstance(uuid, UUID):
        return str(uuid.int)


# https://github.com/EternalWraith/PalEdit/blob/main/PalInfo.py
class PalSkills(Enum):
    UNKNOWN = "Unknown"
    NONE = "None"

    ElementResist_Normal_1_PAL = "Abnormal"
    ElementResist_Dark_1_PAL = "Cheery"
    ElementResist_Dragon_1_PAL = "Dragonkiller"
    ElementResist_Ice_1_PAL = "Heated Body"
    ElementResist_Fire_1_PAL = "Suntan Lover"
    ElementResist_Leaf_1_PAL = "Botanical Barrier"
    ElementResist_Earth_1_PAL = "Earthquake Resistant"
    ElementResist_Thunder_1_PAL = "Insulated Body"
    ElementResist_Aqua_1_PAL = "Waterproof"

    ElementBoost_Normal_1_PAL = "Zen Mind"
    ElementBoost_Dark_1_PAL = "Veil of Darkness"
    ElementBoost_Dragon_1_PAL = "Blood of the Dragon"
    ElementBoost_Ice_1_PAL = "Coldblooded"
    ElementBoost_Fire_1_PAL = "Pyromaniac"
    ElementBoost_Leaf_1_PAL = "Fragrant Foliage"
    ElementBoost_Earth_1_PAL = "Power of Gaia"
    ElementBoost_Thunder_1_PAL = "Capacitor"
    ElementBoost_Aqua_1_PAL = "Hydromaniac"

    ElementBoost_Normal_2_PAL = "Celestial Emperor"
    ElementBoost_Dark_2_PAL = "Lord of the Underworld"
    ElementBoost_Dragon_2_PAL = "Divine Dragon"
    ElementBoost_Ice_2_PAL = "Ice Emperor"
    ElementBoost_Fire_2_PAL = "Flame Emperor"
    ElementBoost_Leaf_2_PAL = "Spirit Emperor"
    ElementBoost_Earth_2_PAL = "Earth Emperor"
    ElementBoost_Thunder_2_PAL = "Lord of Lightning"
    ElementBoost_Aqua_2_PAL = "Lord of the Sea"

    PAL_ALLAttack_up1 = "Brave"
    PAL_ALLAttack_up2 = "Ferocious"
    PAL_ALLAttack_down1 = "Coward"
    PAL_ALLAttack_down2 = "Pacifist"

    Deffence_up1 = "Hard Skin"
    Deffence_up2 = "Burly Body"
    Deffence_down1 = "Downtrodden"
    Deffence_down2 = "Brittle"

    TrainerMining_up1 = "Mine Foreman"
    TrainerLogging_up1 = "Logging Foreman"
    TrainerATK_UP_1 = "Vanguard"
    TrainerWorkSpeed_UP_1 = "Motivational Leader"
    TrainerDEF_UP_1 = "Stronghold Strategist"

    PAL_Sanity_Down_1 = "Positive Thinker"
    PAL_Sanity_Down_2 = "Workaholic"
    PAL_Sanity_Up_1 = "Unstable"
    PAL_Sanity_Up_2 = "Destructive"

    PAL_FullStomach_Down_1 = "Dainty Eater"
    PAL_FullStomach_Down_2 = "Diet Lover"
    PAL_FullStomach_Up_1 = "Glutton"
    PAL_FullStomach_Up_2 = "Bottomless Stomach"

    CraftSpeed_up1 = "Serious"
    CraftSpeed_up2 = "Artisan"
    CraftSpeed_down1 = "Clumsy"
    CraftSpeed_down2 = "Slacker"

    MoveSpeed_up_1 = "Nimble"
    MoveSpeed_up_2 = "Runner"
    MoveSpeed_up_3 = "Swift"

    PAL_CorporateSlave = "Work Slave"

    PAL_rude = "Hooligan"
    Noukin = "Musclehead"

    PAL_oraora = "Aggressive"

    PAL_conceited = "Conceited"

    PAL_masochist = "Masochist"
    PAL_sadist = "Sadist"

    Rare = "Lucky"
    Legend = "Legend"


class PalType(Enum):
    Alpaca = "Melpaca"
    AmaterasuWolf = "Kitsun"
    Anubis = "Anubis"
    Baphomet = "Incineram"
    Baphomet_Dark = "Incineram Noct"
    Bastet = "Mau"
    Bastet_Ice = "Mau Cryst"
    BerryGoat = "Caprity"
    BirdDragon = "Vanwyrm"
    BirdDragon_Ice = "Vanwyrm Cryst"
    BlackCentaur = "Necromus"
    BlackGriffon = "Shadowbeak"
    BlackMetalDragon = "Astegon"
    BlueDragon = "Azurobe"
    BluePlatypus = "Fuack"
    Boar = "Rushoar"
    CaptainPenguin = "Penking"
    Carbunclo = "Lifmunk"
    CatBat = "Tombat"
    CatMage = "Katress"
    CatVampire = "Felbat"
    ChickenPal = "Chikipi"
    ColorfulBird = "Tocotoco"
    CowPal = "Mozzarina"
    CuteButterfly = "Cinnamoth"
    CuteFox = "Vixy"
    CuteMole = "Fuddler"
    DarkCrow = "Cawgnito"
    DarkScorpion = "Menasting"
    Deer = "Eikthyrdeer"
    Deer_Ground = "Eikthyrdeer Terra"
    DreamDemon = "Daedream"
    DrillGame = "Digtoise"
    Eagle = "Galeclaw"
    ElecCat = "Sparkit"
    ElecPanda = "Grizzbolt"
    FairyDragon = "Elphidran"
    FairyDragon_Water = "Elphidran Aqua"
    FengyunDeeper = "Fenglope"
    FireKirin = "Pyrin"
    FireKirin_Dark = "Pyrin Noct"
    FlameBambi = "Rooby"
    FlameBuffalo = "Arsox"
    FlowerDinosaur = "Dinossom"
    FlowerDinosaur_Electric = "Dinossom Lux"
    FlowerDoll = "Petallia"
    FlowerRabbit = "Flopie"
    FlyingManta = "Celaray"
    FoxMage = "Wixen"
    Ganesha = "Teafant"
    Garm = "Direhowl"
    GhostBeast = "Maraith"
    Gorilla = "Gorirat"
    GrassMammoth = "Mammorest"
    GrassMammoth_Ice = "Mammorest Cryst"
    GrassPanda = "Mossanda"
    GrassPanda_Electric = "Mossanda Lux"
    GrassRabbitMan = "Verdash"
    HadesBird = "Helzephyr"
    HawkBird = "Nitewing"
    Hedgehog = "Jolthog"
    Hedgehog_Ice = "Jolthog Cryst"
    HerculesBeetle = "Warsect"
    Horus = "Faleris"
    IceDeer = "Reindrix"
    IceFox = "Foxcicle"
    IceHorse = "Frostallion"
    IceHorse_Dark = "Frostallion Noct"
    JetDragon = "Jetragon"
    Kelpie = "Kelpsea"
    Kelpie_Fire = "Kelpsea Ignis"
    KingAlpaca = "Kingpaca"
    KingAlpaca_Ice = "Ice Kingpaca"
    KingBahamut = "Blazamut"
    Kirin = "Univolt"
    Kitsunebi = "Foxparks"
    LavaGirl = "Flambelle"
    LazyCatfish = "Dumud"
    LazyDragon = "Relaxaurus"
    LazyDragon_Electric = "Relaxaurus Lux"
    LilyQueen = "Lyleen"
    LilyQueen_Dark = "Lyleen Noct"
    LittleBriarRose = "Bristla"
    LizardMan = "Leezpunk"
    LizardMan_Fire = "Leezpunk Ignis"
    Manticore = "Blazehowl"
    Manticore_Dark = "Blazehowl Noct"
    Monkey = "Tanzee"
    MopBaby = "Swee"
    MopKing = "Sweepa"
    Mutant = "Lunaris"
    NaughtyCat = "Grintale"
    NegativeKoala = "Depresso"
    NegativeOctopus = "Killamari"
    NightFox = "Nox"
    Penguin = "Pengullet"
    PinkCat = "Cattiva"
    PinkLizard = "Lovander"
    PinkRabbit = "Ribbuny"
    PlantSlime = "Gumoss"
    PlantSlime_Flower = "Gumoss Flora"
    QueenBee = "Elizabee"
    RaijinDaughter = "Dazzi"
    RedArmorBird = "Ragnahawk"
    RobinHood = "Robinquill"
    RobinHood_Ground = "Robinquill Terra"
    Ronin = "Bushi"
    SaintCentaur = "Paladius"
    SakuraSaurus = "Broncherry"
    SakuraSaurus_Water = "Broncherry Aqua"
    Serpent = "Surfent"
    Serpent_Ground = "Surfent Terra"
    SharkKid = "Gobfin"
    SharkKid_Fire = "Gobfin Ignis"
    Sheepball = "Lamball"
    SkyDragon = "Quivern"
    SoldierBee = "Beegarde"
    Suzaku = "Suzaku"
    Suzaku_Water = "Suzaku Aqua"
    SweetsSheep = "Woolipop"
    ThunderBird = "Beakon"
    ThunderDog = "Rayhound"
    ThunderDragonMan = "Orserk"
    Umihebi = "Jormuntide"
    Umihebi_Fire = "Jormuntide Ignis"
    VioletFairy = "Vaelet"
    VolcanicMonster = "Reptyro"
    VolcanicMonster_Ice = "Reptyro Cryst"
    WeaselDragon = "Chillet"
    Werewolf = "Loupmoon"
    WhiteMoth = "Sibelyx"
    WhiteTiger = "Cryolinx"
    WindChimes = "Hangyu"
    WindChimes_Ice = "Hangyu Cryst"
    WizardOwl = "Hoocrates"
    WoolFox = "Cremis"
    Yeti = "Wumpo"
    Yeti_Grass = "Wumpo Botan"

    # Tower Bosses
    GYM_ThunderDragonMan = "Axel & Orserk"
    GYM_LilyQueen = "Lily & Lyleen"
    GYM_Horus = "Marus & Faleris"
    GYM_BlackGriffon = "Victor & Shadowbeak"
    GYM_ElecPanda = "Zoe & Grizzbolt"

    # Human Entities (Not yet finished)
    Male_DarkTrader01 = "Black Marketeer"
    FireCult_FlameThrower = "Brothers of the Eternal Pyre Martyr"
    FireCult_Rifle = "Brothers of the Eternal Pyre Martyr"
    FireCult_FlameThrower_Invader = "Brothers of the Eternal Pyre Martyr"
    FireCult_Rifle_Invader = "Brothers of the Eternal Pyre Martyr"
    Believer_Bat_Invader = "Free Pal Alliance Devout"
    Believer_CrossBow = "Free Pal Alliance Devout"
    Believer_CrossBow_Invader = "Free Pal Alliance Devout"
    Believer_Bat = "Free Pal Alliance Devout"
    Male_Soldier01 = "Burly Merc"
    Female_Soldier01 = "Expedition Survivor"
    Believer_CrossBow = "Free Pal Alliance Devout"
    Male_Scientist01_LaserRifle = "PAL Genetic Research Unit Executioner"
    Male_Scientist01_LaserRifle_Invader = "PAL Genetic Research Unit Executioner"
    Scientist_FlameThrower = "PAL Genetic Research Unit Executioner"
    Scientist_FlameThrower_Invader = "PAL Genetic Research Unit Executioner"
    Police_Shotgun = "PIDF Elite"
    Police_Handgun = "PIDF Guard"
    Police_Rifle = "PIDF Infantry"
    PalDealer = "Pal Merchant"
    PalTrader = "Pal Merchant"
    Police_Handgun = "PIDF Guard"
    Hunter_Bat = "Syndicate Thug (Bat)"
    Hunter_FlameThrower = "Syndicate Cleaner"
    Hunter_Fat_GatlingGun = "Syndicate Crusher"
    Hunter_RocketLauncher = "Syndicate Elite"
    Hunter_Grenade = "Syndicate Grenadier"
    Hunter_Rifle = "Syndicate Gunner"
    Hunter_Shotgun = "Syndicate Hunter"
    Hunter_Handgun = "Syndicate Thug (Handgun)"
    SalesPerson = "Wandering Merchant"
    SalesPerson_Wander = "Wandering Merchant"
    Female_People03 = "Villager"
    Female_People02 = "Villager"
    RandomEventShop = "Wandering Merchant(Random)"
    PalDealer_Desert = "Pal Merchant(Desert)"
    PalDealer_Volcano = "Pal Merchant(Volcano)"
    SalesPerson_Volcano = "Wandering Merchant(Volcano)"
    SalesPerson_Volcano2 = "Wandering Merchant(Volcano)"
    SalesPerson_Desert = "Wandering Merchant(Desert)"
    SalesPerson_Desert2 = "Wandering Merchant(Desert)"
    Hunter_Bat_Invader = "Syndicate Thug (Bat)"
    Hunter_FlameThrower_Invader = "Invader Syndicate Cleaner"
    Hunter_Fat_GatlingGun_Invader = "Invader Syndicate Crusher"
    Hunter_RocketLauncher_Invader = "Invader Syndicate Elite"
    Hunter_Grenade_Invader = "Invader Syndicate Grenadier"
    Hunter_Rifle_Invader = "Invader Syndicate Gunner"
    Hunter_Shotgun_Invader = "Invader Syndicate Hunter"
    Hunter_Handgun_Invader = "Invader Syndicate Thug (Handgun)"
    Visitor_Hunter_Rifle = "Syndicate Thug"
	Male_People03 = "Villager"
	Male_People02 = "Villager"
	Male_DesertPeople01 = "Villager"
    VisitingMerchant = "Visitor Merchant"




    @classmethod
    def get(self, value, origin_value):
        for i in PalType:
            if i.name.upper() == value:
                return i.value
        log(f"PalType {origin_value} need to be translated", "WARN")
        return origin_value


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

        typename = data["CharacterID"]["value"]
        typename_upper = typename.upper()
        if typename_upper[:5] == "BOSS_":
            typename_upper = typename_upper.replace("BOSS_", "")
            self.is_boss = not self.is_lucky
        self.is_tower = typename_upper.startswith("GYM_")
        self.type = PalType.get(typename_upper, typename)
        self.workspeed = data["CraftSpeed"]["value"]
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
            self.trans_skill(data["PassiveSkillList"]["value"]["values"])
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

    def trans_skill(self, skill_values):
        skills = []
        for skill in skill_values:
            if skill in PalSkills.__members__:
                skills.append(PalSkills[skill].value)
            else:
                skills.append(skill)
                log(f"Unknown skill {skill}", "WARN")
        return skills

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
        self.base_ids = data["base_ids"]
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
