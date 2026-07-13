<script setup>
import { ChevronForward } from "@vicons/ionicons5";
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import GuildOverview from "@/components/GuildOverview.vue";
import pageStore from "@/stores/model/page.js";
import playerToGuildStore from "@/stores/model/playerToGuild";

const props = defineProps({ guilds: { type: Array, default: () => [] } });
const { t } = useI18n();

const pageWidth = computed(() => pageStore().getScreenWidth());
const smallScreen = computed(() => pageWidth.value < 1024);
const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);

const loadingGuild = ref(false);
const loadingGuildDetail = ref(false);
const guildList = ref([]);
const guildInfo = ref({});
const searchValue = ref("");
const namedFilter = ref("all");
const sortBy = ref("level");

const getGuildList = async () => {
  if (props.guilds.length > 0) {
    guildList.value = [...props.guilds];
    return;
  }
  const { data } = await new ApiService().getGuildList();
  guildList.value = Array.isArray(data.value) ? data.value : [];
};

const getGuildInfo = async (adminPlayerUid) => {
  const { data } = await new ApiService().getGuild({ adminPlayerUid });
  guildInfo.value = data.value || {};
};

const clickGetGuildInfo = async (uid) => {
  if (guildInfo.value.admin_player_uid === uid) return;
  loadingGuildDetail.value = true;
  await getGuildInfo(uid);
  loadingGuildDetail.value = false;
};

const toPlayers = (uid) => {
  playerToGuildStore().setCurrentUid(uid);
  playerToGuildStore().setUpdateStatus("players");
};

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
  return guildList.value
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

watch(
  () => props.guilds,
  (guilds) => {
    if (guilds?.length > 0) guildList.value = [...guilds];
  },
  { deep: true },
);

onMounted(async () => {
  loadingGuild.value = true;
  loadingGuildDetail.value = true;
  await getGuildList();
  loadingGuild.value = false;
  if (guildList.value.length > 0) {
    const currentUid = playerToGuildStore().getCurrentUid();
    const initialGuild =
      guildList.value.find((guild) => guild.admin_player_uid === currentUid) ||
      guildList.value[0];
    await getGuildInfo(initialGuild.admin_player_uid);
    playerToGuildStore().setCurrentUid(null);
  }
  loadingGuildDetail.value = false;
});
</script>

<template>
  <div class="guild-list h-full" :class="{ 'is-dark': isDarkMode }">
    <n-layout has-sider class="h-full">
      <n-layout-sider
        :width="smallScreen ? 360 : 400"
        content-style="padding: 0 16px 16px;"
        :native-scrollbar="false"
        bordered
        class="guild-list-sidebar relative"
      >
        <div class="filter-panel">
          <n-input
            v-model:value="searchValue"
            clearable
            size="large"
            :placeholder="$t('filter.searchGuilds')"
            :aria-label="$t('filter.searchGuilds')"
          />
          <n-grid cols="2" :x-gap="8" class="mt-2">
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

        <n-list class="guild-list-items" :show-divider="false">
          <n-list-item
            v-for="guild in filteredGuilds"
            :key="guild.admin_player_uid"
            class="guild-row-item"
            @click="clickGetGuildInfo(guild.admin_player_uid)"
            @keydown.enter="clickGetGuildInfo(guild.admin_player_uid)"
            @keydown.space.prevent="clickGetGuildInfo(guild.admin_player_uid)"
            role="button"
            tabindex="0"
            :aria-label="guildDisplayName(guild)"
            :aria-current="
              guildInfo.admin_player_uid === guild.admin_player_uid
                ? 'true'
                : undefined
            "
          >
            <div
              class="guild-row"
              :class="{
                'is-selected':
                  guildInfo.admin_player_uid === guild.admin_player_uid,
              }"
            >
              <div class="guild-row-main">
                <div class="guild-name-line">
                  <strong class="guild-name" :title="guildDisplayName(guild)">
                    {{ guildDisplayName(guild) }}
                  </strong>
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

        <n-empty
          v-if="!loadingGuild && filteredGuilds.length === 0"
          class="empty-state"
        />
        <n-spin
          v-if="loadingGuild"
          size="small"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description>{{ $t("message.loading") }}</template>
        </n-spin>
      </n-layout-sider>

      <n-layout class="relative" :native-scrollbar="false">
        <guild-overview :guild-info="guildInfo" @view-player="toPlayers" />
        <n-spin
          v-if="loadingGuildDetail"
          size="small"
          class="absolute top-0 left-0 w-full h-full flex items-center justify-center bg-#ffffff40"
        >
          <template #description>{{ $t("message.loading") }}</template>
        </n-spin>
      </n-layout>
    </n-layout>
  </div>
</template>

<style scoped lang="less">
.filter-panel {
  position: sticky;
  top: 0;
  z-index: 3;
  padding: 20px 0 14px;
  background: #fff;
}

.is-dark .filter-panel {
  background: #18181c;
}

.result-count {
  display: block;
  margin-top: 10px;
  padding: 0 2px;
  font-size: 13px;
}

.guild-list-items {
  background: transparent;
}

.guild-row-item {
  padding: 0 0 4px !important;
  outline: none;
}

.guild-row-item:focus-visible .guild-row {
  box-shadow: 0 0 0 2px rgba(64, 152, 252, 0.45);
}

.guild-row {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  min-width: 0;
  padding: 11px 12px;
  border: 1px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition:
    background-color 0.18s ease,
    border-color 0.18s ease,
    box-shadow 0.18s ease;
}

.guild-row:hover {
  background: rgba(64, 152, 252, 0.08);
}

.guild-row.is-selected {
  border-color: rgba(64, 152, 252, 0.55);
  background: rgba(64, 152, 252, 0.12);
  box-shadow: inset 3px 0 0 #4098fc;
}

.is-dark .guild-row:hover {
  background: rgba(64, 152, 252, 0.12);
}

.is-dark .guild-row.is-selected {
  background: rgba(64, 152, 252, 0.17);
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
  margin-top: 6px;
  color: rgba(24, 24, 28, 0.65);
  font-size: 12px;
}

.guild-master {
  margin-top: 4px;
  color: rgba(24, 24, 28, 0.42);
  font-size: 12px;
}

.guild-master span {
  margin-left: 5px;
  color: rgba(24, 24, 28, 0.65);
}

.is-dark .guild-summary,
.is-dark .guild-master span {
  color: rgba(255, 255, 255, 0.62);
}

.is-dark .guild-master {
  color: rgba(255, 255, 255, 0.38);
}

.row-chevron {
  flex: none;
  color: rgba(24, 24, 28, 0.28);
  transition: transform 0.18s ease;
}

.guild-row:hover .row-chevron,
.guild-row.is-selected .row-chevron {
  color: #4098fc;
  transform: translateX(2px);
}

.is-dark .row-chevron {
  color: rgba(255, 255, 255, 0.3);
}

.empty-state {
  padding: 48px 12px;
}
</style>
