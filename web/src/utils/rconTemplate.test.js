import test from "node:test";
import assert from "node:assert/strict";
import {
  extractRconPlaceholders,
  resolveRconTemplate,
} from "./rconTemplate.js";

test("resolves documented and legacy RCON placeholders", () => {
  const result = resolveRconTemplate(
    "{playerid} {steamid} {userID} {itemID} {palID} {steamUserID}",
    {
      player: {
        player_uid: "player-1",
        steam_id: "7656119",
        user_id: "steam_7656119",
      },
      item: "Money",
      pal: "SheepBall",
    },
  );
  assert.equal(
    result.content,
    "player-1 steam_7656119 steam_7656119 Money SheepBall steam_7656119",
  );
  assert.deepEqual(result.missing, []);
  assert.deepEqual(result.unknown, []);
});

test("keeps unresolved placeholders visible", () => {
  const result = resolveRconTemplate("Give {playerid} {itemid} {unknown}");
  assert.equal(result.content, "Give {playerid} {itemid} {unknown}");
  assert.deepEqual(result.missing, ["playerid", "itemid"]);
  assert.deepEqual(result.unknown, ["unknown"]);
  assert.deepEqual(extractRconPlaceholders(result.content), [
    "playerid",
    "itemid",
    "unknown",
  ]);
});
