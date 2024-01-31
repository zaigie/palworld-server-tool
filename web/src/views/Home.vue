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
const avatar =
  "data:image/webp;base64,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";

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
const getPlayerList = async () => {
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = data.value;
  await getPlayerInfo(playerList.value[0].player_uid);
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
        return palZHSkills[skill];
      });
    });
  }
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
  return dayjs() - dayjs(last_online) < 60000;
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
const handleLogin = () => {
  const context = new ApiService()
    .login({
      password: password.value,
    })
    .then((res) => {
      if (res.statusCode.value === 401) {
        message.error(t("message.autherr"));
        password.value = "";
        return;
      }
      let token = res.data.value.token;
      localStorage.setItem(PALWORLD_TOKEN, token);
      authToken.value = token;
      message.success(t("message.authsuccess"));
      showLoginModal.value = false;
      isLogin.value = true;
    });
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
    onPositiveClick: () => {
      if (type === "ban") {
        new ApiService()
          .banPlayer({
            playerUid: playerInfo.value.player_uid,
          })
          .then((res) => {
            if (res.statusCode.value === 200) {
              message.success(t("message.bansuccess"));
            } else {
              message.error(t("message.banfail"));
            }
          });
      } else if (type === "kick") {
        new ApiService()
          .kickPlayer({
            playerUid: playerInfo.value.player_uid,
          })
          .then((res) => {
            if (res.statusCode.value === 200) {
              message.success(t("message.kicksuccess"));
            } else {
              console.log(res);
              message.error(t("message.kickfail"));
            }
          });
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
  const { data } = await new ApiService()
    .sendBroadcast({
      message: broadcastText.value,
    })
    .then((res) => {
      if (res.statusCode.value === 200) {
        message.success(t("message.broadcastsuccess"));
        showBroadcastModal.value = false;
        bradcastText.value = "";
      } else {
        message.error(t("message.broadcastfail"));
      }
    });
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
      onPositiveClick: () => {
        doShutdown().then((res) => {
          if (res.statusCode.value === 200) {
            message.success(t("message.shutdownsuccess"));
            return;
          } else {
            message.error(t("message.shutdownfail"));
          }
        });
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
    getPlayerList();
  }, 60000);
});
</script>

<template>
  <div class="home-page">
    <div class="bg-#fff flex justify-between items-center p-3">
      <n-space>
        <span
          class="line-clamp-1"
          :class="smallScreen ? 'text-lg' : 'text-2xl'"
          >{{ $t("title") }}</span
        >
        <n-tag type="default" size="large">{{
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
            <n-button-group size="large">
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
                size="large"
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
                size="large"
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
              width="450"
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
                    <n-avatar :src="avatar" round></n-avatar>
                  </template>
                  <div>
                    <div>
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
                      <span class="pl-2 font-bold">{{ player.nickname }}</span>
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
              content-style="padding: 24px;"
              :native-scrollbar="false"
            >
              <n-card :bordered="false" v-if="playerInfo.nickname">
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
                    <n-avatar :src="avatar" round></n-avatar>
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
                    : {{ skillDescMap[locale][skill] }}
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
                        <n-avatar :src="avatar" round></n-avatar>
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
      <span class="block pb-2">You need to log in before you can operate.</span>
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
