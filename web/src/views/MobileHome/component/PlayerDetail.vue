<script setup>
import { ContentCopyFilled } from "@vicons/material";
import { LogOut, Ban } from "@vicons/ionicons5";
import { CrownFilled } from "@vicons/antd";
import dayjs from "dayjs";
import { computed } from "vue";
import { NTag, NButton, NAvatar, useMessage, useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";
import palMap from "@/assets/pal.json";
import skillMap from "@/assets/skill.json";
import PalDetail from "./PalDetail.vue";
import userStore from "@/stores/model/user";

const { t, locale } = useI18n();

const message = useMessage();
const dialog = useDialog();

const isLogin = computed(() => userStore().getLoginInfo().isLogin);

const props = defineProps(["playerInfo", "currentPlayerPalsList", "finished"]);
const playerInfo = computed(() => props.playerInfo);
const currentPlayerPalsList = computed(() => props.currentPlayerPalsList);
const finished = computed(() => props.finished);

const emits = defineEmits(["onSearch"]);

const handelPlayerAction = async (type) => {
  if (!isLogin.value) {
    message.error($t("message.requireauth"));
    showLoginModal.value = true;
    return;
  } else {
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
  }
};

const searchValue = ref("");
const clickSearch = () => {
  emits("onSearch", searchValue.value);
};

// 查看帕鲁详情
const showPalDetailModal = ref(false);
const palDetail = ref({});

const showPalDetail = (pal) => {
  palDetail.value = pal;
  showPalDetailModal.value = true;
};

const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 120000;
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

const getUserAvatar = () => {
  return new URL("@/assets/avatar.webp", import.meta.url).href;
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
  return new URL(`../../../assets/pal/${name}.png`, import.meta.url).href;
};
const getUnknowPalAvatar = () => {
  return new URL("@/assets/pal/Unknown.png", import.meta.url).href;
};
</script>

<template>
  <div class="player-detail">
    <n-layout :native-scrollbar="false">
      <!-- ban / kick -->
      <div v-if="isLogin" class="pt-2 px-3 bg-transparent" position="absolute">
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
      <n-card
        :bordered="false"
        v-if="playerInfo.nickname"
        content-style="padding: 12px"
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
            <div class="flex items-center w-full">
              <span class="flex-1 text-sm line-clamp-1 pr-1">
                {{ playerInfo.nickname }}
              </span>
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
            </div>
            <n-tag
              @click="copyText(playerInfo.player_uid)"
              class="mt-1 mr-2"
              type="info"
              size="small"
              icon-placement="right"
              v-if="isLogin"
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
              v-if="isLogin"
              ghost
            >
              Steam64:
              {{
                playerInfo.steam_id && playerInfo.steam_id.length === 17
                  ? playerInfo.steam_id
                  : "--"
              }}
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
        </n-space>
        <div class="flex w-full mt-5 border-b border-b-solid border-b-#eee">
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
              <n-avatar
                class="bg-#c5c5c5 rounded-md"
                :size="32"
                :src="getPalAvatar(pal.type)"
                :fallback-src="getUnknowPalAvatar()"
              ></n-avatar>
              <div class="flex-1 flex items-center justify-between ml-3">
                <van-tag
                  plain
                  :type="pal.gender == 'Male' ? 'primary' : 'danger'"
                  >{{ pal.gender == "Male" ? "♂" : "♀" }}</van-tag
                >
                <span class="px-3 flex-1 line-clamp-1">{{
                  palMap[locale][pal.type] ? palMap[locale][pal.type] : pal.type
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
                >{{
                  skillMap[locale][skill] ? skillMap[locale][skill].name : skill
                }}</van-tag
              >
            </div>
          </div>
        </van-list>
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
        palMap[locale][palDetail.type]
          ? palMap[locale][palDetail.type]
          : palDetail.type
      }}
    </template>
    <pal-detail :palDetail="palDetail"></pal-detail>
  </n-modal>
</template>
