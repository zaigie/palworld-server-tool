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

const props = defineProps(["showWhitelistPlayer"]);
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

// 获取玩家列表
const getPlayerList = async () => {
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = data.value;
};

// 获取玩家详情信息
const getPlayerInfo = async (player_uid) => {
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  playerInfo.value = data.value;
  playerPalsList.value = JSON.parse(JSON.stringify(playerInfo?.value.pals));
  nextTick(() => {
    const playerInfoEL = document.getElementById("player-info");
    if (playerInfoEL) {
      playerInfoEL.scrollIntoView({ behavior: "smooth" });
    }
  });
};

const clickGetPlayerInfo = async (id) => {
  if (playerInfo.value.player_uid !== id) {
    loadingPlayerDetail.value = true;
    await getPlayerInfo(id);
    loadingPlayerDetail.value = false;
  }
};

watch(
  () => showWhitelistPlayer.value,
  async (newVal) => {
    if (playerInfo.value.player_uid !== newVal) {
      loadingPlayerDetail.value = true;
      await getPlayerInfo(newVal);
      loadingPlayerDetail.value = false;
    }
  }
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
  return dayjs() - dayjs(last_online) < 120000;
};
const displayLastOnline = (last_online) => {
  if (dayjs(last_online).year() < 1970) {
    return "Unknown";
  }
  return dayjs(last_online).format("YYYY-MM-DD HH:mm:ss");
};
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
        <n-list hoverable clickable>
          <n-list-item
            v-for="player in playerList"
            :key="player.player_uid"
            style="padding: 12px 8px"
            @click="clickGetPlayerInfo(player.player_uid)"
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
                <span class="flex-1 pl-2 font-bold line-clamp-1">{{
                  player.nickname
                }}</span>
              </div>
              <span
                :class="
                  isDarkMode ? 'bg-#2f69aa text-#fff' : 'bg-#ddd text-#18181c'
                "
                class="inline-block mt-1 rounded-full text-xs px-2 py-0.5"
                >{{ $t("status.last_online") }}:
                {{ displayLastOnline(player.last_online) }}</span
              >
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
          :whiteList="whiteList"
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
