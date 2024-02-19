<script setup>
import { GroupWorkRound, PersonSearchSharp } from "@vicons/material";
import { CrownFilled } from "@vicons/antd";
import ApiService from "@/service/api";
import pageStore from "@/stores/model/page.js";
import { ref, onMounted } from "vue";
import playerToGuildStore from "@/stores/model/playerToGuild";
import whitelistStore from "@/stores/model/whitelist";

const pageWidth = computed(() => pageStore().getScreenWidth());
const smallScreen = computed(() => pageWidth.value < 1024);

const loadingGuild = ref(false);
const loadingGuildDetail = ref(false);
const guildList = ref([]);
const guildInfo = ref({});

// 获取公会列表
const getGuildList = async () => {
  const { data } = await new ApiService().getGuildList();
  guildList.value = data.value;
};

// 获取公会详情信息
const getGuildInfo = async (admin_player_uid) => {
  const { data } = await new ApiService().getGuild({
    adminPlayerUid: admin_player_uid,
  });
  guildInfo.value = data.value;
};

const clickGetGuildInfo = async (id) => {
  loadingGuildDetail.value = true;
  await getGuildInfo(id);
  loadingGuildDetail.value = false;
};

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

// 查看玩家
const ToPlayers = async (uid) => {
  playerToGuildStore().setCurrentUid(uid);
  playerToGuildStore().setUpdateStatus("players");
};

onMounted(async () => {
  loadingGuild.value = true;
  loadingGuildDetail.value = true;
  await getGuildList();
  loadingGuild.value = false;
  if (guildList.value.length > 0) {
    const currentUid = playerToGuildStore().getCurrentUid();
    await getGuildInfo(
      currentUid ? currentUid : guildList.value[0].admin_player_uid
    );
    playerToGuildStore().setCurrentUid(null);
  }
  loadingGuildDetail.value = false;
});

// 其他操作
const getUserAvatar = () => {
  return new URL("@/assets/avatar.webp", import.meta.url).href;
};
</script>

<template>
  <div class="guild-list h-full">
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
            v-for="guild in guildList"
            :key="guild.admin_player_uid"
            @click="clickGetGuildInfo(guild.admin_player_uid)"
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
        <n-spin
          size="small"
          v-if="loadingGuild"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description> 加载中... </template>
        </n-spin>
      </n-layout-sider>
      <n-layout class="relative" :native-scrollbar="false">
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
                <n-space size="large" class="flex items-center flex-wrap">
                  <n-avatar :src="getUserAvatar()" round></n-avatar>
                  {{ player.nickname }}
                  <n-tag :bordered="false" type="info" size="small">
                    UID: {{ player.player_uid }}
                  </n-tag>
                  <n-tag
                    v-if="player.player_uid === guildInfo.admin_player_uid"
                    type="error"
                    size="small"
                  >
                    {{ $t("status.master") }}
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
                  <n-button
                    @click="ToPlayers(player.player_uid)"
                    size="small"
                    type="warning"
                    icon-placement="right"
                    ghost
                  >
                    {{ $t("button.viewPlayer") }}
                    <template #icon>
                      <n-icon><PersonSearchSharp /></n-icon>
                    </template>
                  </n-button>
                </n-space> </n-list-item
            ></n-list>
          </n-space>
        </n-card>
        <n-spin
          size="small"
          v-if="loadingGuildDetail"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description> 加载中... </template>
        </n-spin>
      </n-layout>
    </n-layout>
  </div>
</template>
