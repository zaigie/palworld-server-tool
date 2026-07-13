<script setup>
import { ChevronForward } from "@vicons/ionicons5";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";

const props = defineProps({
  guildList: { type: Array, default: () => [] },
});
const emit = defineEmits(["onGetInfo"]);
const { t } = useI18n();

const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);
const searchValue = ref("");
const namedFilter = ref("all");
const sortBy = ref("level");

const namedOptions = computed(() => [
  { label: t("filter.allGuilds"), value: "all" },
  { label: t("filter.namedGuilds"), value: "named" },
  { label: t("filter.unnamedGuilds"), value: "unnamed" },
]);
const sortOptions = computed(() => [
  { label: t("filter.levelHighToLow"), value: "level" },
  { label: t("filter.memberCount"), value: "members" },
  { label: t("filter.baseCount"), value: "bases" },
  { label: t("filter.guildName"), value: "name" },
]);
const isUnnamedGuild = (guild) => !guild.name || guild.name === "无名公会";
const guildMaster = (guild) =>
  guild.players?.find((player) => player.player_uid === guild.admin_player_uid);
const guildDisplayName = (guild) => {
  if (!isUnnamedGuild(guild)) return guild.name;
  const suffix =
    guildMaster(guild)?.nickname || guild.admin_player_uid?.slice(-6) || "—";
  return `${t("filter.unnamedGuild")} · ${suffix}`;
};
const filteredGuilds = computed(() => {
  const keyword = searchValue.value.trim().toLowerCase();
  return (props.guildList || [])
    .filter((guild) => {
      const searchable = [
        guild.name,
        guild.admin_player_uid,
        ...(guild.players || []).map((player) => player.nickname),
      ]
        .filter(Boolean)
        .join(" ")
        .toLowerCase();
      if (keyword && !searchable.includes(keyword)) return false;
      if (namedFilter.value === "named" && isUnnamedGuild(guild)) return false;
      if (namedFilter.value === "unnamed" && !isUnnamedGuild(guild))
        return false;
      return true;
    })
    .sort((a, b) => {
      if (sortBy.value === "members")
        return (b.players?.length || 0) - (a.players?.length || 0);
      if (sortBy.value === "bases")
        return (b.base_camp?.length || 0) - (a.base_camp?.length || 0);
      if (sortBy.value === "name")
        return guildDisplayName(a).localeCompare(guildDisplayName(b));
      return Number(b.base_camp_level || 0) - Number(a.base_camp_level || 0);
    });
});
</script>

<template>
  <div class="guild-list" :class="{ 'is-dark': isDarkMode }">
    <div class="mobile-list-controls">
      <n-input
        v-model:value="searchValue"
        clearable
        size="large"
        :placeholder="$t('filter.searchGuilds')"
        :aria-label="$t('filter.searchGuilds')"
      />
      <n-grid :cols="2" :x-gap="8" class="mt-2">
        <n-gi>
          <n-select
            v-model:value="namedFilter"
            :options="namedOptions"
            :aria-label="$t('filter.guildType')"
          />
        </n-gi>
        <n-gi>
          <n-select
            v-model:value="sortBy"
            :options="sortOptions"
            :aria-label="$t('filter.guildSorting')"
          />
        </n-gi>
      </n-grid>
      <n-text depth="3" class="result-count">
        {{ $t("filter.resultCount", { count: filteredGuilds.length }) }}
      </n-text>
    </div>

    <n-list :show-divider="false" class="mobile-guild-list">
      <n-list-item
        v-for="guild in filteredGuilds"
        :key="guild.admin_player_uid"
        class="mobile-guild-item"
        @click="emit('onGetInfo', guild.admin_player_uid)"
        @keydown.enter="emit('onGetInfo', guild.admin_player_uid)"
        @keydown.space.prevent="emit('onGetInfo', guild.admin_player_uid)"
        role="button"
        tabindex="0"
        :aria-label="guildDisplayName(guild)"
      >
        <div class="mobile-guild-row">
          <div class="guild-row-main">
            <div class="guild-name-line">
              <strong class="guild-name">{{ guildDisplayName(guild) }}</strong>
              <n-tag :bordered="false" type="primary" size="small" round>
                Lv.{{ guild.base_camp_level || 0 }}
              </n-tag>
            </div>
            <div class="guild-summary">
              {{
                $t("filter.guildSummary", {
                  members: guild.players?.length || 0,
                  bases: guild.base_camp?.length || 0,
                })
              }}
            </div>
            <div class="guild-master">
              {{ $t("guild.master") }}
              <span>{{ guildMaster(guild)?.nickname || "—" }}</span>
            </div>
          </div>
          <n-icon class="row-chevron" size="18">
            <ChevronForward />
          </n-icon>
        </div>
      </n-list-item>
    </n-list>
    <n-empty v-if="filteredGuilds.length === 0" class="empty-state" />
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

.mobile-guild-list {
  padding: 0 10px 16px;
  background: transparent;
}

.mobile-guild-item {
  padding: 0 0 8px !important;
  outline: none;
}

.mobile-guild-item:focus-visible .mobile-guild-row {
  box-shadow: 0 0 0 2px rgba(64, 152, 252, 0.45);
}

.mobile-guild-row {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  min-width: 0;
  padding: 14px;
  border: 1px solid rgba(24, 24, 28, 0.07);
  border-radius: 12px;
  background: rgba(24, 24, 28, 0.025);
}

.is-dark .mobile-guild-row {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
}

.guild-row-main {
  flex: 1;
  min-width: 0;
}

.guild-name-line {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.guild-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  font-size: 17px;
  font-weight: 650;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.guild-summary {
  margin-top: 8px;
  color: rgba(24, 24, 28, 0.62);
  font-size: 12px;
}

.guild-master {
  margin-top: 5px;
  color: rgba(24, 24, 28, 0.42);
  font-size: 12px;
}

.guild-master span {
  margin-left: 5px;
  color: rgba(24, 24, 28, 0.65);
}

.is-dark .guild-summary,
.is-dark .guild-master span {
  color: rgba(255, 255, 255, 0.6);
}

.is-dark .guild-master {
  color: rgba(255, 255, 255, 0.36);
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
