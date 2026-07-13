<script setup>
import { useI18n } from "vue-i18n";
import "leaflet/dist/leaflet.css";
import {
  LCircle,
  LIcon,
  LMap,
  LMarker,
  LPopup,
  LTileLayer,
  LTooltip,
} from "@vue-leaflet/vue-leaflet";
import { AddCircle20Filled, SubtractCircle20Filled } from "@vicons/fluent";
import ApiService from "@/service/api.js";
import IconBase from "@/assets/map/base.webp";
import IconPlayer from "@/assets/map/player.webp";
import IconBossTower from "@/assets/map/boss_tower.webp";
import IconFastTravel from "@/assets/map/fast_travel.webp";
import playerToGuildStore from "@/stores/model/playerToGuild.js";
import points from "@/assets/map/points.json";

const { t } = useI18n();

const LAND_SCAPE = [349400, 724400, -1099400, -724400];

const api = new ApiService();

const mousePosition = ref([0, 0]);
const zoom = ref(2);
const tiles = ref("map/tiles/{z}/{x}/{y}.png");
const playerList = ref([]);
const guildList = ref([]);
const showPlayer = ref(true);
const showBaseCamp = ref(true);
const showBossTower = ref(false);
const showFastTravel = ref(false);
const mapRef = ref(null);
const searchTarget = ref(null);

let timer = null;

const toMapPosition = (position) => {
  // hack
  if (position[0] >= -256 && position[0] <= 256) {
    return position;
  }
  const x =
    -256 +
    (256 * (position[0] - LAND_SCAPE[2])) / (LAND_SCAPE[0] - LAND_SCAPE[2]);
  const y =
    (256 * (position[1] - LAND_SCAPE[3])) / (LAND_SCAPE[1] - LAND_SCAPE[3]);
  return [x, y];
};

const fromMapPosition = (mapPosition) => {
  // 还原 x 坐标
  const worldX =
    ((mapPosition[0] + 256) * (LAND_SCAPE[0] - LAND_SCAPE[2])) / 256 +
    LAND_SCAPE[2];
  // 还原 y 坐标
  const worldY =
    (mapPosition[1] * (LAND_SCAPE[1] - LAND_SCAPE[3])) / 256 + LAND_SCAPE[3];

  // 保留两位小数
  return [worldX.toFixed(2), worldY.toFixed(2)];
};

const toMapDistance = (distance) => {
  return 256 * (distance / (LAND_SCAPE[0] - LAND_SCAPE[2]));
};

const rawBaseMarkers = computed(() =>
  guildList.value.flatMap((guild, guildIndex) =>
    (guild.base_camp || []).map((camp, campIndex) => ({
      key: `${guildIndex}-${campIndex}`,
      guild,
      camp,
      position: toMapPosition([camp.location_x, camp.location_y]),
    })),
  ),
);

const clusteredBaseMarkers = computed(() => {
  if (zoom.value >= 5) {
    return rawBaseMarkers.value.map((marker) => ({
      key: marker.key,
      position: marker.position,
      markers: [marker],
    }));
  }

  const cellSize = 48 / 2 ** zoom.value;
  const clusters = new Map();
  rawBaseMarkers.value.forEach((marker) => {
    const clusterKey = `${Math.floor(marker.position[0] / cellSize)}:${Math.floor(
      marker.position[1] / cellSize,
    )}`;
    if (!clusters.has(clusterKey)) clusters.set(clusterKey, []);
    clusters.get(clusterKey).push(marker);
  });

  return Array.from(clusters.entries()).map(([key, markers]) => ({
    key,
    markers,
    position: [
      markers.reduce((sum, marker) => sum + marker.position[0], 0) /
        markers.length,
      markers.reduce((sum, marker) => sum + marker.position[1], 0) /
        markers.length,
    ],
  }));
});

const searchOptions = computed(() => [
  ...playerList.value.map((player, index) => ({
    label: `${t("status.online")} · ${player.nickname || player.player_uid}`,
    value: `player:${index}`,
  })),
  ...rawBaseMarkers.value.map((marker) => ({
    label: `${t("map.baseCamp")} · ${marker.guild.name || t("filter.unnamedGuild")}`,
    value: `base:${marker.key}`,
  })),
]);

const focusSearchTarget = (value) => {
  if (!value) return;
  const [type, id] = value.split(":");
  let position;
  if (type === "player") {
    const player = playerList.value[Number(id)];
    if (player)
      position = toMapPosition([player.location_x, player.location_y]);
  } else {
    position = rawBaseMarkers.value.find(
      (marker) => marker.key === id,
    )?.position;
  }
  if (position) {
    zoom.value = 5;
    mapRef.value?.leafletObject?.setView(position, 5);
  }
};

const setMarkerAccessibility = (marker, label) => {
  const element = marker?.getElement?.();
  if (!element) return;
  element.setAttribute("aria-label", label);
  element.setAttribute("alt", label);
};

const ToPlayers = async (uid) => {
  playerToGuildStore().setCurrentUid(uid);
  playerToGuildStore().setUpdateStatus("players");
};

const refreshPlayer = async () => {
  const { data } = await api.getOnlinePlayerList();
  const onlinePlayers = Array.isArray(data.value) ? data.value : [];
  for (const i of onlinePlayers) {
    for (const j of playerList.value) {
      if (i.player_uid === j.player_uid) {
        j.location_x = i.location_x;
        j.location_y = i.location_y;
        break;
      }
    }
  }
  timer = setTimeout(refreshPlayer, 5000);
};

const onMapMouseMove = (event) => {
  mousePosition.value = [
    event.latlng.lat.toFixed(2),
    event.latlng.lng.toFixed(2),
  ];
};

// 左下角控件
const onAddZoom = () => {
  if (zoom.value !== 6) {
    zoom.value += 1;
  }
};
const onSubtractZoom = () => {
  if (zoom.value !== 0) {
    zoom.value -= 1;
  }
};

onMounted(async () => {
  let res = await api.getPlayerList({});
  playerList.value = Array.isArray(res.data.value) ? res.data.value : [];
  // 接口中玩家location_x和location_y同时为0时，表示玩家不在线，不显示
  playerList.value = playerList.value.filter(
    (i) => i.location_x !== 0 && i.location_y !== 0,
  );
  res = await api.getGuildList();
  guildList.value = Array.isArray(res.data.value) ? res.data.value : [];

  refreshPlayer();
});

onUnmounted(async () => {
  clearTimeout(timer);
});
</script>

<template>
  <div class="map-view h-full">
    <l-map
      ref="mapRef"
      style="width: 100%; height: 100%"
      crs="Simple"
      v-model:zoom="zoom"
      :use-global-leaflet="false"
      :center="[-128, 128]"
      :min-zoom="0"
      :max-zoom="6"
      :options="{ zoomControl: false, attributionControl: false }"
      @mousemove="onMapMouseMove"
    >
      <l-tile-layer
        :url="tiles"
        no-wrap
        layer-type="base"
        :options="{
          bounds: [
            [0, 0],
            [-256, 256],
          ],
        }"
      ></l-tile-layer>
      <l-marker
        v-if="showFastTravel"
        v-for="i in points.fast_travel"
        :key="`fast-${i[0]}-${i[1]}`"
        :lat-lng="toMapPosition([i[0], i[1]])"
        :options="{ title: $t('map.fastTravel'), alt: $t('map.fastTravel') }"
        @ready="
          (marker) => setMarkerAccessibility(marker, $t('map.fastTravel'))
        "
      >
        <l-icon :icon-url="IconFastTravel" :icon-size="[48, 48]" />
      </l-marker>
      <l-marker
        v-if="showBossTower"
        v-for="i in points.boss_tower"
        :key="`tower-${i[0]}-${i[1]}`"
        :lat-lng="toMapPosition([i[0], i[1]])"
        :options="{ title: $t('map.bossTower'), alt: $t('map.bossTower') }"
        @ready="(marker) => setMarkerAccessibility(marker, $t('map.bossTower'))"
      >
        <l-icon :icon-url="IconBossTower" :icon-size="[48, 48]" />
      </l-marker>
      <l-marker
        v-if="showPlayer"
        v-for="i in playerList"
        :key="i.player_uid"
        :lat-lng="toMapPosition([i.location_x, i.location_y])"
        :options="{
          title: `${$t('map.player')}: ${i.nickname}`,
          alt: `${$t('map.player')}: ${i.nickname}`,
        }"
        @ready="
          (marker) =>
            setMarkerAccessibility(marker, `${$t('map.player')}: ${i.nickname}`)
        "
      >
        <l-icon :icon-url="IconPlayer" :icon-size="[45, 45]" />
        <l-tooltip
          :options="{ direction: 'top', permanent: true, offset: [0, -15] }"
          >{{ i.nickname }}</l-tooltip
        >
      </l-marker>
      <template v-if="showBaseCamp">
        <l-marker
          v-for="cluster in clusteredBaseMarkers"
          :key="cluster.key"
          :lat-lng="cluster.position"
          :options="{
            title:
              cluster.markers.length > 1
                ? $t('map.clusterTitle', { count: cluster.markers.length })
                : $t('map.baseCampTitle', {
                    name:
                      cluster.markers[0].guild.name ||
                      $t('filter.unnamedGuild'),
                  }),
            alt:
              cluster.markers.length > 1
                ? $t('map.clusterTitle', { count: cluster.markers.length })
                : $t('map.baseCampTitle', {
                    name:
                      cluster.markers[0].guild.name ||
                      $t('filter.unnamedGuild'),
                  }),
          }"
          @ready="
            (marker) =>
              setMarkerAccessibility(
                marker,
                cluster.markers.length > 1
                  ? $t('map.clusterTitle', { count: cluster.markers.length })
                  : $t('map.baseCampTitle', {
                      name:
                        cluster.markers[0].guild.name ||
                        $t('filter.unnamedGuild'),
                    }),
              )
          "
        >
          <l-icon
            :icon-url="IconBase"
            :icon-size="cluster.markers.length > 1 ? [62, 62] : [55, 55]"
          />
          <l-tooltip
            v-if="cluster.markers.length > 1"
            :options="{ direction: 'top', permanent: true, offset: [0, -18] }"
          >
            {{ cluster.markers.length }}
          </l-tooltip>
          <l-popup :options="{ interactive: true }">
            <div v-if="cluster.markers.length > 1" class="popup-title">
              {{ $t("map.clusterTitle", { count: cluster.markers.length }) }}
            </div>
            <div
              v-for="marker in cluster.markers"
              :key="marker.key"
              class="base-popup"
            >
              <div class="popup-title">
                {{
                  $t("map.baseCampTitle", {
                    name: marker.guild.name || $t("filter.unnamedGuild"),
                  })
                }}
              </div>
              <div class="member-list">
                {{ $t("map.guildMember") }}
                <button
                  v-for="member in marker.guild.players"
                  :key="member.player_uid"
                  type="button"
                  class="player_name"
                  @click="ToPlayers(member.player_uid)"
                >
                  {{ member.nickname }}
                </button>
              </div>
            </div>
          </l-popup>
        </l-marker>
        <template v-if="zoom >= 5">
          <l-circle
            v-for="marker in rawBaseMarkers"
            :key="`area-${marker.key}`"
            :lat-lng="marker.position"
            :radius="toMapDistance(marker.camp.area)"
          />
        </template>
      </template>
    </l-map>
    <div
      class="min-h-50 p-2 fixed bottom-2 left-2 z-999 flex flex-col justify-end"
    >
      <div class="h-40 flex flex-col justify-between items-center">
        <n-button
          text
          :aria-label="$t('map.zoomIn')"
          size="24"
          color="#fff"
          @click="onAddZoom"
        >
          <template #icon
            ><n-icon><AddCircle20Filled /></n-icon
          ></template>
        </n-button>
        <n-slider
          style="height: 100px"
          class="border border-solid border-#fff rounded-full"
          v-model:value="zoom"
          :tooltip="false"
          :height="4"
          :step="1"
          :min="0"
          :max="6"
          vertical
        />
        <n-button
          text
          :aria-label="$t('map.zoomOut')"
          size="24"
          color="#fff"
          @click="onSubtractZoom"
        >
          <template #icon
            ><n-icon><SubtractCircle20Filled /></n-icon
          ></template>
        </n-button>
      </div>
    </div>
    <div class="control">
      <n-select
        v-model:value="searchTarget"
        :options="searchOptions"
        :placeholder="$t('map.searchTarget')"
        filterable
        clearable
        :aria-label="$t('map.searchTarget')"
        @update:value="focusSearchTarget"
      />
      <div class="visible-summary">
        {{
          $t("map.visibleSummary", {
            players: showPlayer ? playerList.length : 0,
            bases: showBaseCamp ? clusteredBaseMarkers.length : 0,
          })
        }}
      </div>
      <div>
        <span>{{ $t("map.showFastTravel") }}</span>
        <n-switch
          v-model:value="showFastTravel"
          :aria-label="$t('map.showFastTravel')"
        />
      </div>
      <div>
        <span>{{ $t("map.showBossTower") }}</span>
        <n-switch
          v-model:value="showBossTower"
          :aria-label="$t('map.showBossTower')"
        />
      </div>
      <div>
        <span>{{ $t("map.showPlayer") }}</span>
        <n-switch
          v-model:value="showPlayer"
          :aria-label="$t('map.showPlayer')"
        />
      </div>
      <div>
        <span>{{ $t("map.showBaseCamp") }}</span>
        <n-switch
          v-model:value="showBaseCamp"
          :aria-label="$t('map.showBaseCamp')"
        />
      </div>
      <div>
        <span>{{ mousePosition[0] }}, {{ mousePosition[1] }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.leaflet-container {
  background: #102536;
  outline: 0;
}

.player_name {
  border: 0;
  cursor: pointer;
  font: inherit;
  margin: 0 3px;
  padding: 3px;
  color: #fff;
  background-color: #009f5d;
  border-radius: 3px;
}

.control {
  width: 260px;
  max-width: calc(100% - 40px);
  min-height: 230px;
  position: absolute;
  bottom: 20px;
  right: 20px;
  padding: 10px;
  box-sizing: border-box;
  color: #fff;
  background-color: rgb(24, 24, 28);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  z-index: 999;
}

.control > .n-select {
  width: 100%;
  min-width: 0;
  flex: none;
}

.control > div {
  display: flex;
  justify-content: space-between;
  min-width: 0;
  margin: 0;
  padding: 8px 0;
}

.visible-summary {
  justify-content: flex-start !important;
  color: #b8c2cc;
  font-size: 12px;
}

@media (max-width: 480px) {
  .control {
    right: 16px;
    bottom: 16px;
    width: min(260px, calc(100% - 32px));
    max-width: calc(100% - 32px);
  }
}

.popup-title {
  padding-bottom: 3px;
  font-size: 16px;
  font-weight: 600;
}

.base-popup + .base-popup {
  margin-top: 10px;
  padding-top: 8px;
  border-top: 1px solid #ddd;
}

.member-list {
  line-height: 28px;
}
</style>
