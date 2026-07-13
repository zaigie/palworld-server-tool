<script setup>
import { ContentCopyFilled, PersonSearchSharp } from "@vicons/material";
import { LogOut, Ban, ShieldCheckmarkOutline } from "@vicons/ionicons5";
import { CrownFilled } from "@vicons/antd";
import ApiService from "@/service/api";
import { ref, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import palMap from "@/assets/pal.json";
import skillMap from "@/assets/skill.json";
import { useDialog, useMessage, NAvatar, NTag, NButton } from "naive-ui";
import PalDetail from "./PalDetail.vue";
import whitelistStore from "@/stores/model/whitelist.js";
import playerToGuildStore from "@/stores/model/playerToGuild.js";
import userStore from "@/stores/model/user";
import palItems from "@/assets/items.json";
import {
  localizedSkillName,
  statusPointTranslationKey,
} from "@/utils/gameLabels";

const { t, locale } = useI18n();
const PALWORLD_TOKEN = "palworld_token";
const props = defineProps(["playerInfo", "playerPalsList"]);
const playerInfo = computed(() => props.playerInfo);
const playerPalsList = computed(() => props.playerPalsList);

const isLogin = computed(() => userStore().getLoginInfo().isLogin);

const message = useMessage();
const dialog = useDialog();

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);

const platformColors = {
  steam: { color: "#223D58", textColor: "#fff" },
  xbox: { color: "#2B8B2B", textColor: "#fff" },
  ps5: { color: "#00439C", textColor: "#fff" },
  mac: { color: "#777", textColor: "#fff" },
  default: { color: "#d9c36c", textColor: "#fff" },
};

const localeLowerPalMap = ref({});
const skillTypeList = ref([]);

// 帕鲁列表
const currentPalsList = ref([]);
const createPlayerPalsColumns = () => {
  return [
    {
      title: "",
      key: "",
      render(row) {
        return h(NAvatar, {
          size: "small",
          src: getPalAvatar(row.type),
          fallbackSrc: getUnknowPalAvatar(row.is_boss),
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
            },
          ),
          h(
            "div",
            {
              style: {
                display: "inline-block",
                color: row.is_lucky ? "darkorange" : getDarkModeColor(),
                fontWeight: row.is_lucky ? "bold" : "normal",
              },
            },
            {
              default: () => getPalName(row.type),
            },
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
              default: () => localizedSkillName(skill, locale.value, skillMap),
            },
          );
        });
        return skills;
      },
      filterOptions: skillTypeList.value.map((value) => ({
        label: value,
        value,
      })),
      filter(value, row) {
        return row.skills.some((skill) => {
          return localizedSkillName(skill, locale.value, skillMap).includes(
            value,
          );
        });
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
          { default: () => t("button.detail") },
        );
      },
    },
  ];
};

watch(
  () => playerPalsList.value,
  (newVal) => {
    currentPalsList.value = newVal;
    paginationReactive.page = 1;
    paginationReactive.pageSize = 10;
    searchValue.value = "";
    mergeItems();
  },
);

// 游戏用户的帕鲁列表分页，搜索
const paginationReactive = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 15, 20],
  onChange: (page) => {
    paginationReactive.page = page;
  },
  onUpdatePageSize: (pageSize) => {
    paginationReactive.pageSize = pageSize;
    paginationReactive.page = 1;
  },
});

const searchValue = ref("");
const clickSearch = () => {
  const pattern = /^\s*$|(\s)\1/;
  if (searchValue.value && !pattern.test(searchValue.value)) {
    currentPalsList.value = playerInfo?.value.pals.filter((item) => {
      return (
        item.skills.some((skill) => {
          return localizedSkillName(skill, locale.value, skillMap).includes(
            searchValue.value,
          );
        }) || getPalName(item.type).includes(searchValue.value)
      );
    });
  } else {
    currentPalsList.value = JSON.parse(JSON.stringify(playerInfo?.value.pals));
  }
  paginationReactive.page = 1;
};
const clearSearch = () => {
  nextTick(() => {
    clickSearch();
  });
};

// 帕鲁详情
const showPalDetailModal = ref(false);
const palDetail = ref({});

const showPalDetail = (pal) => {
  palDetail.value = pal;
  showPalDetailModal.value = true;
};

// UID、Steam64 复制
const copyText = async (text) => {
  if (navigator.clipboard) {
    try {
      await navigator.clipboard.writeText(text);
      message.success(t("message.copysuccess"));
    } catch (err) {
      message.error(t("message.copyerr", { err }));
    }
  } else {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand("copy");
      message.success(t("message.copysuccess"));
    } catch (err) {
      message.error(t("message.copyerr", { err }));
    }
    document.body.removeChild(textarea);
  }
};
// 查看公会
const toGuilds = async (uid) => {
  playerToGuildStore().setCurrentUid(uid);
  playerToGuildStore().setUpdateStatus("guilds");
};

// 加入白名单
const showAddWhiteListModal = ref(false);
const addWhiteData = ref({
  name: "",
  player_uid: "",
  steam_id: "",
});
const addWhiteList = async () => {
  const { data, statusCode } = await new ApiService().addWhitelist(
    addWhiteData,
  );
  if (statusCode.value === 200) {
    message.success(t("message.addwhitesuccess"));
    showAddWhiteListModal.value = false;
    await getWhiteList();
  } else {
    message.error(t("message.addwhitefail", { err: data.value?.error }));
  }
};
const handleAddWhiteList = () => {
  if (isLogin.value) {
    addWhiteData.value.name = playerInfo.value.nickname;
    addWhiteData.value.player_uid = playerInfo.value.player_uid;
    addWhiteData.value.steam_id = playerInfo.value.steam_id;
    showAddWhiteListModal.value = true;
  } else {
    message.error(t("message.requireauth"));
    showAddWhiteListModal.value = true;
  }
};
// 移除白名单
const removeWhitelist = async (player) => {
  if (isWhite(player)) {
    const { data, statusCode } = await new ApiService().removeWhitelist(player);
    if (statusCode.value === 200) {
      message.success(t("message.removewhitesuccess"));
      await getWhiteList();
    } else {
      message.error(t("message.removewhitefail", { err: data.value?.error }));
    }
  }
};

// 封禁、踢出
const handelPlayerAction = async (type) => {
  if (!isLogin.value) {
    message.error($t("message.requireauth"));
    showLoginModal.value = true;
    return;
  }
  const param = {
    ban: {
      title: t("message.bantitle"),
      content: t("message.banwarn"),
    },
    unban: {
      title: t("message.unbantitle"),
      content: t("message.unbanwarn"),
    },
    kick: {
      title: t("message.kicktitle"),
      content: t("message.kickwarn"),
    },
  }[type];
  dialog.warning({
    ...param,
    positiveText: t("button.confirm"),
    negativeText: t("button.cancel"),
    onPositiveClick: async () => {
      if (type === "ban") {
        const { data, statusCode } = await new ApiService().banPlayer({
          playerUid: playerInfo?.value.player_uid,
        });
        if (statusCode.value === 200) {
          message.success(t("message.bansuccess"));
        } else {
          message.error(t("message.banfail", { err: data.value?.error }));
        }
      } else if (type === "unban") {
        const { data, statusCode } = await new ApiService().unbanPlayer({
          playerUid: playerInfo?.value.player_uid,
        });
        if (statusCode.value === 200) {
          message.success(t("message.unbansuccess"));
        } else {
          message.error(t("message.unbanfail", { err: data.value?.error }));
        }
      } else if (type === "kick") {
        const { data, statusCode } = await new ApiService().kickPlayer({
          playerUid: playerInfo?.value.player_uid,
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

// 获取白名单列表
const whiteList = computed(() => whitelistStore().getWhitelist());
const getWhiteList = async () => {
  if (isLogin.value) {
    const { data, statusCode } = await new ApiService().getWhitelist();
    if (statusCode.value === 200) {
      if (data.value) {
        whitelistStore().setWhitelist(data.value);
      } else {
        whitelistStore().setWhitelist([]);
      }
    }
  }
};

// 是否在白名单中
const isWhite = (player) => {
  if (whiteList.value.length === 0) {
    return false;
  }
  return whiteList.value.some((whitelistItem) => {
    return (
      (whitelistItem.player_uid &&
        whitelistItem.player_uid === player.player_uid) ||
      (whitelistItem.steam_id && whitelistItem.steam_id === player.steam_id)
    );
  });
};

onMounted(async () => {
  skillTypeList.value = getSkillTypeList();
  await getWhiteList();
  localeLowerPalMap.value = Object.keys(palMap[locale.value]).reduce(
    (acc, key) => {
      acc[key.toLowerCase()] = palMap[locale.value][key];
      return acc;
    },
    {},
  );
});

// 其他操作
const getDarkModeColor = () => {
  return isDarkMode.value ? "#fff" : "#000";
};

const getSkillTypeList = () => {
  if (skillMap[locale.value]) {
    return Object.values(skillMap[locale.value]).map((item) => item.name);
  } else {
    return [];
  }
};
const getStatusPointLabel = (rawKey) => {
  const translationKey = statusPointTranslationKey(rawKey);
  return translationKey ? t(`statusPoint.${translationKey}`) : rawKey;
};
const getPalAvatar = (name) => {
  const lowerName = name.toLowerCase();
  return new URL(`../../../assets/pals/${lowerName}.png`, import.meta.url).href;
};
const getPalName = (name) => {
  const lowerName = name.toLowerCase();
  return localeLowerPalMap.value[lowerName]
    ? localeLowerPalMap.value[lowerName]
    : name;
};
const getItemIcon = (id) => {
  return new URL(`../../../assets/items/${id}.webp`, import.meta.url).href;
};
const getUnknowPalAvatar = (is_boss = false) => {
  if (is_boss) {
    return new URL("@/assets/pals/boss_unknown.png", import.meta.url).href;
  }
  return new URL("@/assets/pals/unknown.png", import.meta.url).href;
};
const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 80000;
};
const getPlatformColor = (userId) => {
  if (!userId) return platformColors.default;
  return platformColors[userId.split("_")[0]] || platformColors.default;
};
const displayLastOnline = (lastOnline) => {
  if (dayjs(lastOnline).year() < 1970) return "Unknown";
  return dayjs(lastOnline).format("YYYY-MM-DD HH:mm:ss");
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

const mergedItems = ref({});
const mergeItems = () => {
  mergedItems.value = {};

  if (!playerInfo.value.items) return;
  for (const [containerId, items] of Object.entries(playerInfo.value.items)) {
    mergedItems.value[containerId] = items.map((item) => {
      const frontendItem = palItems[locale.value].find(
        (frontItem) => frontItem.id === item.ItemId,
      );
      if (!frontendItem) {
        return {
          ...item,
          id: item.ItemId,
          name: item.ItemId,
          description: "No description.",
          key: item.ItemId,
        };
      }
      return {
        ...item,
        id: frontendItem.id,
        name: frontendItem.name,
        description: frontendItem.description,
        key: frontendItem.key,
      };
    });
  }
};

const createPlayerItemsColumns = () => {
  return [
    {
      title: "",
      key: "",
      render(row) {
        return h(NAvatar, {
          size: "small",
          src: getItemIcon(row.id),
          fallbackSrc: getUnknowPalAvatar(),
        });
      },
    },
    {
      title: t("item.name"),
      key: "name",
    },
    {
      title: t("item.description"),
      key: "description",
      defaultSortOrder: "descend",
    },
    {
      title: t("item.count"),
      key: "StackCount",
      width: 170,
      defaultSortOrder: "descend",
      sorter: "default",
    },
  ];
};
</script>

<template>
  <div class="player-detail" :class="{ 'is-dark': isDarkMode }">
    <n-card
      content-style="padding: 24px 28px 36px;"
      id="player-info"
      :bordered="false"
      v-if="playerInfo?.nickname"
    >
      <section class="player-overview" aria-labelledby="player-detail-title">
        <div class="overview-heading">
          <div class="player-identity">
            <div class="player-title-row">
              <h1 id="player-detail-title" class="player-title">
                {{ playerInfo?.nickname }}
              </h1>
              <n-tag type="primary" size="large" round strong>
                Lv.{{ playerInfo?.level }}
                <template #icon>
                  <n-icon :component="CrownFilled" />
                </template>
              </n-tag>
            </div>
            <div class="identity-tags">
              <n-tag
                :bordered="false"
                :type="
                  isPlayerOnline(playerInfo?.last_online) ? 'success' : 'error'
                "
                size="small"
                round
              >
                {{
                  isPlayerOnline(playerInfo?.last_online)
                    ? $t("status.online")
                    : $t("status.offline")
                }}
              </n-tag>
              <n-tag
                v-if="playerInfo?.user_id"
                :bordered="false"
                round
                size="small"
                :color="getPlatformColor(playerInfo.user_id)"
              >
                {{ playerInfo.user_id.split("_")[0] }}
              </n-tag>
              <n-tag
                v-if="isWhite(playerInfo)"
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
              <span class="last-online-text">
                {{ $t("status.last_online") }}
                {{ displayLastOnline(playerInfo?.last_online) }}
              </span>
            </div>
          </div>
          <n-button
            @click="toGuilds(playerInfo?.player_uid)"
            type="warning"
            secondary
            strong
          >
            {{ $t("button.viewGuild") }}
            <template #icon>
              <n-icon><PersonSearchSharp /></n-icon>
            </template>
          </n-button>
        </div>

        <div v-if="isLogin" class="identity-copy-list">
          <n-button
            class="identity-copy"
            secondary
            @click="copyText(playerInfo?.player_uid)"
          >
            <span class="copy-label">UID</span>
            <span class="copy-value">{{ playerInfo?.player_uid }}</span>
            <template #icon>
              <n-icon><ContentCopyFilled /></n-icon>
            </template>
          </n-button>
          <n-button
            class="identity-copy"
            secondary
            @click="copyText(playerInfo?.steam_id)"
          >
            <span class="copy-label">Steam64</span>
            <span class="copy-value">
              {{ playerInfo.steam_id ? playerInfo.steam_id : "--" }}
            </span>
            <template #icon>
              <n-icon><ContentCopyFilled /></n-icon>
            </template>
          </n-button>
        </div>

        <div
          v-if="
            playerInfo.ip ||
            playerInfo.ping ||
            playerInfo.location_x ||
            playerInfo.location_y
          "
          class="runtime-meta"
        >
          <span v-if="playerInfo.ip">IP {{ playerInfo.ip }}</span>
          <span v-if="playerInfo.ping">
            Ping {{ playerInfo.ping.toFixed(2) }}
          </span>
          <span v-if="playerInfo.location_x"
            >X {{ playerInfo.location_x }}</span
          >
          <span v-if="playerInfo.location_y"
            >Y {{ playerInfo.location_y }}</span
          >
        </div>

        <div class="status-grid">
          <div
            v-for="status in Object.entries(playerInfo?.status_point || {})"
            :key="status[0]"
            class="status-card"
          >
            <span class="status-label">{{
              getStatusPointLabel(status[0])
            }}</span>
            <strong class="status-value">{{ status[1] }}</strong>
          </div>
        </div>
      </section>
      <!-- <n-space vertical>
        <n-progress
          type="line"
          status="error"
          indicator-placement="inside"
          :percentage="percentageHP(playerInfo?.hp, playerInfo?.max_hp)"
          :height="24"
          :border-radius="4"
          :fill-border-radius="0"
          >HP: {{ displayHP(playerInfo?.hp, playerInfo?.max_hp) }}</n-progress
        >
        <n-progress
          type="line"
          indicator-placement="inside"
          :percentage="
            percentageHP(playerInfo?.shield_hp, playerInfo?.shield_max_hp)
          "
          :height="24"
          :border-radius="4"
          :fill-border-radius="0"
          >SHIELD:
          {{
            displayHP(playerInfo?.shield_hp, playerInfo?.shield_max_hp)
          }}</n-progress
        >
      </n-space> -->
      <div class="detail-tabs">
        <n-tabs type="line" size="large" animated>
          <n-tab-pane :name="$t('item.palList')">
            <div class="w-full mt-5">
              <n-input-group class="w-full flex justify-end">
                <n-input
                  v-model:value="searchValue"
                  clearable
                  :placeholder="$t('input.searchPlaceholder')"
                  :on-clear="clearSearch"
                  @keydown.enter="clickSearch"
                />
                <n-button type="primary" class="w-20" @click="clickSearch">
                  {{ $t("button.search") }}
                </n-button>
              </n-input-group>
            </div>
            <n-data-table
              class="mt-2"
              size="small"
              :columns="createPlayerPalsColumns()"
              :data="currentPalsList"
              :bordered="false"
              striped
              :pagination="paginationReactive"
            />
          </n-tab-pane>
          <n-tab-pane :name="$t('item.itemList')">
            <n-tabs type="segment" animated>
              <n-tab-pane :name="$t('item.commonContainer')">
                <n-data-table
                  size="small"
                  :columns="createPlayerItemsColumns()"
                  :data="mergedItems['CommonContainerId']"
                  :bordered="false"
                  striped
                  :pagination="paginationReactive"
                />
              </n-tab-pane>
              <n-tab-pane :name="$t('item.essentialContainer')">
                <n-data-table
                  size="small"
                  :columns="createPlayerItemsColumns()"
                  :data="mergedItems['EssentialContainerId']"
                  :bordered="false"
                  striped
                  :pagination="paginationReactive"
                />
              </n-tab-pane>
              <n-tab-pane :name="$t('item.weaponContainer')">
                <n-data-table
                  size="small"
                  :columns="createPlayerItemsColumns()"
                  :data="mergedItems['WeaponLoadOutContainerId']"
                  :bordered="false"
                  striped
                />
              </n-tab-pane>
              <n-tab-pane :name="$t('item.armorContainer')">
                <n-data-table
                  class="mt-1"
                  size="small"
                  :columns="createPlayerItemsColumns()"
                  :data="mergedItems['PlayerEquipArmorContainerId']"
                  :bordered="false"
                  striped
                />
              </n-tab-pane>
            </n-tabs>
          </n-tab-pane>
        </n-tabs>
      </div>
    </n-card>
    <!-- 加入白名单，封禁，踢出 -->
    <n-flex
      justify="end"
      class="player-actions"
      v-if="playerInfo?.nickname && isLogin"
    >
      <n-button
        @click="
          isWhite(playerInfo)
            ? removeWhitelist(playerInfo)
            : handleAddWhiteList()
        "
        :type="isWhite(playerInfo) ? 'warning' : 'success'"
        size="large"
        secondary
        strong
        round
      >
        <template #icon>
          <n-icon>
            <ShieldCheckmarkOutline />
          </n-icon>
        </template>
        {{
          isWhite(playerInfo)
            ? $t("button.removeWhitelist")
            : $t("button.joinWhitelist")
        }}
      </n-button>
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
        @click="handelPlayerAction('unban')"
        type="success"
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
        {{ $t("button.unban") }}
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
  </div>
  <!-- 帕鲁详情 modal -->
  <n-modal
    v-model:show="showPalDetailModal"
    preset="card"
    :style="{ width: '90%', maxWidth: '400px' }"
    header-style="padding:12px 20px;"
    content-style="padding:12px 20px;margin:0;"
    size="huge"
    :bordered="false"
    :segmented="{ content: 'soft', footer: 'soft' }"
  >
    <template #header-extra>
      <div class="flex pr-3 space-x-2">
        <n-tag type="primary" round> Lv.{{ palDetail.level }} </n-tag>
        <n-tag :type="palDetail.gender === 'Male' ? 'primary' : 'error'" round>
          {{ palDetail.gender === "Male" ? "♂" : "♀" }}
        </n-tag>
      </div>
    </template>
    <template #header>
      {{
        palDetail.nickname == ""
          ? getPalName(palDetail.type)
          : palDetail.nickname + "(" + getPalName(palDetail.type) + ")"
      }}
    </template>
    <pal-detail :palDetail="palDetail"></pal-detail>
  </n-modal>

  <!-- 添加白名单 modal -->
  <n-modal
    v-model:show="showAddWhiteListModal"
    class="custom-card"
    preset="card"
    style="width: 90%; max-width: 700px"
    footer-style="padding: 12px;"
    content-style="padding: 12px;"
    header-style="padding: 12px;"
    :title="$t('modal.addWhitelist')"
    :bordered="false"
  >
    <n-grid class="flex items-center">
      <n-gi span="5">
        <div class="flex justify-center">
          {{ $t("message.selectVerify") }}
        </div>
      </n-gi>
      <n-gi span="19">
        <n-input-group>
          <n-input
            v-model:value="addWhiteData.name"
            :style="{ width: '33%' }"
            :placeholder="$t('input.nickname')"
          />
          <n-input
            v-model:value="addWhiteData.player_uid"
            :style="{ width: '33%' }"
            :placeholder="$t('input.player_uid')"
          />
          <n-input
            v-model:value="addWhiteData.steam_id"
            :style="{ width: '33%' }"
            :placeholder="$t('input.steam_id')"
          />
        </n-input-group>
      </n-gi>
    </n-grid>
    <template #footer>
      <div class="flex justify-end">
        <n-button
          type="tertiary"
          @click="
            () => {
              showAddWhiteListModal = false;
            }
          "
        >
          {{ $t("button.cancel") }}
        </n-button>
        <n-button
          class="ml-3 w-40"
          type="primary"
          @click="addWhiteList"
          :disabled="
            !addWhiteData.name ||
            (!addWhiteData.player_uid && !addWhiteData.steam_id)
          "
        >
          {{ $t("button.confirm") }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<style scoped lang="less">
.player-detail {
  min-height: 100%;
}

.player-overview {
  padding: 22px;
  border: 1px solid rgba(64, 152, 252, 0.22);
  border-radius: 16px;
  background: linear-gradient(
    135deg,
    rgba(64, 152, 252, 0.1),
    rgba(64, 152, 252, 0.025) 58%,
    transparent
  );
}

.is-dark .player-overview {
  border-color: rgba(64, 152, 252, 0.28);
  background: linear-gradient(
    135deg,
    rgba(64, 152, 252, 0.16),
    rgba(64, 152, 252, 0.045) 58%,
    transparent
  );
}

.overview-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
}

.player-identity {
  min-width: 0;
}

.player-title-row,
.identity-tags,
.runtime-meta,
.identity-copy-list {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.player-title-row {
  gap: 12px;
}

.player-title {
  min-width: 0;
  overflow-wrap: anywhere;
  font-size: clamp(24px, 2.2vw, 30px);
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.02em;
}

.identity-tags {
  gap: 8px;
  margin-top: 10px;
}

.last-online-text {
  color: rgba(24, 24, 28, 0.52);
  font-size: 13px;
  font-variant-numeric: tabular-nums;
}

.is-dark .last-online-text {
  color: rgba(255, 255, 255, 0.5);
}

.identity-copy-list {
  gap: 8px;
  margin-top: 18px;
}

.identity-copy {
  max-width: min(100%, 420px);
}

.copy-label {
  flex: none;
  font-weight: 650;
}

.copy-value {
  min-width: 0;
  overflow: hidden;
  color: rgba(24, 24, 28, 0.58);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.is-dark .copy-value {
  color: rgba(255, 255, 255, 0.54);
}

.runtime-meta {
  gap: 8px;
  margin-top: 12px;
}

.runtime-meta span {
  padding: 4px 9px;
  border-radius: 999px;
  background: rgba(24, 24, 28, 0.06);
  color: rgba(24, 24, 28, 0.62);
  font-size: 12px;
}

.is-dark .runtime-meta span {
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.62);
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 10px;
  margin-top: 20px;
}

.status-card {
  min-width: 0;
  padding: 13px 14px;
  border: 1px solid rgba(24, 24, 28, 0.06);
  border-radius: 11px;
  background: rgba(255, 255, 255, 0.7);
}

.is-dark .status-card {
  border-color: rgba(255, 255, 255, 0.07);
  background: rgba(255, 255, 255, 0.055);
}

.status-label {
  display: block;
  overflow: hidden;
  color: rgba(24, 24, 28, 0.52);
  font-size: 12px;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.is-dark .status-label {
  color: rgba(255, 255, 255, 0.5);
}

.status-value {
  display: block;
  margin-top: 6px;
  font-size: 22px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
  line-height: 1;
}

.detail-tabs {
  margin-top: 20px;
}

.player-actions {
  position: sticky;
  bottom: 0;
  z-index: 5;
  padding: 12px 20px;
  border-top: 1px solid rgba(24, 24, 28, 0.08);
  background: rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(12px);
}

.is-dark .player-actions {
  border-top-color: rgba(255, 255, 255, 0.08);
  background: rgba(24, 24, 28, 0.94);
}

@media (max-width: 1100px) {
  .status-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .overview-heading {
    align-items: stretch;
    flex-direction: column;
  }

  .overview-heading > .n-button {
    align-self: flex-start;
  }
}
</style>
