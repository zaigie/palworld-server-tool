<script setup>
import {
  AdminPanelSettingsOutlined,
  SupervisedUserCircleRound,
  GroupWorkRound,
  ContentCopyFilled,
  SettingsPowerRound,
} from "@vicons/material";
import { ChevronsLeft } from "@vicons/tabler";
import { GameController, LogOut, Ban, LanguageSharp } from "@vicons/ionicons5";
import { BroadcastTower } from "@vicons/fa";
import { CrownFilled } from "@vicons/antd";
import { computed, onMounted, ref, h } from "vue";
import { NTag, NButton, NAvatar, useMessage, useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import dayjs from "dayjs";
import skillDescMap from "@/assets/skillDesc.json";
import palZHTypes from "@/assets/zhTypes.json";
import palZHSkills from "@/assets/zhSkills.json";
import palJATypes from "@/assets/jaTypes.json";
import palJASkills from "@/assets/jaSkills.json";

const { t, locale } = useI18n();

const message = useMessage();
const dialog = useDialog();

const PALWORLD_TOKEN = "palworld_token";

const loading = ref(false);
const serverInfo = ref({});
const currentDisplay = ref("players");
const isShowDetail = ref(false);
const playerList = ref([]);
const guildList = ref([]);
const playerInfo = ref({});
const playerPalsList = ref([]);
const currentPlayerPalsList = ref([]);
const guildInfo = ref({});
const skillTypeList = ref([]);
const languageOptions = ref([]);

const contentRef = ref(null);

const isLogin = ref(false);
const authToken = ref("");

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches
);

const updateDarkMode = (e) => {
  isDarkMode.value = e.matches;
};

const getUserAvatar = () => {
  return new URL("../assets/avatar.webp", import.meta.url).href;
};

const handleSelectLanguage = (key) => {
  message.info(t("message.changelanguage"));
  if (key === "zh") {
    localStorage.setItem("locale", "zh");
    // locale.value = "zh";
  } else if (key === "ja") {
    localStorage.setItem("locale", "ja");
    // locale.value = "ja";
  } else {
    localStorage.setItem("locale", "en");
    // locale.value = "en";
  }
  setTimeout(() => {
    location.reload();
  }, 1000);
};

const getSkillTypeList = () => {
  if (locale.value === "zh") {
    return Object.values(palZHSkills);
  } else if (local.value === "ja") {
    return Object.values(palJASkills);
  } else if (locale.value === "en") {
    return Object.keys(palZHSkills);
  }
};

const getPalAvatar = (name) => {
  return new URL(`../assets/pal/${name}.png`, import.meta.url).href;
};
const getUnknowPalAvatar = () => {
  return new URL("../assets/pal/Unknown.png", import.meta.url).href;
};

// get data
const getServerInfo = async () => {
  const { data } = await new ApiService().getServerInfo();
  serverInfo.value = data.value;
};
const getPlayerList = async (is_update_info = true) => {
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = data.value;
};
const getGuildList = async () => {
  const { data } = await new ApiService().getGuildList();
  guildList.value = data.value;
};

const getPlayerInfo = async (player_uid) => {
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  playerInfo.value = data.value;
  if (locale.value === "zh") {
    playerInfo.value.pals.forEach((pal) => {
      pal.skills = pal.skills.map((skill) => {
        return palZHSkills[skill] ? palZHSkills[skill] : skill;
      });
      pal.typeName = palZHTypes[pal.type] ? palZHTypes[pal.type] : pal.type;
    });
  } else if (locale.value === "ja") {
    playerInfo.value.pals.forEach((pal) => {
      pal.skills = pal.skills.map((skill) => {
        return palJASkills[skill] ? palJASkills[skill] : skill;
      });
      pal.typeName = palJATypes[pal.type] ? palJATypes[pal.type] : pal.type;
    });
  } else {
    playerInfo.value.pals.forEach((pal) => {
      pal.typeName = pal.type;
    });
  }
  playerPalsList.value = JSON.parse(JSON.stringify(playerInfo.value.pals));
  currentPlayerPalsList.value = playerPalsList.value.slice(0, pageSize.value);
  isShowDetail.value = true;
  contentRef.value.scrollTo(0, 0);
};

const getGuildInfo = async (admin_player_uid) => {
  const { data } = await new ApiService().getGuild({
    adminPlayerUid: admin_player_uid,
  });
  guildInfo.value = data.value;
  isShowDetail.value = true;
  contentRef.value.scrollTo(0, 0);
};

// 游戏用户的帕鲁列表分页，搜索
const searchValue = ref("");
const clickSearch = () => {
  const pattern = /^\s*$|(\s)\1/;
  if (searchValue.value && !pattern.test(searchValue.value)) {
    playerPalsList.value = playerInfo.value.pals.filter((item) => {
      return (
        item.skills.some((skill) => skill.includes(searchValue.value)) ||
        item.typeName.includes(searchValue.value)
      );
    });
  } else {
    playerPalsList.value = JSON.parse(JSON.stringify(playerInfo.value.pals));
  }
  currentPage.value = 1;
  if (playerPalsList.value.length <= 10) {
    finished.value = true;
    currentPlayerPalsList.value = playerPalsList.value ?? [];
  } else {
    finished.value = false;
    currentPlayerPalsList.value = playerPalsList.value.slice(0, pageSize.value);
  }
};
// 滚动加载更多
const palsLoading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const finished = ref(false);
const onLoadPals = () => {
  if (playerPalsList.value.length <= currentPage.value * pageSize.value) {
    finished.value = true;
  } else {
    currentPage.value += 1;
    currentPlayerPalsList.value = playerPalsList.value.slice(
      0,
      pageSize.value * currentPage.value
    );
  }
};
const onContentScroll = () => {
  if (currentDisplay.value === "players" && isShowDetail.value) {
    const dom = document.getElementsByClassName("n-layout-scroll-container");
    if (dom[1].scrollTop + dom[1].clientHeight > dom[1].scrollHeight - 6) {
      onLoadPals();
    }
  }
};

const showPalDetailModal = ref(false);
const palDetail = ref({});

const showPalDetail = (pal) => {
  palDetail.value = pal;
  showPalDetailModal.value = true;
};
const dataRowProps = (row) => {
  return {
    onClick: () => showPalDetail(row),
  };
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
const getOnlineList = () => {
  return playerList.value.filter((player) =>
    isPlayerOnline(player.last_online)
  );
};

const displayHP = (hp, max_hp) => {
  return (hp / 1000).toFixed(0) + "/" + (max_hp / 1000).toFixed(0);
};

const percentageHP = (hp, max_hp) => {
  if (max_hp === 0) {
    return 0;
  }
  return ((hp / max_hp) * 100).toFixed(2);
};

const copyText = (text) => {
  const textarea = document.createElement("textarea");
  textarea.value = text;
  document.body.appendChild(textarea);
  textarea.select();

  try {
    const successful = document.execCommand("copy");
    message.success(t("message.copysuccess"));
  } catch (err) {
    message.error(t("message.copyerr", { err: err }));
  }

  document.body.removeChild(textarea);
};

// login
const showLoginModal = ref(false);
const password = ref("");
const handleLogin = async () => {
  const { data, statusCode } = await new ApiService().login({
    password: password.value,
  });
  if (statusCode.value === 401) {
    message.error(t("message.autherr"));
    password.value = "";
    return;
  }
  let token = data.value.token;
  localStorage.setItem(PALWORLD_TOKEN, token);
  authToken.value = token;
  message.success(t("message.authsuccess"));
  showLoginModal.value = false;
  isLogin.value = true;
};

const handelPlayerAction = async (type) => {
  if (!checkAuthToken()) {
    message.error($t("message.requireauth"));
    showLoginModal.value = true;
    return;
  }
  dialog.warning({
    title: type === "ban" ? t("message.bantitle") : t("message.kicktitle"),
    content: type === "ban" ? t("message.banwarn") : t("message.kickwarn"),
    positiveText: t("button.confirm"),
    negativeText: t("button.cancel"),
    onPositiveClick: async () => {
      if (type === "ban") {
        const { data, statusCode } = await new ApiService().banPlayer({
          playerUid: playerInfo.value.player_uid,
        });
        if (statusCode.value === 200) {
          message.success(t("message.bansuccess"));
        } else {
          message.error(t("message.banfail", { err: data.value?.error }));
        }
      } else if (type === "kick") {
        const { data, statusCode } = await new ApiService().kickPlayer({
          playerUid: playerInfo.value.player_uid,
        });
        if (statusCode.value === 200) {
          message.success(t("message.kicksuccess"));
        } else {
          message.error(t("message.kickfail", { err: data.value?.error }));
        }
      }
    },
  });
};

// broadcast
const showBroadcastModal = ref(false);
const broadcastText = ref("");
const handleStartBrodcast = () => {
  // broadcast start
  if (checkAuthToken()) {
    showBroadcastModal.value = true;
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};
const handleBroadcast = async () => {
  const { data, statusCode } = await new ApiService().sendBroadcast({
    message: broadcastText.value,
  });
  if (statusCode.value === 200) {
    message.success(t("message.broadcastsuccess"));
    showBroadcastModal.value = false;
    broadcastText.value = "";
  } else {
    if (data.value?.error.includes("contain non-ascii")) {
      message.error(t("message.broadcastasciierr"));
      return;
    }
    message.error(t("message.broadcastfail", { err: data.value?.error }));
  }
};

//shutdown
const doShutdown = async () => {
  return await new ApiService().shutdownServer({
    seconds: 60,
    message: "Server Will Shutdown After 60 Seconds",
  });
};

const handleShutdown = () => {
  if (checkAuthToken()) {
    dialog.warning({
      title: t("message.warn"),
      content: t("message.shutdowntip"),
      positiveText: t("button.confirm"),
      negativeText: t("button.cancel"),
      onPositiveClick: async () => {
        const { data, statusCode } = await doShutdown();
        if (statusCode.value === 200) {
          message.success(t("message.shutdownsuccess"));
          return;
        } else {
          message.error(t("message.shutdownfail", { err: data.value?.error }));
        }
      },
      onNegativeClick: () => {},
    });
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};

const toPlayers = async () => {
  if (currentDisplay.value === "players") {
    return;
  }
  await getPlayerList();
  currentDisplay.value = "players";
  isShowDetail.value = false;

  palsLoading.value = false;
  finished.value = false;
  currentPage.value = 1;
  searchValue.value = "";

  contentRef.value.scrollTo(0, 0);
};
const toGuilds = async () => {
  if (currentDisplay.value === "guilds") {
    return;
  }
  await getGuildList();
  currentDisplay.value = "guilds";
  isShowDetail.value = false;

  palsLoading.value = false;
  finished.value = false;
  currentPage.value = 1;
  searchValue.value = "";

  contentRef.value.scrollTo(0, 0);
};
const returnList = () => {
  isShowDetail.value = false;

  palsLoading.value = false;
  finished.value = false;
  currentPage.value = 1;
  searchValue.value = "";

  contentRef.value.scrollTo(0, 0);
};

/**
 * check auth token
 */
const checkAuthToken = () => {
  const token = localStorage.getItem(PALWORLD_TOKEN);
  if (token && token !== "") {
    isLogin.value = true;
    authToken.value = token;
    return true;
  }
  return false;
};

onMounted(async () => {
  locale.value = localStorage.getItem("locale");
  languageOptions.value = [
    {
      label: "简体中文",
      key: "zh",
      disabled: locale.value == "zh",
    },
    {
      label: "English",
      key: "en",
      disabled: locale.value == "en",
    },
    {
      label: "日本語",
      key: "ja",
      disabled: locale.value == "ja",
    },
  ];
  const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  mediaQuery.addEventListener("change", updateDarkMode);
  isDarkMode.value = mediaQuery.matches;

  skillTypeList.value = getSkillTypeList();
  loading.value = true;
  checkAuthToken();
  getServerInfo();
  await getPlayerList();
  loading.value = false;
  setInterval(() => {
    getPlayerList(false);
  }, 60000);
});
</script>

<template>
  <div class="home-page overflow-hidden">
    <div
      :class="isDarkMode ? 'bg-#18181c text-#fff' : 'bg-#fff text-#18181c'"
      class="flex justify-between items-center p-3"
    >
      <div>
        <span class="line-clamp-1 text-base">{{ $t("title") }}</span>
        <n-tag type="default" size="small">{{
          serverInfo?.name
            ? `${serverInfo.name + " " + serverInfo.version}`
            : "获取中..."
        }}</n-tag>
      </div>
      <n-space vertical>
        <n-space justify="end">
          <n-tag type="info" round size="small">{{
            $t("status.player_number", { number: playerList.length })
          }}</n-tag>
          <n-tag type="success" round size="small">{{
            $t("status.online_number", { number: getOnlineList().length })
          }}</n-tag>
        </n-space>
        <n-space justify="end" class="flex items-center">
          <n-dropdown
            trigger="hover"
            :options="languageOptions"
            @select="handleSelectLanguage"
          >
            <n-button type="default" secondary strong circle size="small">
              <template #icon>
                <n-icon><LanguageSharp /></n-icon>
              </template>
            </n-button>
          </n-dropdown>

          <n-button
            type="primary"
            size="small"
            secondary
            strong
            @click="showLoginModal = true"
            v-if="!isLogin"
          >
            <template #icon>
              <n-icon>
                <AdminPanelSettingsOutlined />
              </n-icon>
            </template>
            {{ $t("button.auth") }}
          </n-button>
          <n-tag v-else type="success" size="small" round>
            <template #icon>
              <n-icon>
                <AdminPanelSettingsOutlined />
              </n-icon>
            </template>
            {{ $t("status.authenticated") }}
          </n-tag>
        </n-space>
      </n-space>
    </div>
    <div class="w-full">
      <div class="rounded-lg" v-if="!loading && playerList.length > 0">
        <n-layout style="height: calc(100vh - 86px)" has-sider>
          <n-layout-header
            class="flex flex-col justify-between"
            :class="isLogin ? 'h-16' : 'h-10'"
            bordered
          >
            <div v-if="isLogin" class="flex justify-center items-center px-3">
              <n-button
                size="small"
                type="success"
                class="mr-2"
                secondary
                strong
                round
                @click="handleStartBrodcast"
              >
                <template #icon>
                  <n-icon>
                    <BroadcastTower />
                  </n-icon>
                </template>
                {{ $t("button.broadcast") }}
              </n-button>
              <n-button
                size="small"
                type="error"
                secondary
                strong
                round
                @click="handleShutdown"
              >
                <template #icon>
                  <n-icon>
                    <SettingsPowerRound />
                  </n-icon>
                </template>
                {{ $t("button.shutdown") }}
              </n-button>
            </div>
            <div v-else></div>
            <div class="flex justify-end">
              <n-button-group size="small" class="w-full">
                <n-button
                  v-if="isShowDetail"
                  class="w-20%"
                  @click="returnList"
                  type="tertiary"
                  strong
                  secondary
                >
                  <n-icon size="24">
                    <ChevronsLeft />
                  </n-icon>
                </n-button>
                <n-button
                  :class="isShowDetail ? 'w-40%' : 'w-50%'"
                  @click="toPlayers"
                  :type="currentDisplay === 'players' ? 'primary' : 'tertiary'"
                  secondary
                  strong
                >
                  <template #icon>
                    <n-icon>
                      <GameController />
                    </n-icon>
                  </template>
                  {{ $t("button.players") }}
                </n-button>
                <n-button
                  :class="isShowDetail ? 'w-40%' : 'w-50%'"
                  @click="toGuilds"
                  :type="currentDisplay === 'guilds' ? 'primary' : 'tertiary'"
                  secondary
                  strong
                >
                  <template #icon>
                    <n-icon>
                      <SupervisedUserCircleRound />
                    </n-icon>
                  </template>
                  {{ $t("button.guilds") }}
                </n-button>
              </n-button-group>
            </div>
          </n-layout-header>
          <n-layout
            position="absolute"
            style="top: 64px"
            ref="contentRef"
            @scroll="onContentScroll"
          >
            <div v-if="!isShowDetail">
              <!-- list -->
              <n-list v-if="currentDisplay === 'players'" hoverable clickable>
                <n-list-item
                  v-for="player in playerList"
                  :key="player.player_uid"
                  @click="getPlayerInfo(player.player_uid)"
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
                          isPlayerOnline(player.last_online)
                            ? 'success'
                            : 'error'
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
                      <span class="flex-1 pl-2 font-bold line-clamp-1">{{
                        player.nickname
                      }}</span>
                    </div>
                    <span
                      :class="
                        isDarkMode
                          ? 'bg-#2f69aa text-#fff'
                          : 'bg-#ddd text-#18181c'
                      "
                      class="inline-block mt-1 rounded-full text-xs px-2 py-0.5"
                      >{{ $t("status.last_online") }}:
                      {{ displayLastOnline(player.last_online) }}</span
                    >
                  </div>
                </n-list-item>
              </n-list>
              <n-list v-if="currentDisplay === 'guilds'" hoverable clickable>
                <n-list-item
                  v-for="guild in guildList"
                  :key="guild.admin_player_uid"
                  @click="getGuildInfo(guild.admin_player_uid)"
                >
                  <template #prefix>
                    <n-avatar
                      :style="{ color: 'white', backgroundColor: 'darkorange' }"
                      round
                    >
                      <n-icon>
                        <GroupWorkRound />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <n-tag type="primary" size="small" round>
                    Lv.{{ guild.base_camp_level }}
                  </n-tag>
                  <span class="pl-2 font-bold">{{ guild.name }}</span>
                </n-list-item>
              </n-list>
            </div>
            <!-- detail -->
            <div v-else class="relative">
              <!-- ban / kick -->
              <div
                v-if="
                  currentDisplay === 'players' &&
                  isLogin &&
                  playerInfo.player_uid &&
                  isShowDetail
                "
                class="pt-2 px-3 bg-transparent"
                position="absolute"
              >
                <n-flex justify="space-between">
                  <n-button
                    @click="handelPlayerAction('ban')"
                    type="error"
                    size="small"
                    secondary
                    strong
                    round
                  >
                    <template #icon>
                      <n-icon>
                        <Ban />
                      </n-icon>
                    </template>
                    {{ $t("button.ban") }}
                  </n-button>
                  <n-button
                    @click="handelPlayerAction('kick')"
                    type="warning"
                    size="small"
                    secondary
                    strong
                    round
                  >
                    <template #icon>
                      <n-icon>
                        <LogOut />
                      </n-icon>
                    </template>
                    {{ $t("button.kick") }}
                  </n-button>
                </n-flex>
              </div>

              <n-layout
                v-if="currentDisplay === 'players'"
                :native-scrollbar="false"
              >
                <n-card
                  :bordered="false"
                  v-if="playerInfo.nickname"
                  content-style="padding: 12px"
                >
                  <n-page-header>
                    <n-grid :cols="6">
                      <n-gi
                        v-for="status in Object.entries(
                          playerInfo.status_point
                        )"
                        :key="status[0]"
                      >
                        <n-statistic :label="status[0]" :value="status[1]" />
                      </n-gi>
                    </n-grid>
                    <template #title>
                      <div class="flex items-center w-full">
                        <span class="flex-1 text-sm line-clamp-1 pr-1">
                          {{ playerInfo.nickname }}
                        </span>
                        <n-tag
                          :bordered="false"
                          :type="
                            isPlayerOnline(playerInfo.last_online)
                              ? 'success'
                              : 'error'
                          "
                          round
                          size="small"
                        >
                          {{
                            isPlayerOnline(playerInfo.last_online)
                              ? $t("status.online")
                              : $t("status.offline")
                          }}
                        </n-tag>
                      </div>
                      <n-tag
                        @click="copyText(playerInfo.player_uid)"
                        class="mt-1"
                        type="info"
                        size="small"
                        icon-placement="right"
                        ghost
                      >
                        UID: {{ playerInfo.player_uid }}
                        <template #icon>
                          <n-icon><ContentCopyFilled /></n-icon>
                        </template>
                      </n-tag>
                      <n-tag
                        @click="copyText(playerInfo.steam_id)"
                        class="mt-1"
                        type="info"
                        size="small"
                        icon-placement="right"
                        ghost
                      >
                        Steam64:
                        {{ playerInfo.steam_id ? playerInfo.steam_id : "--" }}
                        <template #icon>
                          <n-icon><ContentCopyFilled /></n-icon>
                        </template>
                      </n-tag>
                    </template>
                    <template #avatar>
                      <n-avatar :src="getUserAvatar()" round></n-avatar>
                    </template>
                    <template #extra>
                      <n-space>
                        <n-tag type="primary" size="small" round strong>
                          Lv.{{ playerInfo.level }}
                          <template #icon>
                            <n-icon :component="CrownFilled" />
                          </template>
                        </n-tag>
                      </n-space>
                    </template>
                    <template #footer>
                      <!-- <n-flex justify="end">Updated at 2022-01-01</n-flex> -->
                    </template>
                  </n-page-header>
                  <n-space vertical>
                    <n-progress
                      type="line"
                      status="error"
                      indicator-placement="inside"
                      :percentage="
                        percentageHP(playerInfo.hp, playerInfo.max_hp)
                      "
                      :height="24"
                      :border-radius="4"
                      :fill-border-radius="0"
                      >HP:
                      {{
                        displayHP(playerInfo.hp, playerInfo.max_hp)
                      }}</n-progress
                    >
                    <n-progress
                      type="line"
                      indicator-placement="inside"
                      :percentage="
                        percentageHP(
                          playerInfo.shield_hp,
                          playerInfo.shield_max_hp
                        )
                      "
                      :height="24"
                      :border-radius="4"
                      :fill-border-radius="0"
                      >SHIELD:
                      {{
                        displayHP(
                          playerInfo.shield_hp,
                          playerInfo.shield_max_hp
                        )
                      }}</n-progress
                    >
                  </n-space>
                  <div
                    class="flex w-full mt-5 border-b border-b-solid border-b-#eee"
                  >
                    <van-field
                      v-model="searchValue"
                      :placeholder="$t('input.searchPlaceholder')"
                      @update:model-value="clickSearch"
                      right-icon="search"
                    >
                    </van-field>
                  </div>
                  <van-list :finished="finished" finished-text="没有更多了">
                    <div
                      v-for="(pal, index) in currentPlayerPalsList"
                      :key="pal"
                      class="py-2"
                      :class="
                        index < currentPlayerPalsList.length - 1
                          ? 'border-b border-b-solid border-b-#eee'
                          : ''
                      "
                      @click="showPalDetail(pal)"
                    >
                      <div class="flex justify-between items-center">
                        <van-image
                          class="bg-#c5c5c5 rounded-md"
                          width="32"
                          height="32"
                          :src="getPalAvatar(pal.type)"
                        />
                        <div
                          class="flex-1 flex items-center justify-between ml-3"
                        >
                          <van-tag
                            plain
                            :type="pal.gender == 'Male' ? 'primary' : 'danger'"
                            >{{ pal.gender == "Male" ? "♂" : "♀" }}</van-tag
                          >
                          <span class="px-3 flex-1 line-clamp-1">{{
                            pal.typeName
                          }}</span>
                          <span>{{ "Lv." + pal.level }}</span>
                        </div>
                      </div>
                      <div class="ml-11 mt-1 flex flex-wrap">
                        <van-tag
                          v-for="skill in pal.skills"
                          class="rounded-sm mr-2"
                          size="medium"
                          :key="skill"
                          color="#fcf0e0"
                          text-color="#ee9b2f"
                          >{{ skill }}</van-tag
                        >
                      </div>
                    </div>
                  </van-list>
                  <div class="h-10"></div>
                </n-card>
                <n-modal
                  v-model:show="showPalDetailModal"
                  preset="card"
                  :style="{ width: '90%', maxWidth: '400px' }"
                  header-style="padding:12px;"
                  content-style="margin:0 28px;"
                  size="huge"
                  :bordered="false"
                  :segmented="{ content: 'soft', footer: 'soft' }"
                >
                  <template #header-extra>
                    <n-tag class="mr-2" type="primary" round>
                      Lv.{{ palDetail.level }}
                    </n-tag>
                    <n-tag
                      class="mr-3"
                      :type="palDetail.gender === 'Male' ? 'primary' : 'error'"
                      round
                    >
                      {{ palDetail.gender === "Male" ? "♂" : "♀" }}
                    </n-tag>
                  </template>
                  <template #header>
                    {{
                      locale === "zh"
                        ? palZHTypes[palDetail.type]
                          ? palZHTypes[palDetail.type]
                          : palDetail.type
                        : locale === "ja"
                          ? palJATypes[palDetail.type]
                            ? palJATypes[palDetail.type]
                            : palDetail.type
                          : palDetail.type
                    }}
                  </template>
                  <n-space class="mb-2" justify="center">
                    <n-avatar
                      :size="64"
                      :src="getPalAvatar(palDetail.type)"
                      :fallback-src="getUnknowPalAvatar()"
                    ></n-avatar>
                  </n-space>
                  <n-space class="mb-2" justify="center">
                    <n-tag v-if="palDetail.is_boss" type="success" round
                      >Boss</n-tag
                    >
                    <n-tag
                      v-else-if="palDetail.is_lucky"
                      type="warning"
                      round
                      >{{ $t("pal.lucky") }}</n-tag
                    >
                    <n-tag v-else-if="palDetail.is_tower" type="error" round>{{
                      $t("pal.tower")
                    }}</n-tag>
                  </n-space>
                  <n-space vertical>
                    <n-progress
                      type="line"
                      status="error"
                      indicator-placement="inside"
                      :percentage="percentageHP(palDetail.hp, palDetail.max_hp)"
                      :height="24"
                      :border-radius="4"
                      :fill-border-radius="0"
                      >HP:
                      {{
                        displayHP(palDetail.hp, palDetail.max_hp)
                      }}</n-progress
                    >
                    <n-grid cols="4">
                      <!-- <n-gi>
                          <n-statistic label="Exp" :value="palDetail.exp" />
                        </n-gi> -->
                      <n-gi>
                        <n-statistic
                          :label="$t('pal.ranged')"
                          :value="palDetail.ranged"
                        />
                      </n-gi>
                      <n-gi>
                        <n-statistic
                          :label="$t('pal.defense')"
                          :value="palDetail.defense"
                        />
                      </n-gi>
                      <n-gi>
                        <n-statistic
                          :label="$t('pal.melee')"
                          :value="palDetail.melee"
                        />
                      </n-gi>
                      <n-gi>
                        <n-statistic
                          :label="$t('pal.rank')"
                          :value="palDetail.rank"
                        />
                      </n-gi>
                    </n-grid>
                  </n-space>
                  <n-space vertical>
                    <div v-for="skill in palDetail.skills" :key="skill">
                      <n-tag type="warning">{{ skill }}</n-tag>
                      :
                      {{
                        skillDescMap[locale][skill]
                          ? skillDescMap[locale][skill]
                          : "Unknown"
                      }}
                    </div>
                  </n-space>
                </n-modal>
              </n-layout>

              <n-layout
                v-if="currentDisplay === 'guilds'"
                :native-scrollbar="false"
              >
                <n-card
                  :bordered="false"
                  v-if="guildInfo.name"
                  content-style="padding:0;"
                >
                  <n-page-header class="px-3 pt-3">
                    <template #title>
                      <span class="text-sm">{{ guildInfo.name }}</span>
                    </template>
                    <template #avatar>
                      <n-avatar
                        :style="{
                          color: 'white',
                          backgroundColor: 'darkorange',
                        }"
                        round
                      >
                        <n-icon>
                          <GroupWorkRound />
                        </n-icon>
                      </n-avatar>
                    </template>
                    <template #extra>
                      <n-space>
                        <n-tag type="primary" size="small" round strong>
                          Lv.{{ guildInfo.base_camp_level }}
                          <template #icon>
                            <n-icon :component="CrownFilled" />
                          </template>
                        </n-tag>
                      </n-space>
                    </template>
                    <template #footer> </template>
                  </n-page-header>
                  <n-space vertical>
                    <n-list hoverable clickable>
                      <n-list-item
                        v-for="player in guildInfo.players"
                        :key="player.player_uid"
                      >
                        <n-space size="large" style="margin-top: 4px">
                          <n-avatar :src="getUserAvatar()" round></n-avatar>
                          {{ player.nickname }}
                          <n-tag :bordered="false" type="info" size="small">
                            UID: {{ player.player_uid }}
                          </n-tag>
                          <n-tag
                            v-if="
                              player.player_uid === guildInfo.admin_player_uid
                            "
                            type="error"
                            size="small"
                          >
                            {{ $t("status.master") }}
                          </n-tag>
                        </n-space>
                      </n-list-item></n-list
                    >
                  </n-space>
                </n-card>
              </n-layout>
            </div>
          </n-layout>
        </n-layout>
      </div>
    </div>
  </div>
  <!-- login modal -->
  <n-modal
    v-model:show="showLoginModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 600px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.auth')"
    size="huge"
    :bordered="false"
    :segmented="segmented"
  >
    <div>
      <span class="block pb-2">{{ $t("message.authdesc") }}</span>
      <n-input
        type="password"
        show-password-on="click"
        size="large"
        v-model:value="password"
      ></n-input>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <n-button
          type="tertiary"
          @click="
            () => {
              showLoginModal = false;
              password = '';
            }
          "
          >{{ $t("button.cancel") }}</n-button
        >
        <n-button class="ml-3 w-40" type="primary" @click="handleLogin">{{
          $t("button.confirm")
        }}</n-button>
      </div>
    </template>
  </n-modal>
  <!-- broadcast modal -->
  <n-modal
    v-model:show="showBroadcastModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 600px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.broadcast')"
    size="huge"
    :bordered="false"
    :segmented="segmented"
  >
    <div>
      <n-input
        type="text"
        show-password-on="click"
        v-model:value="broadcastText"
      ></n-input>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <n-button
          type="tertiary"
          @click="
            () => {
              showBroadcastModal = false;
              broadcastText = '';
            }
          "
          >{{ $t("button.cancel") }}</n-button
        >
        <n-button class="ml-3 w-40" type="primary" @click="handleBroadcast">{{
          $t("button.confirm")
        }}</n-button>
      </div>
    </template>
  </n-modal>
</template>
<style scoped lang="less">
:deep .n-layout-scroll-container {
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
