export const hasMapLocation = (player) => {
  const x = Number(player?.location_x);
  const y = Number(player?.location_y);
  return Number.isFinite(x) && Number.isFinite(y) && (x !== 0 || y !== 0);
};

export const mergeMapPlayers = (players = [], onlinePlayers = []) => {
  const merged = new Map(
    players
      .filter((player) => player?.player_uid)
      .map((player) => [player.player_uid, { ...player }]),
  );

  onlinePlayers.forEach((player) => {
    if (!player?.player_uid) return;
    merged.set(player.player_uid, {
      ...merged.get(player.player_uid),
      ...player,
    });
  });

  return Array.from(merged.values());
};

export const selectVisibleMapPlayers = (
  players = [],
  onlinePlayerIds = new Set(),
  visibility = "online",
) =>
  players.filter(
    (player) =>
      hasMapLocation(player) &&
      (visibility === "all" || onlinePlayerIds.has(player.player_uid)),
  );

export const buildPlayerGuildMap = (guilds = []) => {
  const playerGuilds = new Map();
  guilds.forEach((guild) => {
    (guild?.players || []).forEach((player) => {
      if (player?.player_uid) playerGuilds.set(player.player_uid, guild);
    });
  });
  return playerGuilds;
};
