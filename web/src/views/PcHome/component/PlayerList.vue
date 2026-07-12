<script setup>
import ApiService from "@/service/api";
import pageStore from "@/stores/model/page.js";
import { ref, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import skillMap from "@/assets/skill.json";
import { NAvatar, NTag } from "naive-ui";
import PlayerDetail from "./PlayerDetail.vue";
import playerToGuildStore from "@/stores/model/playerToGuild";
import whitelistStore from "@/stores/model/whitelist";

const { t, locale } = useI18n();

const props = defineProps({
  showWhitelistPlayer: String,
  players: { type: Array, default: () => [] },
});
const showWhitelistPlayer = computed(() => props.showWhitelistPlayer);

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches
);

const pageWidth = computed(() => pageStore().getScreenWidth());
const smallScreen = computed(() => pageWidth.value < 1024);

const loadingPlayer = ref(false);
const loadingPlayerDetail = ref(false);
const playerList = ref([]);
const playerInfo = ref(null);
const playerPalsList = ref([]);
const skillTypeList = ref([]);
const searchValue = ref("");
const statusFilter = ref("all");
const platformFilter = ref("all");
const whitelistFilter = ref("all");
const sortBy = ref("last_online");
// 平台标记颜色
const platformColors = {
  steam: { color: '#223D58', textColor: '#fff' },   // 青底白字
  xbox:  { color: '#2B8B2B', textColor: '#fff' },   // 绿底白字
  ps5:   { color: '#00439C', textColor: '#fff' },   // 蓝底白字
  mac:   { color: '#999999', textColor: '#fff' },   // 灰底白字
  default: { color: '#d9c36c', textColor: '#fff' }  // 其他平台
}

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
  }
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
  skillTypeList.value = getSkillTypeList();
  loadingPlayer.value = false;
  if (playerList.value.length > 0) {
    const currentUid = playerToGuildStore().getCurrentUid();
    await getPlayerInfo(
      currentUid ? currentUid : playerList.value[0].player_uid
    );
    playerToGuildStore().setCurrentUid(null);
  }
  loadingPlayerDetail.value = false;
});

// 其他操作
const getUserAvatar = () => {
  return new URL("@/assets/avatar.webp", import.meta.url).href;
};
const getSkillTypeList = () => {
  if (skillMap[locale.value]) {
    return Object.values(skillMap[locale.value]).map((item) => item.name);
  } else {
    return [];
  }
};
const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 80000;
};
const getPlatformColor = (userId) => {
  if (!userId) return platformColors.default;
  const platform = userId.split('_')[0];
  return platformColors[platform] || platformColors.default;
}
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
    ...[...platforms].sort().map((platform) => ({ label: platform, value: platform })),
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
    if (platformFilter.value !== "all" && platform !== platformFilter.value) return false;
    const whitelisted = isWhite.value(player);
    if (whitelistFilter.value === "whitelist" && !whitelisted) return false;
    if (whitelistFilter.value === "non-whitelist" && whitelisted) return false;
    return true;
  });
  return filtered.sort((a, b) => {
    if (sortBy.value === "level") return Number(b.level || 0) - Number(a.level || 0);
    if (sortBy.value === "nickname") return (a.nickname || "").localeCompare(b.nickname || "");
    return dayjs(b.last_online).valueOf() - dayjs(a.last_online).valueOf();
  });
});
</script>
<template>
  <div class="paler-list h-full">
    <n-layout has-sider class="h-full">
      <n-layout-sider
        :width="smallScreen ? 360 : 400"
        content-style="padding: 24px;"
        :native-scrollbar="false"
        bordered
        class="relative"
      >
        <div class="mb-3">
          <n-input
            v-model:value="searchValue"
            clearable
            :placeholder="$t('filter.searchPlayers')"
            aria-label="Search players"
          />
          <n-grid cols="2" :x-gap="8" :y-gap="8" class="mt-2">
            <n-gi><n-select v-model:value="statusFilter" :options="statusOptions" aria-label="Player status" /></n-gi>
            <n-gi><n-select v-model:value="platformFilter" :options="platformOptions" aria-label="Player platform" /></n-gi>
            <n-gi><n-select v-model:value="whitelistFilter" :options="whitelistOptions" aria-label="Whitelist status" /></n-gi>
            <n-gi><n-select v-model:value="sortBy" :options="sortOptions" aria-label="Player sorting" /></n-gi>
          </n-grid>
          <n-text depth="3" class="block mt-2">
            {{ $t("filter.resultCount", { count: filteredPlayers.length }) }}
          </n-text>
        </div>
        <n-list hoverable clickable>
          <n-list-item
            v-for="player in filteredPlayers"
            :key="player.player_uid"
            style="padding: 12px 8px"
            @click="clickGetPlayerInfo(player.player_uid)"
            @keydown.enter="clickGetPlayerInfo(player.player_uid)"
            @keydown.space.prevent="clickGetPlayerInfo(player.player_uid)"
            role="button"
            tabindex="0"
            :aria-label="`${player.nickname}, Lv.${player.level}`"
          >
            <template #prefix>
              <n-avatar :src="getUserAvatar()" round></n-avatar>
            </template>
            <div>
              <div class="flex">
                <n-tag
                  :bordered="false"
                  size="small"
                  :type="
                    isPlayerOnline(player.last_online) ? 'success' : 'error'
                  "
                  round
                >
                  {{
                    isPlayerOnline(player.last_online)
                      ? $t("status.online")
                      : $t("status.offline")
                  }}
                </n-tag>
                <n-tag class="ml-2" type="primary" size="small" round>
                  Lv.{{ player.level }}
                </n-tag>
                <n-tag
                  v-if="isWhite(player)"
                  class="ml-2"
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
                <span class="flex-1 pl-2 font-bold line-clamp-1">
                  {{ player.nickname }}
                  <n-tag
                      v-if="player.user_id"
                      class=""
                      style="line-height: 22px;"
                      :bordered="false"
                      round
                      size="small"
                      :color="getPlatformColor(player.user_id)"
                  >
                  {{ player.user_id.split("_")[0] }}
                </n-tag>
                </span>
              </div>
              <n-tag :bordered="false" round size="small" class="mt-2">
                {{ $t("status.last_online") }}:
                {{ displayLastOnline(player.last_online) }}
              </n-tag>
            </div>
          </n-list-item>
        </n-list>
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
