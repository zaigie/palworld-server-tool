<script setup>
import {
  AdminPanelSettingsOutlined,
  SupervisedUserCircleRound,
} from "@vicons/material";
import { ChevronsLeft } from "@vicons/tabler";
import { GameController, LanguageSharp } from "@vicons/ionicons5";
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import palMap from "@/assets/pal.json";
import PlayerList from "./component/PlayerList.vue";
import GuildList from "./component/GuildList.vue";
import PlayerDetail from "./component/PlayerDetail.vue";
import GuildDetail from "./component/GuildDetail.vue";
import userStore from "@/stores/model/user";
import AdminOverview from "@/components/AdminOverview.vue";
import BackupManager from "@/components/BackupManager.vue";
import BroadcastComposer from "@/components/BroadcastComposer.vue";
import RconManager from "@/components/RconManager.vue";
import ShutdownDialog from "@/components/ShutdownDialog.vue";
import WhitelistManager from "@/components/WhitelistManager.vue";
import MapView from "@/views/PcHome/component/MapView.vue";

const emit = defineEmits(["open-config"]);

const { t, locale } = useI18n();

const message = useMessage();

const PALWORLD_TOKEN = "palworld_token";

const loading = ref(false);
const serverInfo = ref({});
const serverMetrics = ref({});
const localeLowerPalMap = ref({});
const currentDisplay = ref("players");
const isShowDetail = ref(false);
const playerList = ref([]);
const onlinePlayerList = ref([]);
const guildList = ref([]);
const playerInfo = ref({});
const playerPalsList = ref([]);
const currentPlayerPalsList = ref([]);
const guildInfo = ref({});
const languageOptions = ref([]);

const contentRef = ref(null);

const isLogin = ref(false);
const authToken = ref("");
let refreshTimer = null;
let mediaQuery = null;
const asArray = (value) => (Array.isArray(value) ? value : []);

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);

const updateDarkMode = (e) => {
  isDarkMode.value = e.matches;
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

// get data
const getServerInfo = async () => {
  const { data } = await new ApiService().getServerInfo();
  serverInfo.value = data.value || {};
};
const getServerMetrics = async () => {
  const { data } = await new ApiService().getServerMetrics();
  serverMetrics.value = data.value || {};
};
const getPlayerList = async (is_update_info = true) => {
  getOnlineList();
  const { data } = await new ApiService().getPlayerList({
    order_by: "last_online",
    desc: true,
  });
  playerList.value = asArray(data.value);
};
const getGuildList = async () => {
  const { data } = await new ApiService().getGuildList();
  guildList.value = asArray(data.value);
};

const getPlayerInfo = async (player_uid) => {
  const { data } = await new ApiService().getPlayer({ playerUid: player_uid });
  playerInfo.value = data.value;
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

// 接受子组件
const getChoosePlayer = (uid) => {
  getPlayerInfo(uid);
};
const getChooseGuild = (uid) => {
  getGuildInfo(uid);
};

const viewPlayerFromGuild = async (uid) => {
  currentDisplay.value = "players";
  await getPlayerInfo(uid);
};

const getPalName = (name) => {
  const lowerName = name.toLowerCase();
  return localeLowerPalMap.value[lowerName]
    ? localeLowerPalMap.value[lowerName]
    : name;
};

// 游戏用户的帕鲁列表分页，搜索
const clickSearch = (searchValue) => {
  const pattern = /^\s*$|(\s)\1/;
  if (searchValue && !pattern.test(searchValue)) {
    playerPalsList.value = playerInfo.value.pals.filter((item) => {
      return (
        item.skills.some((skill) => {
          return (
            skillMap[locale.value][skill]
              ? skillMap[locale.value][skill].name
              : skill
          ).includes(searchValue);
        }) || getPalName(item.type).includes(searchValue)
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
      pageSize.value * currentPage.value,
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

const getOnlineList = async () => {
  const { data } = await new ApiService().getOnlinePlayerList();
  onlinePlayerList.value = asArray(data.value);
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
  userStore().setIsLogin(true, token);
  authToken.value = token;
  message.success(t("message.authsuccess"));
  showLoginModal.value = false;
  isLogin.value = true;
  currentDisplay.value = "overview";
};

// broadcast
const showBroadcastModal = ref(false);
const handleStartBrodcast = () => {
  // broadcast start
  if (checkAuthToken()) {
    showBroadcastModal.value = true;
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};
const showShutdownDialog = ref(false);
const showRconDrawer = ref(false);
const showBackupManager = ref(false);
const showWhitelistManager = ref(false);
const handleShutdown = () => {
  if (checkAuthToken()) {
    showShutdownDialog.value = true;
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};

const openAuthenticated = (target) => {
  if (checkAuthToken()) {
    target.value = true;
  } else {
    message.error(t("message.requireauth"));
    showLoginModal.value = true;
  }
};

const adminOptions = computed(() => [
  { label: t("configuration.title"), key: "settings" },
  { label: t("button.rcon"), key: "rcon" },
  { label: t("button.backup"), key: "backup" },
  { label: t("button.whitelist"), key: "whitelist" },
  { label: t("button.broadcast"), key: "broadcast" },
  { label: t("button.palconf"), key: "palconf" },
  { type: "divider", key: "divider" },
  {
    label: t("button.shutdown"),
    key: "shutdown",
    props: { style: "color: #d03050" },
  },
]);

const handleAdminAction = (key) => {
  if (key === "settings") {
    if (checkAuthToken()) emit("open-config");
    else {
      message.error(t("message.requireauth"));
      showLoginModal.value = true;
    }
  }
  if (key === "rcon") openAuthenticated(showRconDrawer);
  if (key === "backup") openAuthenticated(showBackupManager);
  if (key === "whitelist") openAuthenticated(showWhitelistManager);
  if (key === "broadcast") handleStartBrodcast();
  if (key === "shutdown") handleShutdown();
  if (key === "config") {
    if (checkAuthToken()) emit("open-config");
    else {
      message.error(t("message.requireauth"));
      showLoginModal.value = true;
    }
  }
  if (key === "palconf") {
    window.open(
      "https://pal-conf.bluefissure.com/",
      "_blank",
      "noopener,noreferrer",
    );
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

  contentRef.value.scrollTo(0, 0);
};
const toOverview = () => {
  currentDisplay.value = "overview";
  isShowDetail.value = false;
};
const toMap = () => {
  currentDisplay.value = "map";
  isShowDetail.value = false;
};
const returnList = () => {
  isShowDetail.value = false;

  palsLoading.value = false;
  finished.value = false;
  currentPage.value = 1;

  contentRef.value.scrollTo(0, 0);
};

/**
 * check auth token
 */
const checkAuthToken = () => {
  const token = localStorage.getItem(PALWORLD_TOKEN);
  if (token && token !== "") {
    if (isTokenExpired(token)) {
      localStorage.removeItem(PALWORLD_TOKEN);
      return false;
    }
    isLogin.value = true;
    authToken.value = token;
    return true;
  }
  return false;
};
const isTokenExpired = (token) => {
  const payload = JSON.parse(atob(token.split(".")[1]));
  return payload.exp < Date.now() / 1000;
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
  localeLowerPalMap.value = Object.keys(palMap[locale.value]).reduce(
    (acc, key) => {
      acc[key.toLowerCase()] = palMap[locale.value][key];
      return acc;
    },
    {},
  );
  mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  mediaQuery.addEventListener("change", updateDarkMode);
  isDarkMode.value = mediaQuery.matches;

  loading.value = true;
  checkAuthToken();
  await Promise.all([getServerInfo(), getServerMetrics(), getPlayerList()]);
  if (isLogin.value) currentDisplay.value = "overview";
  loading.value = false;
  refreshTimer = setInterval(() => {
    getPlayerList(false);
    getServerMetrics();
  }, 60000);
});

onBeforeUnmount(() => {
  clearInterval(refreshTimer);
  mediaQuery?.removeEventListener("change", updateDarkMode);
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
            : loading
              ? $t("message.loading")
              : $t("status.serverUnavailable")
        }}</n-tag>
      </div>
      <n-space vertical>
        <n-space justify="end">
          <n-tag type="info" round size="small">{{
            $t("status.player_number", { number: playerList?.length })
          }}</n-tag>
          <n-tag type="success" round size="small">{{
            $t("status.online_number", {
              number: onlinePlayerList?.length ?? 0,
            })
          }}</n-tag>
        </n-space>
        <n-space justify="end" class="flex items-center">
          <n-dropdown
            trigger="hover"
            :options="languageOptions"
            @select="handleSelectLanguage"
          >
            <n-button
              type="default"
              secondary
              strong
              circle
              size="small"
              :aria-label="$t('button.language')"
            >
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
      <div class="rounded-lg" v-if="!loading">
        <n-layout style="height: calc(100vh - 86px)" has-sider>
          <n-layout-header
            class="flex flex-col justify-between"
            :class="isLogin ? 'h-16' : 'h-10'"
            bordered
          >
            <div v-if="isLogin" class="flex justify-center items-center px-3">
              <n-dropdown
                trigger="click"
                :options="adminOptions"
                @select="handleAdminAction"
              >
                <n-button size="small" type="primary" secondary strong round>
                  {{ $t("button.management") }}
                </n-button>
              </n-dropdown>
            </div>
            <div v-else></div>
            <div class="flex justify-end">
              <n-button-group size="small" class="w-full">
                <n-button
                  v-if="isShowDetail"
                  class="flex-1"
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
                  v-if="isLogin && !isShowDetail"
                  class="flex-1"
                  @click="toOverview"
                  :type="currentDisplay === 'overview' ? 'primary' : 'tertiary'"
                  secondary
                  strong
                >
                  {{ $t("button.overview") }}
                </n-button>
                <n-button
                  class="flex-1"
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
                  class="flex-1"
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
                <n-button
                  v-if="!isShowDetail"
                  class="flex-1"
                  @click="toMap"
                  :type="currentDisplay === 'map' ? 'primary' : 'tertiary'"
                  secondary
                  strong
                >
                  {{ $t("button.map") }}
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
            <admin-overview
              v-if="currentDisplay === 'overview'"
              :server-info="serverInfo"
              :server-metrics="serverMetrics"
              :players="playerList"
              @open-rcon="openAuthenticated(showRconDrawer)"
              @open-backup="openAuthenticated(showBackupManager)"
              @open-broadcast="handleStartBrodcast"
              @open-config="handleAdminAction('config')"
            />
            <map-view v-if="currentDisplay === 'map'" />
            <div v-if="!isShowDetail">
              <!-- list -->
              <player-list
                v-if="currentDisplay === 'players'"
                :playerList="playerList"
                @onGetInfo="getChoosePlayer"
              ></player-list>
              <guild-list
                v-if="currentDisplay === 'guilds'"
                :guildList="guildList"
                @onGetInfo="getChooseGuild"
              >
              </guild-list>
            </div>
            <!-- detail -->
            <div v-else class="relative">
              <player-detail
                v-if="currentDisplay === 'players'"
                :playerInfo="playerInfo"
                :currentPlayerPalsList="currentPlayerPalsList"
                :finished="finished"
                @onSearch="clickSearch"
              ></player-detail>
              <guild-detail
                v-if="currentDisplay === 'guilds'"
                :guildInfo="guildInfo"
                @view-player="viewPlayerFromGuild"
              ></guild-detail>
            </div>
          </n-layout>
        </n-layout>
      </div>
    </div>
  </div>
  <!-- 登录 modal -->
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
  >
    <div>
      <span class="block pb-2">{{ $t("message.authdesc") }}</span>
      <n-input
        type="password"
        show-password-on="click"
        size="large"
        v-model:value="password"
        :aria-label="$t('modal.auth')"
        autocomplete="current-password"
        @keyup.enter="handleLogin"
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
  <rcon-manager v-model:show="showRconDrawer" />
  <broadcast-composer v-model:show="showBroadcastModal" />
  <shutdown-dialog v-model:show="showShutdownDialog" />
  <backup-manager v-model:show="showBackupManager" />
  <whitelist-manager
    v-model:show="showWhitelistManager"
    :players="playerList"
  />
</template>
<style scoped lang="less">
:deep .n-layout-scroll-container {
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
