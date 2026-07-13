<script setup>
import { PersonSearchSharp } from "@vicons/material";
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import whitelistStore from "@/stores/model/whitelist";

const LANDSCAPE = {
  maxX: 349400,
  maxY: 724400,
  minX: -1099400,
  minY: -724400,
};

const props = defineProps({
  guildInfo: { type: Object, default: () => ({}) },
  compact: Boolean,
});

const emit = defineEmits(["view-player"]);
const { t } = useI18n();
const quickMapTile = "map/tiles/0/0/0.png";
const activeBaseIndex = ref(null);
const isDarkMode = ref(
  window.matchMedia("(prefers-color-scheme: dark)").matches,
);

const guildInfo = computed(() => props.guildInfo || {});
const members = computed(() => guildInfo.value.players || []);
const bases = computed(() => guildInfo.value.base_camp || []);
const guildMaster = computed(
  () =>
    members.value.find(
      (player) => player.player_uid === guildInfo.value.admin_player_uid,
    ) || null,
);
const isUnnamedGuild = computed(
  () => !guildInfo.value.name || guildInfo.value.name === "无名公会",
);
const guildDisplayName = computed(() => {
  if (!isUnnamedGuild.value) return guildInfo.value.name;
  const suffix =
    guildMaster.value?.nickname ||
    guildInfo.value.admin_player_uid?.slice(-6) ||
    "—";
  return `${t("filter.unnamedGuild")} · ${suffix}`;
});

const whiteList = computed(() => whitelistStore().getWhitelist());
const isWhite = (player) => {
  if (!player) return false;
  return whiteList.value.some(
    (item) =>
      (item.player_uid && item.player_uid === player.player_uid) ||
      (item.steam_id && item.steam_id === player.steam_id),
  );
};

const clamp = (value, min, max) => Math.min(Math.max(value, min), max);
const baseMarkers = computed(() =>
  bases.value.map((base, index) => ({
    ...base,
    index,
    left: clamp(
      ((Number(base.location_y) - LANDSCAPE.minY) /
        (LANDSCAPE.maxY - LANDSCAPE.minY)) *
        100,
      2,
      98,
    ),
    top: clamp(
      ((LANDSCAPE.maxX - Number(base.location_x)) /
        (LANDSCAPE.maxX - LANDSCAPE.minX)) *
        100,
      2,
      98,
    ),
  })),
);

const formatCoordinate = (value) => {
  const number = Number(value);
  return Number.isFinite(number) ? Math.round(number).toLocaleString() : "—";
};
const formatRange = (value) => {
  const number = Number(value);
  return Number.isFinite(number) ? Math.round(number).toLocaleString() : "—";
};

const focusBase = (index) => {
  activeBaseIndex.value = index;
  document
    .getElementById(`guild-base-${index}`)
    ?.scrollIntoView({ behavior: "smooth", block: "nearest" });
};

watch(
  () => guildInfo.value.admin_player_uid,
  () => {
    activeBaseIndex.value = null;
  },
);
</script>

<template>
  <article
    v-if="guildInfo.admin_player_uid"
    class="guild-overview"
    :class="{ 'is-compact': compact, 'is-dark': isDarkMode }"
  >
    <header class="guild-hero">
      <div class="guild-heading">
        <span class="guild-eyebrow">{{ $t("guild.profile") }}</span>
        <h1 class="guild-title">{{ guildDisplayName }}</h1>
        <p class="guild-master-line">
          <span>{{ $t("guild.master") }}</span>
          <strong>{{ guildMaster?.nickname || "—" }}</strong>
          <n-tag v-if="isUnnamedGuild" :bordered="false" round size="small">
            {{ $t("filter.unnamedGuild") }}
          </n-tag>
        </p>
      </div>
      <div class="guild-level-badge" :aria-label="$t('guild.level')">
        <span>{{ $t("guild.level") }}</span>
        <strong>Lv.{{ guildInfo.base_camp_level || 0 }}</strong>
      </div>
    </header>

    <div class="guild-stats" :aria-label="$t('guild.summary')">
      <div class="guild-stat-card">
        <span>{{ $t("guild.members") }}</span>
        <strong>{{ members.length }}</strong>
      </div>
      <div class="guild-stat-card">
        <span>{{ $t("guild.bases") }}</span>
        <strong>{{ bases.length }}</strong>
      </div>
      <div class="guild-stat-card guild-stat-master">
        <span>{{ $t("guild.master") }}</span>
        <strong>{{ guildMaster?.nickname || "—" }}</strong>
      </div>
    </div>

    <div class="guild-content-grid">
      <section class="guild-section member-section">
        <div class="section-heading">
          <div>
            <span class="section-kicker">{{ $t("guild.roster") }}</span>
            <h2>{{ $t("guild.memberTitle", { count: members.length }) }}</h2>
          </div>
          <n-tag :bordered="false" round>{{ members.length }}</n-tag>
        </div>

        <div v-if="members.length" class="member-grid">
          <div
            v-for="player in members"
            :key="player.player_uid"
            class="member-card"
          >
            <div class="member-card-main">
              <div class="member-name-line">
                <strong class="member-name">{{
                  player.nickname || player.player_uid
                }}</strong>
                <n-tag
                  v-if="player.player_uid === guildInfo.admin_player_uid"
                  type="error"
                  size="small"
                  :bordered="false"
                  round
                >
                  {{ $t("status.master") }}
                </n-tag>
                <n-tag
                  v-if="isWhite(player)"
                  type="warning"
                  size="small"
                  :bordered="false"
                  round
                >
                  {{ $t("status.whitelist") }}
                </n-tag>
              </div>
              <span class="member-uid">UID {{ player.player_uid }}</span>
            </div>
            <n-button
              size="small"
              secondary
              type="warning"
              class="view-player-button"
              @click="emit('view-player', player.player_uid)"
            >
              <template #icon>
                <n-icon><PersonSearchSharp /></n-icon>
              </template>
              {{ $t("button.viewPlayer") }}
            </n-button>
          </div>
        </div>
        <n-empty v-else size="small" :description="$t('guild.noMembers')" />
      </section>

      <section class="guild-section base-section">
        <div class="section-heading">
          <div>
            <span class="section-kicker">{{ $t("guild.territory") }}</span>
            <h2>{{ $t("guild.baseTitle", { count: bases.length }) }}</h2>
          </div>
          <n-tag :bordered="false" round>{{ bases.length }}</n-tag>
        </div>

        <div v-if="bases.length" class="quick-map-wrap">
          <div class="quick-map-heading">
            <div>
              <strong>{{ $t("guild.quickMap") }}</strong>
              <span>{{
                $t("guild.quickMapHint", { count: bases.length })
              }}</span>
            </div>
            <span class="map-axis">X / Y</span>
          </div>
          <div class="quick-map" role="img" :aria-label="$t('guild.quickMap')">
            <img
              class="quick-map-background"
              :src="quickMapTile"
              alt=""
              aria-hidden="true"
            />
            <div class="quick-map-grid" aria-hidden="true"></div>
            <button
              v-for="marker in baseMarkers"
              :key="marker.id || marker.index"
              type="button"
              class="base-marker"
              :class="{ 'is-active': activeBaseIndex === marker.index }"
              :style="{ left: `${marker.left}%`, top: `${marker.top}%` }"
              :aria-label="
                $t('guild.baseMarker', {
                  index: marker.index + 1,
                  x: formatCoordinate(marker.location_x),
                  y: formatCoordinate(marker.location_y),
                })
              "
              @click="focusBase(marker.index)"
            >
              {{ marker.index + 1 }}
            </button>
          </div>
        </div>

        <div v-if="bases.length" class="base-list">
          <article
            v-for="(base, index) in bases"
            :id="`guild-base-${index}`"
            :key="base.id || index"
            class="base-card"
            :class="{ 'is-active': activeBaseIndex === index }"
            @mouseenter="activeBaseIndex = index"
          >
            <div class="base-card-index">{{ index + 1 }}</div>
            <div class="base-card-content">
              <div class="base-card-title">
                <strong>{{
                  $t("guild.baseLabel", { index: index + 1 })
                }}</strong>
                <span
                  >{{ $t("guild.range") }} {{ formatRange(base.area) }}</span
                >
              </div>
              <div class="base-coordinate">
                <span>X {{ formatCoordinate(base.location_x) }}</span>
                <span>Y {{ formatCoordinate(base.location_y) }}</span>
              </div>
              <span class="base-id">ID {{ base.id || "—" }}</span>
            </div>
          </article>
        </div>
        <n-empty v-else size="small" :description="$t('guild.noBases')" />
      </section>
    </div>
  </article>
</template>

<style scoped lang="less">
.guild-overview {
  --panel-border: rgba(24, 24, 28, 0.08);
  --panel-bg: rgba(24, 24, 28, 0.025);
  --muted: rgba(24, 24, 28, 0.56);
  --subtle: rgba(24, 24, 28, 0.4);
  min-width: 0;
  padding: 28px;
}

.guild-overview.is-dark {
  --panel-border: rgba(255, 255, 255, 0.09);
  --panel-bg: rgba(255, 255, 255, 0.04);
  --muted: rgba(255, 255, 255, 0.58);
  --subtle: rgba(255, 255, 255, 0.4);
}

.guild-hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  padding: 4px 2px 24px;
}

.guild-heading {
  min-width: 0;
}

.guild-eyebrow,
.section-kicker {
  display: block;
  color: #4098fc;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  line-height: 1.4;
  text-transform: uppercase;
}

.guild-title {
  margin: 5px 0 8px;
  font-size: clamp(24px, 2.4vw, 34px);
  font-weight: 720;
  line-height: 1.22;
  overflow-wrap: anywhere;
}

.guild-master-line {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 7px;
  margin: 0;
  color: var(--muted);
  font-size: 14px;
}

.guild-master-line strong {
  color: inherit;
  font-weight: 650;
}

.guild-level-badge {
  display: flex;
  flex: none;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  padding: 10px 14px;
  border: 1px solid rgba(64, 152, 252, 0.3);
  border-radius: 12px;
  background: rgba(64, 152, 252, 0.09);
}

.guild-level-badge span {
  color: var(--muted);
  font-size: 11px;
}

.guild-level-badge strong {
  color: #4098fc;
  font-size: 20px;
  line-height: 1.1;
}

.guild-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 18px;
}

.guild-stat-card {
  min-width: 0;
  padding: 14px 16px;
  border: 1px solid var(--panel-border);
  border-radius: 12px;
  background: var(--panel-bg);
}

.guild-stat-card span {
  display: block;
  margin-bottom: 3px;
  color: var(--muted);
  font-size: 12px;
}

.guild-stat-card strong {
  display: block;
  overflow: hidden;
  font-size: 21px;
  line-height: 1.25;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.guild-stat-master strong {
  font-size: 17px;
}

.guild-content-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(min(100%, 360px), 1fr));
  align-items: start;
  gap: 16px;
}

.guild-section {
  min-width: 0;
  padding: 18px;
  border: 1px solid var(--panel-border);
  border-radius: 16px;
  background: var(--panel-bg);
}

.section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 14px;
}

.section-heading h2 {
  margin: 3px 0 0;
  font-size: 19px;
  font-weight: 680;
  line-height: 1.3;
}

.member-grid,
.base-list {
  display: grid;
  gap: 8px;
}

.member-card,
.base-card {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  padding: 12px;
  border: 1px solid var(--panel-border);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.62);
}

.is-dark .member-card,
.is-dark .base-card {
  background: rgba(0, 0, 0, 0.12);
}

.member-card-main,
.base-card-content {
  flex: 1;
  min-width: 0;
}

.member-name-line {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
}

.member-name {
  min-width: 0;
  font-size: 16px;
  font-weight: 660;
  line-height: 1.35;
  overflow-wrap: anywhere;
}

.member-uid,
.base-id {
  display: block;
  margin-top: 5px;
  overflow: hidden;
  color: var(--subtle);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.view-player-button {
  flex: none;
}

.quick-map-wrap {
  overflow: hidden;
  margin-bottom: 10px;
  border: 1px solid var(--panel-border);
  border-radius: 12px;
  background: rgba(64, 152, 252, 0.06);
}

.quick-map-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
}

.quick-map-heading strong,
.quick-map-heading span {
  display: block;
}

.quick-map-heading strong {
  font-size: 13px;
}

.quick-map-heading span,
.map-axis {
  margin-top: 2px;
  color: var(--muted);
  font-size: 11px;
}

.quick-map {
  position: relative;
  aspect-ratio: 4 / 3;
  overflow: hidden;
  border-top: 1px solid var(--panel-border);
  background:
    radial-gradient(
      circle at 62% 42%,
      rgba(64, 152, 252, 0.22),
      transparent 22%
    ),
    linear-gradient(145deg, #182536, #0f1722);
}

.quick-map-background {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0.68;
  object-fit: fill;
  filter: saturate(0.72) contrast(1.04) brightness(0.72);
}

.quick-map-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.11) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.11) 1px, transparent 1px);
  background-size: 25% 50%;
}

.base-marker {
  position: absolute;
  z-index: 2;
  width: 25px;
  height: 25px;
  padding: 0;
  transform: translate(-50%, -50%);
  border: 2px solid #fff;
  border-radius: 50%;
  outline: none;
  background: #f0a020;
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.42);
  color: #fff;
  cursor: pointer;
  font-size: 11px;
  font-weight: 800;
  line-height: 21px;
  transition:
    transform 0.16s ease,
    background-color 0.16s ease;
}

.base-marker:hover,
.base-marker:focus-visible,
.base-marker.is-active {
  z-index: 3;
  transform: translate(-50%, -50%) scale(1.18);
  background: #4098fc;
}

.base-card {
  align-items: flex-start;
  transition:
    border-color 0.16s ease,
    background-color 0.16s ease;
}

.base-card.is-active {
  border-color: rgba(64, 152, 252, 0.55);
  background: rgba(64, 152, 252, 0.09);
}

.base-card-index {
  display: grid;
  width: 28px;
  height: 28px;
  flex: none;
  place-items: center;
  border-radius: 8px;
  background: rgba(240, 160, 32, 0.14);
  color: #f0a020;
  font-size: 13px;
  font-weight: 800;
}

.base-card-title,
.base-coordinate {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.base-card-title strong {
  font-size: 15px;
}

.base-card-title span {
  color: var(--muted);
  font-size: 11px;
}

.base-coordinate {
  justify-content: flex-start;
  flex-wrap: wrap;
  margin-top: 5px;
  color: var(--muted);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.guild-overview.is-compact {
  padding: 12px 12px 24px;
}

.is-compact .guild-hero {
  gap: 12px;
  padding: 4px 2px 18px;
}

.is-compact .guild-title {
  font-size: 24px;
}

.is-compact .guild-level-badge {
  padding: 9px 11px;
}

.is-compact .guild-stats {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.is-compact .guild-stat-master {
  grid-column: 1 / -1;
}

.is-compact .guild-section {
  padding: 14px;
}

.is-compact .member-card {
  align-items: flex-start;
  flex-direction: column;
}

.is-compact .view-player-button {
  width: 100%;
}

@media (max-width: 700px) {
  .guild-overview:not(.is-compact) {
    padding: 18px;
  }

  .guild-stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .guild-stat-master {
    grid-column: 1 / -1;
  }
}
</style>
