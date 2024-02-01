<script setup>
import {
  AdminPanelSettingsOutlined,
  SupervisedUserCircleRound,
  GroupWorkRound,
  ContentCopyFilled,
  SettingsPowerRound,
} from "@vicons/material";
import { GameController, LogOut, Ban, LanguageSharp } from "@vicons/ionicons5";
import { BroadcastTower } from "@vicons/fa";
import { CrownFilled } from "@vicons/antd";
import { computed, onMounted, ref, h } from "vue";
import { NTag, NButton, NAvatar, useMessage, useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import pageStore from "@/stores/model/page.js";
import dayjs from "dayjs";
import skillDescMap from "@/assets/skillDesc.json";
import palZHTypes from "@/assets/zhTypes.json";
import palZHSkills from "@/assets/zhSkills.json";

const { t, locale } = useI18n();

const message = useMessage();
const dialog = useDialog();

const PALWORLD_TOKEN = "palworld_token";

const pageWidth = computed(() => pageStore().getScreenWidth());
const smallScreen = computed(() => pageWidth.value < 1024);

const loading = ref(false);
const serverInfo = ref({});
const currentDisplay = ref("players");
const playerList = ref([]);
const guildList = ref([]);
const playerInfo = ref({});
const guildInfo = ref({});
const skillTypeList = ref([]);
const languageOptions = ref([]);

const isLogin = ref(false);
const authToken = ref("");

const getUserAvatar = () => {
  return new URL("../assets/avatar.webp", import.meta.url).href;
};

const handleSelectLanguage = (key) => {
  message.info(t("message.changelanguage"));
  if (key === "zh") {
    localStorage.setItem("locale", "zh");
    // locale.value = "zh";
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
  if (is_update_info && playerList.value.length > 0) {
    getPlayerInfo(playerList.value[0].player_uid);
  }
};
const getGuildList = async () => {
  const { data } = await new ApiService().getGuildList();
  guildList.value = data.value;
  if (guildList.value.length > 0) {
    getGuildInfo(guildList.value[0].admin_player_uid);
  }
};

const getPlayerInfo = async (player_uid) => {
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  playerInfo.value = data.value;
  if (locale.value === "zh") {
    playerInfo.value.pals.forEach((pal) => {
      pal.skills = pal.skills.map((skill) => {
        return palZHSkills[skill] ? palZHSkills[skill] : skill;
      });
    });
  }
  nextTick(() => {
    const playerInfoEL = document.getElementById("player-info");
    if (playerInfoEL) {
      playerInfoEL.scrollIntoView({ behavior: "smooth" });
    }
  });
};

const getGuildInfo = async (admin_player_uid) => {
  const { data } = await new ApiService().getGuild({
    adminPlayerUid: admin_player_uid,
  });
  guildInfo.value = data.value;
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

const displayHP = (hp, max_hp) => {
  return (hp / 1000).toFixed(0) + "/" + (max_hp / 1000).toFixed(0);
};

const percentageHP = (hp, max_hp) => {
  if (max_hp === 0) {
    return 0;
  }
  return ((hp / max_hp) * 100).toFixed(2);
};

const createPlayerPalsColumns = () => {
  return [
    {
      title: "",
      key: "",
      render(row) {
        return h(NAvatar, {
          size: "small",
          src: getPalAvatar(row.type),
          fallbackSrc: getUnknowPalAvatar(),
        });
      },
    },
    {
      title: t("pal.type"),
      key: "type",
      // defaultSortOrder: 'ascend',
      sorter: "default",
      render(row) {
        return [
          h(
            NTag,
            {
              style: {
                marginRight: "6px",
              },
              type: row.gender == "Male" ? "primary" : "error",
              bordered: false,
            },
            {
              default: () => (row.gender == "Male" ? "♂" : "♀"),
            }
          ),
          h(
            "div",
            {
              style: {
                display: "inline-block",
                color: row.is_lucky ? "darkorange" : "black",
                fontWeight: row.is_lucky ? "bold" : "normal",
              },
            },
            {
              default: () =>
                locale.value === "zh"
                  ? palZHTypes[row.type]
                    ? palZHTypes[row.type]
                    : row.type
                  : row.type,
            }
          ),
        ];
      },
    },
    {
      title: t("pal.level"),
      key: "level",
      width: 70,
      defaultSortOrder: "descend",
      sorter: "default",
      render(row) {
        return "Lv." + row.level;
      },
    },
    {
      title: t("pal.skills"),
      key: "skills",
      render(row) {
        const skills = row.skills.map((skill) => {
          return h(
            NTag,
            {
              style: {
                marginRight: "6px",
              },
              type: "warning",
              bordered: false,
            },
            {
              default: () => skill,
            }
          );
        });
        return skills;
      },
      filterOptions: skillTypeList.value.map((value) => ({
        label: value,
        value,
      })),
      filter(value, row) {
        return ~row.skills.indexOf(value);
      },
    },
    {
      title: "",
      key: "actions",
      render(row) {
        return h(
          NButton,
          {
            size: "small",
            onClick: () => showPalDetail(row),
          },
          { default: () => t("button.detail") }
        );
      },
    },
  ];
};

const copyText = async (text) => {
  if (!navigator.clipboard) {
    message.error(t("message.copyfail"));
    return;
  }

  try {
    await navigator.clipboard.writeText(text);
    message.success(t("message.copysuccess"));
  } catch (err) {
    message.error(t("message.copyerr", { err: err }));
  }
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
    message.error(t("message.broadcastfail", { err: data.value?.error }));
  }
};

const doShutdown = async () => {
  return await new ApiService().shutdownServer({
    seconds: 60,
    message: "Server Will Shutdown After 60 Seconds",
  });
};

// shutdown
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
};
const toGuilds = async () => {
  if (currentDisplay.value === "guilds") {
    return;
  }
  await getGuildList();
  currentDisplay.value = "guilds";
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
  ];
  skillTypeList.value = getSkillTypeList();
  loading.value = true;
  checkAuthToken();
  await getServerInfo();
  await getPlayerList();
  loading.value = false;
  setInterval(() => {
    getPlayerList(false);
  }, 60000);
});
</script>

<template>
  <div class="home-page">
    <div class="bg-#fff flex justify-between items-center p-3">
      <n-space class="flex items-center">
        <span
          class="line-clamp-1"
          :class="smallScreen ? 'text-lg' : 'text-2xl'"
          >{{ $t("title") }}</span
        >
        <n-tag type="default" :size="smallScreen ? 'medium' : 'large'">{{
          serverInfo.name + " " + serverInfo.version
        }}</n-tag>
      </n-space>

      <n-space>
        <n-dropdown
          trigger="hover"
          :options="languageOptions"
          @select="handleSelectLanguage"
        >
          <n-button type="default" secondary strong circle>
            <template #icon>
              <n-icon><LanguageSharp /></n-icon>
            </template>
          </n-button>
        </n-dropdown>

        <n-button
          type="primary"
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
        <n-tag v-else type="success" size="large" round>
          <template #icon>
            <n-icon>
              <AdminPanelSettingsOutlined />
            </n-icon>
          </template>
          {{ $t("status.authenticated") }}
        </n-tag>
      </n-space>
    </div>
    <div class="w-full">
      <div class="rounded-lg" v-if="!loading && playerList.length > 0">
        <n-layout style="height: calc(100vh - 64px)" has-sider>
          <n-layout-header class="p-3 flex justify-between h-16" bordered>
            <n-button-group :size="smallScreen ? 'medium' : 'large'">
              <n-button
                @click="toPlayers"
                :type="currentDisplay === 'players' ? 'primary' : 'tertiary'"
                secondary
                strong
                round
              >
                <template #icon>
                  <n-icon>
                    <GameController />
                  </n-icon>
                </template>
                {{ $t("button.players") }}
              </n-button>
              <n-button
                @click="toGuilds"
                :type="currentDisplay === 'guilds' ? 'primary' : 'tertiary'"
                secondary
                strong
                round
              >
                <template #icon>
                  <n-icon>
                    <SupervisedUserCircleRound />
                  </n-icon>
                </template>
                {{ $t("button.guilds") }}
              </n-button>
            </n-button-group>
            <n-space v-if="isLogin">
              <n-button
                :size="smallScreen ? 'medium' : 'large'"
                type="success"
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
                :size="smallScreen ? 'medium' : 'large'"
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
            </n-space>
          </n-layout-header>
          <n-layout position="absolute" style="top: 64px" has-sider>
            <n-layout-sider
              :width="smallScreen ? 360 : 400"
              content-style="padding: 24px;"
              :native-scrollbar="false"
              bordered
            >
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
                      class="inline-block mt-1 rounded-full bg-#ddd text-xs px-2 py-0.5"
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
            </n-layout-sider>
            <n-layout
              v-if="currentDisplay === 'players'"
              :native-scrollbar="false"
            >
              <n-card
                id="player-info"
                :bordered="false"
                v-if="playerInfo.nickname"
              >
                <n-page-header>
                  <n-grid :cols="6">
                    <n-gi
                      v-for="status in Object.entries(playerInfo.status_point)"
                      :key="status[0]"
                    >
                      <n-statistic :label="status[0]" :value="status[1]" />
                    </n-gi>
                  </n-grid>
                  <template #title>
                    {{ playerInfo.nickname }}
                    <n-tag
                      :bordered="false"
                      :type="
                        isPlayerOnline(playerInfo.last_online)
                          ? 'success'
                          : 'error'
                      "
                      round
                      >{{
                        isPlayerOnline(playerInfo.last_online)
                          ? $t("status.online")
                          : $t("status.offline")
                      }}</n-tag
                    >
                    <n-button
                      @click="copyText(playerInfo.player_uid)"
                      class="ml-3"
                      type="info"
                      size="small"
                      icon-placement="right"
                      ghost
                    >
                      UID: {{ playerInfo.player_uid }}
                      <template #icon>
                        <n-icon><ContentCopyFilled /></n-icon>
                      </template>
                    </n-button>
                  </template>
                  <template #avatar>
                    <n-avatar :src="getUserAvatar()" round></n-avatar>
                  </template>
                  <template #extra>
                    <n-space>
                      <n-tag type="primary" size="large" round strong>
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
                    :percentage="percentageHP(playerInfo.hp, playerInfo.max_hp)"
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
                      displayHP(playerInfo.shield_hp, playerInfo.shield_max_hp)
                    }}</n-progress
                  >
                </n-space>
                <n-data-table
                  class="mt-5"
                  size="small"
                  :columns="createPlayerPalsColumns()"
                  :row-props="dataRowProps"
                  :data="playerInfo.pals"
                  :bordered="false"
                  striped
                />
              </n-card>
              <n-modal
                v-model:show="showPalDetailModal"
                preset="card"
                :style="{ width: '90%', maxWidth: '400px' }"
                size="huge"
                :bordered="false"
                :segmented="{ content: 'soft', footer: 'soft' }"
              >
                <template #header-extra>
                  <n-tag class="mr-2" type="primary" round>
                    Lv.{{ palDetail.level }}
                  </n-tag>
                  <n-tag
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
                  <n-tag v-else-if="palDetail.is_lucky" type="warning" round>{{
                    $t("pal.lucky")
                  }}</n-tag>
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
                    {{ displayHP(palDetail.hp, palDetail.max_hp) }}</n-progress
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
              content-style="padding: 24px;"
              :native-scrollbar="false"
            >
              <n-card :bordered="false" v-if="guildInfo.name">
                <n-page-header>
                  <template #title>
                    {{ guildInfo.name }}
                  </template>
                  <template #avatar>
                    <n-avatar
                      :style="{ color: 'white', backgroundColor: 'darkorange' }"
                      round
                    >
                      <n-icon>
                        <GroupWorkRound />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #extra>
                    <n-space>
                      <n-tag type="primary" size="large" round strong>
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
          </n-layout>
          <n-layout-footer
            v-if="
              currentDisplay === 'players' && isLogin && playerInfo.player_uid
            "
            class="pt-3 pr-3 bg-transparent"
            position="absolute"
            style="height: 64px"
          >
            <n-flex justify="end">
              <n-button
                @click="handelPlayerAction('ban')"
                type="error"
                size="large"
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
                size="large"
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
          </n-layout-footer>
          <!-- <n-layout
            v-if="!loading && playerList.length === 0"
            class="w-full h-25 flex justify-center items-center"
          >
            <n-empty description="什么都没有"> </n-empty>
          </n-layout>
          <n-layout
            v-if="loading"
            class="w-full h-25 flex justify-center items-center"
          >
            <n-spin size="small" />
          </n-layout> -->
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
        type="textarea"
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
              showLoginModal = false;
              password = '';
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
