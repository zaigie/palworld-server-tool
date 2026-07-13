import test from "node:test";
import assert from "node:assert/strict";
import { statusPointTranslationKey } from "./gameLabels.js";

const currentStatusPointKeys = {
  ジャンプ力: "jumpPower",
  スタミナ消費軽減: "staminaConsumptionReduction",
  パルスフィアホーミング: "palSphereHoming",
  作業速度: "workSpeed",
  崖登り速度: "climbingSpeed",
  所持重量: "carryWeight",
  捕獲率: "captureRate",
  攻撃力: "attack",
  最大HP: "maxHp",
  最大SP: "maxSp",
  泳ぎ速度: "swimSpeed",
  滑空速度: "glideSpeed",
  状態異常耐性: "statusAilmentResistance",
  移動速度アップ: "movementSpeedUp",
  空腹率低減: "hungerRateReduction",
  経験値ボーナス: "expBonus",
  虹パッシブ率: "rainbowPassiveRate",
  食料腐敗低減: "foodSpoilageReduction",
};

test("maps every current Palworld status point key to a semantic key", () => {
  for (const [rawKey, translationKey] of Object.entries(
    currentStatusPointKeys,
  )) {
    assert.equal(statusPointTranslationKey(rawKey), translationKey, rawKey);
  }
});

test("normalizes legacy and English status point aliases", () => {
  assert.equal(statusPointTranslationKey("work_speed"), "workSpeed");
  assert.equal(statusPointTranslationKey("Max HP"), "maxHp");
  assert.equal(statusPointTranslationKey("jump-force"), "jumpPower");
  assert.equal(
    statusPointTranslationKey("Food Decay Slowdown"),
    "foodSpoilageReduction",
  );
});

test("keeps unknown status point keys available for raw fallback", () => {
  assert.equal(statusPointTranslationKey("FutureStatKey"), null);
});
