const statusPointKeys = {
  "作業速度": "workSpeed",
  workspeed: "workSpeed",
  work_speed: "workSpeed",
  "所持重量": "carryWeight",
  maxweight: "carryWeight",
  max_weight: "carryWeight",
  "捕獲率": "captureRate",
  capturerate: "captureRate",
  capture_rate: "captureRate",
  "攻撃力": "attack",
  attack: "attack",
  "最大hp": "maxHp",
  maxhp: "maxHp",
  max_hp: "maxHp",
  "最大sp": "maxSp",
  maxsp: "maxSp",
  max_sp: "maxSp",
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

export const statusPointTranslationKey = (rawKey = "") =>
  statusPointKeys[String(rawKey).toLowerCase()] || null;

export const localizedSkillName = (skill, locale, skillMap) =>
  skillMap?.[locale]?.[skill]?.name ||
  skillAliases[skill]?.[locale] ||
  skillAliases[skill]?.en ||
  skill;
