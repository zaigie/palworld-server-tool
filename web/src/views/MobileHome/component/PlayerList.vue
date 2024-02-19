<script setup>
import dayjs from "dayjs";
import { computed } from "vue";

const props = defineProps(["playerList"]);
const playerList = computed(() => props.playerList);

const emits = defineEmits(["onGetInfo"]);

const getUserAvatar = () => {
  return new URL("@/assets/avatar.webp", import.meta.url).href;
};

const onGetInfo = (uid) => {
  emits("onGetInfo", uid);
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
  <div class="player-list">
    <n-list hoverable clickable>
      <n-list-item
        v-for="player in playerList"
        :key="player.player_uid"
        @click="onGetInfo(player.player_uid)"
      >
        <template #prefix>
          <n-avatar :src="getUserAvatar()" round></n-avatar>
        </template>
        <div>
          <div class="flex">
            <n-tag
              :bordered="false"
              size="small"
              :type="isPlayerOnline(player.last_online) ? 'success' : 'error'"
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
              isDarkMode ? 'bg-#2f69aa text-#fff' : 'bg-#ddd text-#18181c'
            "
            class="inline-block mt-1 rounded-full text-xs px-2 py-0.5"
            >{{ $t("status.last_online") }}:
            {{ displayLastOnline(player.last_online) }}</span
          >
        </div>
      </n-list-item>
    </n-list>
  </div>
</template>
