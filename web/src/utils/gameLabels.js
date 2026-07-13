const statusPointKeys = {
  ジャンプ力: "jumpPower",
  jumppower: "jumpPower",
  jumpforce: "jumpPower",
  スタミナ消費軽減: "staminaConsumptionReduction",
  staminaconsumptionreduction: "staminaConsumptionReduction",
  staminacostreduction: "staminaConsumptionReduction",
  パルスフィアホーミング: "palSphereHoming",
  palspherehoming: "palSphereHoming",
  作業速度: "workSpeed",
  workspeed: "workSpeed",
  崖登り速度: "climbingSpeed",
  climbingspeed: "climbingSpeed",
  climbspeed: "climbingSpeed",
  所持重量: "carryWeight",
  maxweight: "carryWeight",
  carryweight: "carryWeight",
  捕獲率: "captureRate",
  capturerate: "captureRate",
  攻撃力: "attack",
  attack: "attack",
  最大hp: "maxHp",
  maxhp: "maxHp",
  最大sp: "maxSp",
  maxsp: "maxSp",
  泳ぎ速度: "swimSpeed",
  swimspeed: "swimSpeed",
  滑空速度: "glideSpeed",
  glidespeed: "glideSpeed",
  状態異常耐性: "statusAilmentResistance",
  statusailmentresistance: "statusAilmentResistance",
  statusresistance: "statusAilmentResistance",
  移動速度アップ: "movementSpeedUp",
  movementspeedup: "movementSpeedUp",
  空腹率低減: "hungerRateReduction",
  hungerratereduction: "hungerRateReduction",
  satietydecelerate: "hungerRateReduction",
  経験値ボーナス: "expBonus",
  expbonus: "expBonus",
  experiencebonus: "expBonus",
  虹パッシブ率: "rainbowPassiveRate",
  rainbowpassiverate: "rainbowPassiveRate",
  食料腐敗低減: "foodSpoilageReduction",
  foodspoilagereduction: "foodSpoilageReduction",
  fooddecayslowdown: "foodSpoilageReduction",
};

const skillAliases = {
  Deffence_up2_2: {
    en: "Defense Boost Lv.2",
    zh: "防御强化 Lv.2",
    ja: "防御強化 Lv.2",
  },
  ReloadSpeedUp_Passive: {
    en: "Reload Speed Boost",
    zh: "装填速度提升",
    ja: "リロード速度上昇",
  },
};

const normalizeStatusPointKey = (rawKey) =>
  String(rawKey)
    .trim()
    .replace(/[\s_-]/g, "")
    .toLowerCase();

export const statusPointTranslationKey = (rawKey = "") =>
  statusPointKeys[normalizeStatusPointKey(rawKey)] || null;

export const localizedSkillName = (skill, locale, skillMap) =>
  skillMap?.[locale]?.[skill]?.name ||
  skillAliases[skill]?.[locale] ||
  skillAliases[skill]?.en ||
  skill;
