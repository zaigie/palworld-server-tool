<script setup>
import { ContentCopyFilled } from "@vicons/material";
import { LogOut, Ban, Search } from "@vicons/ionicons5";
import { CrownFilled } from "@vicons/antd";
import dayjs from "dayjs";
import { onMounted, computed } from "vue";
import { NTag, NButton, NAvatar, useMessage, useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";
import palMap from "@/assets/pal.json";
import skillMap from "@/assets/skill.json";
import PalDetail from "./PalDetail.vue";
import userStore from "@/stores/model/user";
import ApiService from "@/service/api.js";
import {
  localizedSkillName,
  statusPointTranslationKey,
} from "@/utils/gameLabels";

const { t, locale } = useI18n();

const message = useMessage();
const dialog = useDialog();

const localeLowerPalMap = ref({});
const getStatusPointLabel = (rawKey) => {
  const translationKey = statusPointTranslationKey(rawKey);
  return translationKey ? t(`statusPoint.${translationKey}`) : rawKey;
};
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

const isLogin = computed(() => userStore().getLoginInfo().isLogin);

const props = defineProps(["playerInfo", "currentPlayerPalsList", "finished"]);
const playerInfo = computed(() => props.playerInfo);
const currentPlayerPalsList = computed(() => props.currentPlayerPalsList);
const finished = computed(() => props.finished);

const emits = defineEmits(["onSearch"]);

const handelPlayerAction = async (type) => {
  if (!isLogin.value) {
    message.error(t("message.requireauth"));
    return;
  } else {
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
            playerUid: playerInfo.value.player_uid,
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
  }
};

const searchValue = ref("");
const clickSearch = (input) => {
  emits("onSearch", input);
};

// 查看帕鲁详情
const showPalDetailModal = ref(false);
const palDetail = ref({});

const showPalDetail = (pal) => {
  palDetail.value = pal;
  showPalDetailModal.value = true;
};

const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 80000;
};

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

const displayHP = (hp, max_hp) => {
  return (hp / 1000).toFixed(0) + "/" + (max_hp / 1000).toFixed(0);
};

const percentageHP = (hp, max_hp) => {
  if (max_hp === 0) {
    return 0;
  }
  return ((hp / max_hp) * 100).toFixed(2);
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
const getUnknowPalAvatar = (is_boss = false) => {
  if (is_boss) {
    return new URL("@/assets/pals/boss_unknown.png", import.meta.url).href;
  }
  return new URL("@/assets/pals/unknown.png", import.meta.url).href;
};
const getPlatformColor = (userId) => {
  if (!userId) return platformColors.default;
  return platformColors[userId.split("_")[0]] || platformColors.default;
};
const displayLastOnline = (lastOnline) => {
  if (dayjs(lastOnline).year() < 1970) return "Unknown";
  return dayjs(lastOnline).format("YYYY-MM-DD HH:mm:ss");
};

onMounted(async () => {
  locale.value = localStorage.getItem("locale");
  localeLowerPalMap.value = Object.keys(palMap[locale.value]).reduce(
    (acc, key) => {
      acc[key.toLowerCase()] = palMap[locale.value][key];
      return acc;
    },
    {},
  );
});
</script>

<template>
  <div class="player-detail" :class="{ 'is-dark': isDarkMode }">
    <n-layout :native-scrollbar="false">
      <n-card
        :bordered="false"
        v-if="playerInfo.nickname"
        content-style="padding: 12px 12px 24px"
      >
        <section
          class="mobile-player-overview"
          aria-labelledby="mobile-player-title"
        >
          <div class="mobile-title-row">
            <h1 id="mobile-player-title" class="mobile-player-title">
              {{ playerInfo.nickname }}
            </h1>
            <n-tag type="primary" round strong>
              Lv.{{ playerInfo.level }}
              <template #icon>
                <n-icon :component="CrownFilled" />
              </template>
            </n-tag>
          </div>
          <div class="mobile-identity-tags">
            <n-tag
              :bordered="false"
              :type="
                isPlayerOnline(playerInfo.last_online) ? 'success' : 'error'
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
            <n-tag
              v-if="playerInfo.user_id"
              :bordered="false"
              round
              size="small"
              :color="getPlatformColor(playerInfo.user_id)"
            >
              {{ playerInfo.user_id.split("_")[0] }}
            </n-tag>
            <span class="mobile-last-online">
              {{ $t("status.last_online") }}
              {{ displayLastOnline(playerInfo.last_online) }}
            </span>
          </div>

          <div v-if="isLogin" class="mobile-copy-list">
            <n-button
              block
              secondary
              class="mobile-copy-button"
              @click="copyText(playerInfo.player_uid)"
            >
              <span class="copy-label">UID</span>
              <span class="copy-value">{{ playerInfo.player_uid }}</span>
              <template #icon>
                <n-icon><ContentCopyFilled /></n-icon>
              </template>
            </n-button>
            <n-button
              block
              secondary
              class="mobile-copy-button"
              @click="copyText(playerInfo.steam_id)"
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

          <div class="mobile-status-grid">
            <div
              v-for="status in Object.entries(playerInfo.status_point || {})"
              :key="status[0]"
              class="mobile-status-card"
            >
              <span>{{ getStatusPointLabel(status[0]) }}</span>
              <strong>{{ status[1] }}</strong>
            </div>
          </div>
        </section>

        <n-button-group v-if="isLogin" class="mobile-admin-actions">
          <n-button
            @click="handelPlayerAction('unban')"
            type="success"
            secondary
            strong
          >
            <template #icon>
              <n-icon><Ban /></n-icon>
            </template>
            {{ $t("button.unban") }}
          </n-button>
          <n-button
            @click="handelPlayerAction('ban')"
            type="error"
            secondary
            strong
          >
            <template #icon>
              <n-icon><Ban /></n-icon>
            </template>
            {{ $t("button.ban") }}
          </n-button>
          <n-button
            @click="handelPlayerAction('kick')"
            type="warning"
            secondary
            strong
          >
            <template #icon>
              <n-icon><LogOut /></n-icon>
            </template>
            {{ $t("button.kick") }}
          </n-button>
        </n-button-group>
        <!-- <n-space vertical>
          <n-progress
            type="line"
            status="error"
            indicator-placement="inside"
            :percentage="percentageHP(playerInfo.hp, playerInfo.max_hp)"
            :height="24"
            :border-radius="4"
            :fill-border-radius="0"
            >HP: {{ displayHP(playerInfo.hp, playerInfo.max_hp) }}</n-progress
          >
          <n-progress
            type="line"
            indicator-placement="inside"
            :percentage="
              percentageHP(playerInfo.shield_hp, playerInfo.shield_max_hp)
            "
            :height="24"
            :border-radius="4"
            :fill-border-radius="0"
            >SHIELD:
            {{
              displayHP(playerInfo.shield_hp, playerInfo.shield_max_hp)
            }}</n-progress
          >
        </n-space> -->
        <div class="pal-search">
          <n-input
            v-model:value="searchValue"
            :placeholder="$t('input.searchPlaceholder')"
            @update:value="clickSearch"
          >
            <template #suffix>
              <n-icon>
                <Search />
              </n-icon>
            </template>
          </n-input>
        </div>
        <n-list>
          <n-list-item
            v-for="(pal, index) in currentPlayerPalsList"
            :key="pal"
            class="py-2"
            @click="showPalDetail(pal)"
          >
            <div class="flex justify-between items-center">
              <n-avatar
                class="bg-#c5c5c5 rounded-md"
                :size="32"
                :src="getPalAvatar(pal.type)"
                :fallback-src="getUnknowPalAvatar(pal.is_boss)"
              ></n-avatar>
              <div class="flex-1 flex items-center justify-between ml-3">
                <n-tag
                  size="small"
                  :type="pal.gender == 'Male' ? 'primary' : 'warning'"
                  >{{ pal.gender == "Male" ? "♂" : "♀" }}</n-tag
                >
                <span class="px-3 flex-1 line-clamp-1">{{
                  getPalName(pal.type)
                }}</span>
                <span>{{ "Lv." + pal.level }}</span>
              </div>
            </div>
            <div class="ml-11 mt-1 flex flex-wrap">
              <n-tag
                v-for="skill in pal.skills"
                class="rounded-sm mr-2"
                size="small"
                :key="skill"
                :color="{
                  color: isDarkMode ? 'rgba(238, 155, 47, 0.15)' : '#fcf0e0',
                  textColor: '#ee9b2f',
                  borderColor: 'transparent',
                }"
                >{{ localizedSkillName(skill, locale, skillMap) }}</n-tag
              >
            </div>
          </n-list-item>
        </n-list>
        <div v-if="finished" class="text-center pt-4 color-#999">
          没有更多了
        </div>
        <div class="h-10"></div>
      </n-card>
    </n-layout>
  </div>
  <!-- 帕鲁详情 modal -->
  <n-modal
    v-model:show="showPalDetailModal"
    preset="card"
    :style="{ width: '95%', maxWidth: '400px' }"
    header-style="padding:12px;"
    content-style="margin:0;padding:12px;"
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
        palDetail.nickname == ""
          ? getPalName(palDetail.type)
          : palDetail.nickname + "(" + getPalName(palDetail.type) + ")"
      }}
    </template>
    <pal-detail :palDetail="palDetail"></pal-detail>
  </n-modal>
</template>

<style scoped lang="less">
.mobile-player-overview {
  padding: 16px;
  border: 1px solid rgba(64, 152, 252, 0.22);
  border-radius: 14px;
  background: linear-gradient(
    145deg,
    rgba(64, 152, 252, 0.11),
    rgba(64, 152, 252, 0.025) 62%,
    transparent
  );
}

.is-dark .mobile-player-overview {
  border-color: rgba(64, 152, 252, 0.28);
  background: linear-gradient(
    145deg,
    rgba(64, 152, 252, 0.17),
    rgba(64, 152, 252, 0.045) 62%,
    transparent
  );
}

.mobile-title-row,
.mobile-identity-tags {
  display: flex;
  align-items: center;
}

.mobile-title-row {
  gap: 10px;
}

.mobile-player-title {
  flex: 1;
  min-width: 0;
  overflow-wrap: anywhere;
  font-size: 24px;
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.02em;
}

.mobile-identity-tags {
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.mobile-last-online {
  color: rgba(24, 24, 28, 0.5);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.is-dark .mobile-last-online {
  color: rgba(255, 255, 255, 0.48);
}

.mobile-copy-list {
  display: grid;
  gap: 8px;
  margin-top: 14px;
}

.mobile-copy-button {
  min-width: 0;
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
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.is-dark .copy-value {
  color: rgba(255, 255, 255, 0.54);
}

.mobile-status-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin-top: 16px;
}

.mobile-status-card {
  min-width: 0;
  padding: 10px;
  border: 1px solid rgba(24, 24, 28, 0.06);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.72);
}

.is-dark .mobile-status-card {
  border-color: rgba(255, 255, 255, 0.07);
  background: rgba(255, 255, 255, 0.055);
}

.mobile-status-card span {
  display: block;
  overflow: hidden;
  color: rgba(24, 24, 28, 0.5);
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.is-dark .mobile-status-card span {
  color: rgba(255, 255, 255, 0.48);
}

.mobile-status-card strong {
  display: block;
  margin-top: 6px;
  font-size: 19px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
  line-height: 1;
}

.mobile-admin-actions {
  width: 100%;
  display: flex;
  margin-top: 10px;
}

.mobile-admin-actions > .n-button {
  flex: 1;
}

.pal-search {
  width: 100%;
  margin-top: 18px;
}
</style>
