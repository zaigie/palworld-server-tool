<script setup>
import ApiService from "@/service/api";
import pageStore from "@/stores/model/page.js";
import { ref, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import { ChevronForward } from "@vicons/ionicons5";
import PlayerDetail from "./PlayerDetail.vue";
import playerToGuildStore from "@/stores/model/playerToGuild";
import whitelistStore from "@/stores/model/whitelist";

const { t } = useI18n();

const props = defineProps({
  showWhitelistPlayer: String,
  players: { type: Array, default: () => [] },
});
const showWhitelistPlayer = computed(() => props.showWhitelistPlayer);

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);

const pageWidth = computed(() => pageStore().getScreenWidth());
const smallScreen = computed(() => pageWidth.value < 1024);

const loadingPlayer = ref(false);
const loadingPlayerDetail = ref(false);
const playerList = ref([]);
const playerInfo = ref(null);
const playerPalsList = ref([]);
const searchValue = ref("");
const statusFilter = ref("all");
const platformFilter = ref("all");
const whitelistFilter = ref("all");
const sortBy = ref("last_online");
// 平台标记颜色
const platformColors = {
  steam: { color: "#223D58", textColor: "#fff" }, // 青底白字
  xbox: { color: "#2B8B2B", textColor: "#fff" }, // 绿底白字
  ps5: { color: "#00439C", textColor: "#fff" }, // 蓝底白字
  mac: { color: "#999999", textColor: "#fff" }, // 灰底白字
  default: { color: "#d9c36c", textColor: "#fff" }, // 其他平台
};

// 获取玩家列表
const getPlayerList = async () => {
  if (props.players.length > 0) {
    playerList.value = [...props.players];
    return;
  }
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = Array.isArray(data.value) ? data.value : [];
};

// 获取玩家详情信息
const getPlayerInfo = async (player_uid) => {
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  playerInfo.value = data.value;
  playerPalsList.value = playerInfo?.value.pals
    ? JSON.parse(JSON.stringify(playerInfo?.value.pals))
    : [];
  nextTick(() => {
    const playerInfoEL = document.getElementById("player-info");
    if (playerInfoEL) {
      playerInfoEL.scrollIntoView({ behavior: "smooth" });
    }
  });
};

const clickGetPlayerInfo = async (id) => {
  if (playerInfo.value?.player_uid !== id) {
    loadingPlayerDetail.value = true;
    await getPlayerInfo(id);
    loadingPlayerDetail.value = false;
  }
};

watch(
  () => showWhitelistPlayer.value,
  async (newVal) => {
    if (newVal && playerInfo.value?.player_uid !== newVal) {
      loadingPlayerDetail.value = true;
      await getPlayerInfo(newVal);
      loadingPlayerDetail.value = false;
    }
  },
);

watch(
  () => props.players,
  (players) => {
    if (players?.length > 0) playerList.value = [...players];
  },
  { deep: true },
);

// 白名单
const whiteList = computed(() => whitelistStore().getWhitelist());
const isWhite = computed(() => (player) => {
  if (player) {
    return whiteList.value.some((whitelistItem) => {
      return (
        (whitelistItem.player_uid &&
          whitelistItem.player_uid === player.player_uid) ||
        (whitelistItem.steam_id && whitelistItem.steam_id === player.steam_id)
      );
    });
  } else {
    return false;
  }
});

onMounted(async () => {
  loadingPlayerDetail.value = true;
  loadingPlayer.value = true;
  await getPlayerList();
  loadingPlayer.value = false;
  if (playerList.value.length > 0) {
    const currentUid = playerToGuildStore().getCurrentUid();
    await getPlayerInfo(
      currentUid ? currentUid : playerList.value[0].player_uid,
    );
    playerToGuildStore().setCurrentUid(null);
  }
  loadingPlayerDetail.value = false;
});

// 其他操作
const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 80000;
};
const getPlatformColor = (userId) => {
  if (!userId) return platformColors.default;
  const platform = userId.split("_")[0];
  return platformColors[platform] || platformColors.default;
};
const displayLastOnline = (last_online) => {
  if (dayjs(last_online).year() < 1970) {
    return "Unknown";
  }
  return dayjs(last_online).format("YYYY-MM-DD HH:mm:ss");
};

const platformOptions = computed(() => {
  const platforms = new Set(
    playerList.value
      .map((player) => player.user_id?.split("_")[0])
      .filter(Boolean),
  );
  return [
    { label: t("filter.allPlatforms"), value: "all" },
    ...[...platforms]
      .sort()
      .map((platform) => ({ label: platform, value: platform })),
  ];
});
const statusOptions = computed(() => [
  { label: t("filter.allStatuses"), value: "all" },
  { label: t("status.online"), value: "online" },
  { label: t("status.offline"), value: "offline" },
]);
const whitelistOptions = computed(() => [
  { label: t("filter.allPlayers"), value: "all" },
  { label: t("filter.whitelistOnly"), value: "whitelist" },
  { label: t("filter.nonWhitelistOnly"), value: "non-whitelist" },
]);
const sortOptions = computed(() => [
  { label: t("filter.lastOnline"), value: "last_online" },
  { label: t("filter.levelHighToLow"), value: "level" },
  { label: t("filter.nickname"), value: "nickname" },
]);
const filteredPlayers = computed(() => {
  const keyword = searchValue.value.trim().toLowerCase();
  const filtered = playerList.value.filter((player) => {
    const searchable = [
      player.nickname,
      player.player_uid,
      player.user_id,
      player.steam_id,
    ]
      .filter(Boolean)
      .join(" ")
      .toLowerCase();
    if (keyword && !searchable.includes(keyword)) return false;
    const online = isPlayerOnline(player.last_online);
    if (statusFilter.value === "online" && !online) return false;
    if (statusFilter.value === "offline" && online) return false;
    const platform = player.user_id?.split("_")[0];
    if (platformFilter.value !== "all" && platform !== platformFilter.value)
      return false;
    const whitelisted = isWhite.value(player);
    if (whitelistFilter.value === "whitelist" && !whitelisted) return false;
    if (whitelistFilter.value === "non-whitelist" && whitelisted) return false;
    return true;
  });
  return filtered.sort((a, b) => {
    if (sortBy.value === "level")
      return Number(b.level || 0) - Number(a.level || 0);
    if (sortBy.value === "nickname")
      return (a.nickname || "").localeCompare(b.nickname || "");
    return dayjs(b.last_online).valueOf() - dayjs(a.last_online).valueOf();
  });
});
</script>
<template>
  <div class="paler-list h-full" :class="{ 'is-dark': isDarkMode }">
    <n-layout has-sider class="h-full">
      <n-layout-sider
        :width="smallScreen ? 360 : 400"
        content-style="padding: 0 16px 16px;"
        :native-scrollbar="false"
        bordered
        class="player-list-sidebar relative"
      >
        <div class="filter-panel">
          <n-input
            v-model:value="searchValue"
            clearable
            size="large"
            :placeholder="$t('filter.searchPlayers')"
            :aria-label="$t('filter.searchPlayers')"
          />
          <n-grid cols="2" :x-gap="8" :y-gap="8" class="mt-2">
            <n-gi
              ><n-select v-model:value="statusFilter" :options="statusOptions"
            /></n-gi>
            <n-gi
              ><n-select
                v-model:value="platformFilter"
                :options="platformOptions"
            /></n-gi>
            <n-gi
              ><n-select
                v-model:value="whitelistFilter"
                :options="whitelistOptions"
            /></n-gi>
            <n-gi
              ><n-select v-model:value="sortBy" :options="sortOptions"
            /></n-gi>
          </n-grid>
          <n-text depth="3" class="result-count">
            {{ $t("filter.resultCount", { count: filteredPlayers.length }) }}
          </n-text>
        </div>
        <n-list class="player-list" :show-divider="false">
          <n-list-item
            v-for="player in filteredPlayers"
            :key="player.player_uid"
            class="player-row-item"
            @click="clickGetPlayerInfo(player.player_uid)"
            @keydown.enter="clickGetPlayerInfo(player.player_uid)"
            @keydown.space.prevent="clickGetPlayerInfo(player.player_uid)"
            role="button"
            tabindex="0"
            :aria-label="`${player.nickname}, Lv.${player.level}`"
            :aria-current="
              playerInfo?.player_uid === player.player_uid ? 'true' : undefined
            "
          >
            <div
              class="player-row"
              :class="{
                'is-selected': playerInfo?.player_uid === player.player_uid,
              }"
            >
              <div class="player-row-main">
                <div class="player-name-line">
                  <span class="player-name" :title="player.nickname">
                    {{ player.nickname || "--" }}
                  </span>
                  <n-tag
                    v-if="player.user_id"
                    :bordered="false"
                    round
                    size="small"
                    :color="getPlatformColor(player.user_id)"
                  >
                    {{ player.user_id.split("_")[0] }}
                  </n-tag>
                </div>
                <div class="player-meta-line">
                  <span class="player-status">
                    <span
                      class="status-dot"
                      :class="
                        isPlayerOnline(player.last_online)
                          ? 'is-online'
                          : 'is-offline'
                      "
                    ></span>
                    {{
                      isPlayerOnline(player.last_online)
                        ? $t("status.online")
                        : $t("status.offline")
                    }}
                  </span>
                  <n-tag :bordered="false" type="primary" size="small" round>
                    Lv.{{ player.level }}
                  </n-tag>
                  <n-tag
                    v-if="isWhite(player)"
                    :bordered="false"
                    round
                    size="small"
                    :color="{
                      color: isDarkMode ? '#fff' : '#d9c36c',
                      textColor: isDarkMode ? '#d9c36c' : '#fff',
                    }"
                  >
                    {{ $t("status.whitelist") }}
                  </n-tag>
                </div>
                <div class="last-online">
                  {{ $t("status.last_online") }}
                  <span>{{ displayLastOnline(player.last_online) }}</span>
                </div>
              </div>
              <n-icon class="row-chevron" size="18">
                <ChevronForward />
              </n-icon>
            </div>
          </n-list-item>
        </n-list>
        <n-empty
          v-if="!loadingPlayer && filteredPlayers.length === 0"
          class="empty-state"
        />
        <n-spin
          size="small"
          v-if="loadingPlayer"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description> 加载中... </template>
        </n-spin>
      </n-layout-sider>
      <n-layout :native-scrollbar="false" class="relative">
        <player-detail
          :playerInfo="playerInfo"
          :playerPalsList="playerPalsList"
        ></player-detail>
        <n-spin
          size="small"
          v-if="loadingPlayerDetail"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description> 加载中... </template>
        </n-spin>
      </n-layout>
    </n-layout>
  </div>
</template>

<style scoped lang="less">
.filter-panel {
  position: sticky;
  top: 0;
  z-index: 3;
  padding: 20px 0 14px;
  background: #fff;
}

.is-dark .filter-panel {
  background: #18181c;
}

.result-count {
  display: block;
  margin-top: 10px;
  padding: 0 2px;
  font-size: 13px;
}

.player-list {
  background: transparent;
}

.player-row-item {
  padding: 0 0 4px !important;
  outline: none;
}

.player-row-item:focus-visible .player-row {
  box-shadow: 0 0 0 2px rgba(64, 152, 252, 0.45);
}

.player-row {
  width: 100%;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition:
    background-color 0.18s ease,
    border-color 0.18s ease,
    box-shadow 0.18s ease;

  &:hover {
    background: rgba(64, 152, 252, 0.08);
  }

  &.is-selected {
    border-color: rgba(64, 152, 252, 0.55);
    background: rgba(64, 152, 252, 0.12);
    box-shadow: inset 3px 0 0 #4098fc;
  }
}

.is-dark .player-row:hover {
  background: rgba(64, 152, 252, 0.12);
}

.is-dark .player-row.is-selected {
  background: rgba(64, 152, 252, 0.17);
}

.player-row-main {
  flex: 1;
  min-width: 0;
}

.player-name-line {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.player-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  color: inherit;
  font-size: 17px;
  font-weight: 650;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-meta-line {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 7px;
  margin-top: 5px;
}

.player-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: rgba(24, 24, 28, 0.65);
  font-size: 12px;
}

.is-dark .player-status {
  color: rgba(255, 255, 255, 0.62);
}

.status-dot {
  width: 8px;
  height: 8px;
  flex: none;
  border-radius: 50%;

  &.is-online {
    background: #18a058;
    box-shadow: 0 0 0 3px rgba(24, 160, 88, 0.14);
  }

  &.is-offline {
    background: #d03050;
  }
}

.last-online {
  margin-top: 5px;
  color: rgba(24, 24, 28, 0.45);
  font-size: 12px;
  line-height: 1.25;
}

.last-online span {
  margin-left: 5px;
  color: rgba(24, 24, 28, 0.68);
  font-variant-numeric: tabular-nums;
}

.is-dark .last-online {
  color: rgba(255, 255, 255, 0.38);
}

.is-dark .last-online span {
  color: rgba(255, 255, 255, 0.62);
}

.row-chevron {
  flex: none;
  color: rgba(24, 24, 28, 0.28);
  transition: transform 0.18s ease;
}

.player-row:hover .row-chevron,
.player-row.is-selected .row-chevron {
  color: #4098fc;
  transform: translateX(2px);
}

.is-dark .row-chevron {
  color: rgba(255, 255, 255, 0.3);
}

.empty-state {
  padding: 48px 12px;
}
</style>
