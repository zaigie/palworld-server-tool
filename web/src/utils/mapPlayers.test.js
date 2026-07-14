import test from "node:test";
import assert from "node:assert/strict";
import {
  buildPlayerGuildMap,
  hasMapLocation,
  mergeMapPlayers,
  selectVisibleMapPlayers,
} from "./mapPlayers.js";

const players = [
  {
    player_uid: "online-player",
    nickname: "Online",
    level: 20,
    location_x: 10,
    location_y: 20,
  },
  {
    player_uid: "offline-player",
    nickname: "Offline",
    level: 10,
    location_x: 30,
    location_y: 40,
  },
  {
    player_uid: "unknown-location",
    nickname: "Unknown",
    location_x: 0,
    location_y: 0,
  },
];

test("defaults map visibility to the authoritative online UID set", () => {
  const visible = selectVisibleMapPlayers(players, new Set(["online-player"]));
  assert.deepEqual(
    visible.map((player) => player.player_uid),
    ["online-player"],
  );
});

test("all-player visibility includes offline players with a last known position", () => {
  const visible = selectVisibleMapPlayers(
    players,
    new Set(["online-player"]),
    "all",
  );
  assert.deepEqual(
    visible.map((player) => player.player_uid),
    ["online-player", "offline-player"],
  );
  assert.equal(hasMapLocation(players[2]), false);
});

test("online refresh overrides live fields without losing saved player data", () => {
  const merged = mergeMapPlayers(players, [
    {
      player_uid: "online-player",
      nickname: "Online now",
      location_x: 50,
      location_y: 60,
    },
    {
      player_uid: "new-online-player",
      nickname: "New",
      location_x: 70,
      location_y: 80,
    },
  ]);
  const refreshed = merged.find(
    (player) => player.player_uid === "online-player",
  );
  assert.equal(refreshed.level, 20);
  assert.equal(refreshed.nickname, "Online now");
  assert.equal(refreshed.location_x, 50);
  assert.equal(merged.length, 4);
});

test("indexes each guild member for player card navigation", () => {
  const guild = {
    name: "Builders",
    players: [{ player_uid: "online-player", nickname: "Online" }],
  };
  assert.equal(buildPlayerGuildMap([guild]).get("online-player"), guild);
});
