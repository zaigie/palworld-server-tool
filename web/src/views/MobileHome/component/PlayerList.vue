<script setup>
import dayjs from "dayjs";
import { computed, ref } from "vue";
import { ChevronForward } from "@vicons/ionicons5";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);

const props = defineProps({
  playerList: { type: Array, default: () => [] },
});
const playerList = computed(() => props.playerList || []);

const emits = defineEmits(["onGetInfo"]);

const searchValue = ref("");
const statusFilter = ref("all");
const sortBy = ref("last_online");

const platformColors = {
  steam: { color: "#223D58", textColor: "#fff" },
  xbox: { color: "#2B8B2B", textColor: "#fff" },
  ps5: { color: "#00439C", textColor: "#fff" },
  mac: { color: "#777", textColor: "#fff" },
  default: { color: "#d9c36c", textColor: "#fff" },
};

const onGetInfo = (uid) => {
  emits("onGetInfo", uid);
};

const isPlayerOnline = (last_online) => {
  return dayjs() - dayjs(last_online) < 80000;
};
const displayLastOnline = (last_online) => {
  if (dayjs(last_online).year() < 1970) {
    return "Unknown";
  }
  return dayjs(last_online).format("YYYY-MM-DD HH:mm:ss");
};
const getPlatformColor = (userId) => {
  if (!userId) return platformColors.default;
  return platformColors[userId.split("_")[0]] || platformColors.default;
};

const statusOptions = computed(() => [
  { label: t("filter.allStatuses"), value: "all" },
  { label: t("status.online"), value: "online" },
  { label: t("status.offline"), value: "offline" },
]);
const sortOptions = computed(() => [
  { label: t("filter.lastOnline"), value: "last_online" },
  { label: t("filter.levelHighToLow"), value: "level" },
  { label: t("filter.nickname"), value: "nickname" },
]);
const filteredPlayers = computed(() => {
  const keyword = searchValue.value.trim().toLowerCase();
  return playerList.value
    .filter((player) => {
      const searchable = [
        player.nickname,
        player.player_uid,
        player.user_id,
        player.steam_id,
      ]
        .filter(Boolean)
        .join(" ")
        .toLowerCase();
      if (keyword && !searchable.includes(keyword)) return false;
      const online = isPlayerOnline(player.last_online);
      if (statusFilter.value === "online" && !online) return false;
      if (statusFilter.value === "offline" && online) return false;
      return true;
    })
    .sort((a, b) => {
      if (sortBy.value === "level") {
        return Number(b.level || 0) - Number(a.level || 0);
      }
      if (sortBy.value === "nickname") {
        return (a.nickname || "").localeCompare(b.nickname || "");
      }
      return dayjs(b.last_online).valueOf() - dayjs(a.last_online).valueOf();
    });
});
</script>

<template>
  <div class="player-list" :class="{ 'is-dark': isDarkMode }">
    <div class="mobile-list-controls">
      <n-input
        v-model:value="searchValue"
        clearable
        size="large"
        :placeholder="$t('filter.searchPlayers')"
        :aria-label="$t('filter.searchPlayers')"
      />
      <n-grid :cols="2" :x-gap="8" class="mt-2">
        <n-gi>
          <n-select v-model:value="statusFilter" :options="statusOptions" />
        </n-gi>
        <n-gi>
          <n-select v-model:value="sortBy" :options="sortOptions" />
        </n-gi>
      </n-grid>
      <n-text depth="3" class="result-count">
        {{ $t("filter.resultCount", { count: filteredPlayers.length }) }}
      </n-text>
    </div>
    <n-list :show-divider="false" class="mobile-player-list">
      <n-list-item
        v-for="player in filteredPlayers"
        :key="player.player_uid"
        class="mobile-player-item"
        @click="onGetInfo(player.player_uid)"
        @keydown.enter="onGetInfo(player.player_uid)"
        @keydown.space.prevent="onGetInfo(player.player_uid)"
        role="button"
        tabindex="0"
        :aria-label="`${player.nickname}, Lv.${player.level}`"
      >
        <div class="mobile-player-row">
          <div class="player-row-main">
            <div class="player-name-line">
              <span class="player-name">{{ player.nickname || "--" }}</span>
              <n-tag
                v-if="player.user_id"
                :bordered="false"
                round
                size="small"
                :color="getPlatformColor(player.user_id)"
              >
                {{ player.user_id.split("_")[0] }}
              </n-tag>
            </div>
            <div class="player-meta-line">
              <span class="player-status">
                <span
                  class="status-dot"
                  :class="
                    isPlayerOnline(player.last_online)
                      ? 'is-online'
                      : 'is-offline'
                  "
                ></span>
                {{
                  isPlayerOnline(player.last_online)
                    ? $t("status.online")
                    : $t("status.offline")
                }}
              </span>
              <n-tag :bordered="false" type="primary" size="small" round>
                Lv.{{ player.level }}
              </n-tag>
            </div>
            <div class="last-online">
              {{ $t("status.last_online") }}
              <span>{{ displayLastOnline(player.last_online) }}</span>
            </div>
          </div>
          <n-icon class="row-chevron" size="18">
            <ChevronForward />
          </n-icon>
        </div>
      </n-list-item>
    </n-list>
    <n-empty v-if="filteredPlayers.length === 0" class="empty-state" />
  </div>
</template>

<style scoped lang="less">
.mobile-list-controls {
  padding: 14px 14px 10px;
}

.result-count {
  display: block;
  margin-top: 10px;
  padding: 0 2px;
  font-size: 13px;
}

.mobile-player-list {
  padding: 0 10px 16px;
  background: transparent;
}

.mobile-player-item {
  padding: 0 0 8px !important;
  outline: none;
}

.mobile-player-item:focus-visible .mobile-player-row {
  box-shadow: 0 0 0 2px rgba(64, 152, 252, 0.45);
}

.mobile-player-row {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px;
  border: 1px solid rgba(24, 24, 28, 0.07);
  border-radius: 12px;
  background: rgba(24, 24, 28, 0.025);
}

.is-dark .mobile-player-row {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
}

.player-row-main {
  flex: 1;
  min-width: 0;
}

.player-name-line,
.player-meta-line,
.player-status {
  display: flex;
  align-items: center;
}

.player-name-line {
  min-width: 0;
  gap: 8px;
}

.player-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  font-size: 17px;
  font-weight: 650;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-meta-line {
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.player-status {
  gap: 6px;
  color: rgba(24, 24, 28, 0.6);
  font-size: 12px;
}

.is-dark .player-status {
  color: rgba(255, 255, 255, 0.6);
}

.status-dot {
  width: 8px;
  height: 8px;
  flex: none;
  border-radius: 50%;

  &.is-online {
    background: #18a058;
    box-shadow: 0 0 0 3px rgba(24, 160, 88, 0.14);
  }

  &.is-offline {
    background: #d03050;
  }
}

.last-online {
  margin-top: 8px;
  color: rgba(24, 24, 28, 0.42);
  font-size: 12px;
}

.last-online span {
  margin-left: 5px;
  color: rgba(24, 24, 28, 0.65);
  font-variant-numeric: tabular-nums;
}

.is-dark .last-online {
  color: rgba(255, 255, 255, 0.36);
}

.is-dark .last-online span {
  color: rgba(255, 255, 255, 0.6);
}

.row-chevron {
  flex: none;
  color: rgba(24, 24, 28, 0.28);
}

.is-dark .row-chevron {
  color: rgba(255, 255, 255, 0.3);
}

.empty-state {
  padding: 50px 16px;
}
</style>
