<script setup>
import { GroupWorkRound } from "@vicons/material";
import { CrownFilled } from "@vicons/antd";
import { computed } from "vue";

const props = defineProps(["guildInfo"]);
const guildInfo = computed(() => props.guildInfo);

const getUserAvatar = () => {
  return new URL("@/assets/avatar.webp", import.meta.url).href;
};
</script>
<template>
  <div class="guile-detail">
    <n-layout :native-scrollbar="false">
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
                  v-if="player.player_uid === guildInfo.admin_player_uid"
                  type="error"
                  size="small"
                >
                  {{ $t("status.master") }}
                </n-tag>
              </n-space>
            </n-list-item>
          </n-list>
        </n-space>
      </n-card>
    </n-layout>
  </div>
</template>
