export const RCON_PLACEHOLDERS = [
  { key: "playerid", aliases: ["playerid"] },
  { key: "steamid", aliases: ["steamid", "steamuserid"] },
  { key: "userid", aliases: ["userid"] },
  { key: "itemid", aliases: ["itemid"] },
  { key: "palid", aliases: ["palid"] },
];

const placeholderAliases = RCON_PLACEHOLDERS.reduce((result, definition) => {
  for (const alias of definition.aliases) result[alias] = definition.key;
  return result;
}, {});

export const extractRconPlaceholders = (template = "") => {
  const found = new Set();
  for (const match of template.matchAll(/{([a-z]+)}/gi)) {
    found.add(match[1].toLowerCase());
  }
  return [...found];
};

export const resolveRconTemplate = (template = "", context = {}) => {
  const values = {
    playerid: context.player?.player_uid || "",
    steamid: context.player?.steam_id
      ? `steam_${context.player.steam_id}`
      : "",
    userid: context.player?.user_id || "",
    itemid: context.item || "",
    palid: context.pal || "",
  };
  const missing = new Set();
  const unknown = new Set();
  const content = template.replace(/{([a-z]+)}/gi, (match, rawKey) => {
    const alias = rawKey.toLowerCase();
    const key = placeholderAliases[alias];
    if (!key) {
      unknown.add(alias);
      return match;
    }
    if (!values[key]) {
      missing.add(key);
      return match;
    }
    return values[key];
  });
  return { content, missing: [...missing], unknown: [...unknown] };
};

export const describeCron = (expression = "") => {
  const presets = {
    "*/5 * * * *": "every5Minutes",
    "*/15 * * * *": "every15Minutes",
    "*/30 * * * *": "every30Minutes",
    "0 * * * *": "hourly",
  };
  return presets[expression] || "custom";
};
