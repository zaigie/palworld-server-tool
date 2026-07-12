<script setup>
import { computed, onMounted, ref } from "vue";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";

const props = defineProps({
  serverInfo: { type: Object, default: () => ({}) },
  serverMetrics: { type: Object, default: () => ({}) },
  players: { type: Array, default: () => [] },
});
const emit = defineEmits([
  "open-rcon",
  "open-backup",
  "open-broadcast",
  "open-config",
]);
const { t } = useI18n();
const api = new ApiService();
const loading = ref(false);
const onlinePlayers = ref([]);
const backups = ref([]);
const tasks = ref([]);
const asArray = (value) => (Array.isArray(value) ? value : []);

const latestBackup = computed(
  () =>
    [...backups.value].sort(
      (a, b) => new Date(b.save_time) - new Date(a.save_time),
    )[0],
);
const activeTasks = computed(() => tasks.value.filter((task) => task.enabled));
const nextTask = computed(
  () =>
    activeTasks.value
      .filter((task) => task.next_run_at)
      .sort((a, b) => new Date(a.next_run_at) - new Date(b.next_run_at))[0],
);
const uptime = computed(() => {
  const seconds = Number(props.serverMetrics?.uptime || 0);
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  return days > 0
    ? t("overview.uptimeDays", { days, hours })
    : t("overview.uptimeHours", { hours });
});
const healthType = computed(() => {
  if (!props.serverInfo?.name) return "warning";
  const fps = Number(props.serverMetrics?.server_fps || 0);
  return fps > 0 && fps < 30 ? "warning" : "success";
});

const loadOverview = async () => {
  loading.value = true;
  try {
    const [onlineResponse, backupResponse, taskResponse] = await Promise.all([
      api.getOnlinePlayerList(),
      api.getBackupList({}),
      api.getRconTasks(),
    ]);
    onlinePlayers.value = asArray(onlineResponse.data.value);
    backups.value = asArray(backupResponse.data.value);
    tasks.value = asArray(taskResponse.data.value);
  } finally {
    loading.value = false;
  }
};

const formatTime = (value) =>
  value ? dayjs(value).format("YYYY-MM-DD HH:mm:ss") : "—";

onMounted(loadOverview);
</script>

<template>
  <n-scrollbar class="h-full">
    <div class="p-5 max-w-1400px mx-auto">
      <n-flex justify="space-between" align="center" class="mb-4">
        <div>
          <n-h2 class="m-0">{{ $t("overview.title") }}</n-h2>
          <n-text depth="3">{{ $t("overview.subtitle") }}</n-text>
        </div>
        <n-button secondary :loading="loading" @click="loadOverview">
          {{ $t("overview.refresh") }}
        </n-button>
      </n-flex>

      <n-grid cols="1 640:2 1050:4" :x-gap="16" :y-gap="16">
        <n-gi>
          <n-card size="small">
            <n-statistic :label="$t('overview.serverStatus')">
              <n-flex align="center">
                <n-badge dot :type="healthType" />
                <n-text strong>{{
                  serverInfo?.name || $t("status.serverUnavailable")
                }}</n-text>
              </n-flex>
            </n-statistic>
            <n-text depth="3">{{ serverInfo?.version || "—" }}</n-text>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic
              :label="$t('overview.onlinePlayers')"
              :value="serverMetrics?.current_player_num ?? onlinePlayers.length"
            >
              <template #suffix
                >/ {{ serverMetrics?.max_player_num ?? "—" }}</template
              >
            </n-statistic>
            <n-text depth="3">{{
              $t("overview.totalPlayers", { count: players.length })
            }}</n-text>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic
              :label="$t('item.serverFps')"
              :value="serverMetrics?.server_fps ?? '—'"
            />
            <n-text depth="3"
              >{{ $t("item.serverFrameTime") }}:
              {{ serverMetrics?.server_frame_time ?? "—" }} ms</n-text
            >
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic :label="$t('item.serverUptime')" :value="uptime" />
            <n-text depth="3"
              >{{ $t("item.serverDays") }}:
              {{ serverMetrics?.days ?? "—" }}</n-text
            >
          </n-card>
        </n-gi>
      </n-grid>

      <n-grid cols="1 760:2" :x-gap="16" :y-gap="16" class="mt-4">
        <n-gi>
          <n-card :title="$t('overview.operations')">
            <n-grid cols="2 560:4" :x-gap="12" :y-gap="12">
              <n-gi
                ><n-button
                  block
                  type="primary"
                  secondary
                  @click="emit('open-rcon')"
                  >{{ $t("button.rcon") }}</n-button
                ></n-gi
              >
              <n-gi
                ><n-button
                  block
                  type="success"
                  secondary
                  @click="emit('open-backup')"
                  >{{ $t("button.backup") }}</n-button
                ></n-gi
              >
              <n-gi
                ><n-button
                  block
                  type="warning"
                  secondary
                  @click="emit('open-broadcast')"
                  >{{ $t("button.broadcast") }}</n-button
                ></n-gi
              >
              <n-gi
                ><n-button block secondary @click="emit('open-config')">{{
                  $t("configuration.title")
                }}</n-button></n-gi
              >
            </n-grid>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card :title="$t('overview.automation')">
            <n-descriptions :column="1" label-placement="left">
              <n-descriptions-item :label="$t('overview.activeTasks')">{{
                activeTasks.length
              }}</n-descriptions-item>
              <n-descriptions-item :label="$t('overview.nextTask')">{{
                nextTask?.name || "—"
              }}</n-descriptions-item>
              <n-descriptions-item :label="$t('overview.nextRun')">{{
                formatTime(nextTask?.next_run_at)
              }}</n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-gi>
      </n-grid>

      <n-grid cols="1 760:2" :x-gap="16" :y-gap="16" class="mt-4">
        <n-gi>
          <n-card :title="$t('overview.onlineNow')">
            <n-empty
              v-if="onlinePlayers.length === 0"
              :description="$t('overview.noOnlinePlayers')"
            />
            <n-list v-else hoverable>
              <n-list-item
                v-for="player in onlinePlayers.slice(0, 6)"
                :key="player.player_uid"
              >
                <n-flex justify="space-between">
                  <n-text>{{ player.nickname }}</n-text>
                  <n-tag size="small" type="success"
                    >Lv.{{ player.level }}</n-tag
                  >
                </n-flex>
              </n-list-item>
            </n-list>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card :title="$t('overview.backupStatus')">
            <n-empty
              v-if="!latestBackup"
              :description="$t('overview.noBackup')"
            />
            <n-descriptions v-else :column="1" label-placement="left">
              <n-descriptions-item :label="$t('overview.latestBackup')">{{
                formatTime(latestBackup.save_time)
              }}</n-descriptions-item>
              <n-descriptions-item :label="$t('overview.backupCount')">{{
                backups.length
              }}</n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </n-scrollbar>
</template>
