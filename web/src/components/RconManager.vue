<script setup>
import { computed, ref, watch } from "vue";
import dayjs from "dayjs";
import { useDialog, useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import itemMap from "@/assets/items.json";
import palMap from "@/assets/pal.json";
import {
  RCON_PLACEHOLDERS,
  describeCron,
  extractRconPlaceholders,
  resolveRconTemplate,
} from "@/utils/rconTemplate";

const props = defineProps({ show: Boolean });
const emit = defineEmits(["update:show"]);
const { t, locale } = useI18n();
const message = useMessage();
const dialog = useDialog();
const api = new ApiService();
const asArray = (value) => (Array.isArray(value) ? value : []);

const visible = computed({
  get: () => props.show,
  set: (value) => emit("update:show", value),
});

const loading = ref(false);
const activeTab = ref("commands");
const commands = ref([]);
const tasks = ref([]);
const players = ref([]);
const commandContent = ref({});
const drawerWidth = computed(() => Math.min(680, window.innerWidth));

const selectedPlayerUid = ref(null);
const selectedItem = ref(null);
const selectedPal = ref(null);
const selectedPlayer = computed(() =>
  players.value.find((player) => player.player_uid === selectedPlayerUid.value),
);
const playerOptions = computed(() =>
  players.value.map((player) => ({
    label: `${player.nickname} (${player.player_uid})`,
    value: player.player_uid,
  })),
);
const itemOptions = computed(() =>
  (itemMap[locale.value] || itemMap.zh).map((item) => ({
    label: `${item.name} · ${item.key}`,
    value: item.key,
  })),
);
const palOptions = computed(() =>
  Object.entries(palMap[locale.value] || palMap.zh).map(([key, value]) => ({
    label: `${value} · ${key}`,
    value: key,
  })),
);
const commandOptions = computed(() =>
  commands.value.map((command) => ({
    label: command.remark || command.command,
    value: command.uuid,
  })),
);
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem("palworld_token") || ""}`,
}));

const placeholderLabels = computed(() =>
  RCON_PLACEHOLDERS.map((placeholder) => ({
    ...placeholder,
    token: `{${placeholder.key}}`,
    label: t(`rconManager.placeholder.${placeholder.key}`),
  })),
);

const loadData = async () => {
  loading.value = true;
  try {
    const [commandResponse, taskResponse, playerResponse] = await Promise.all([
      api.getRconCommands(),
      api.getRconTasks(),
      api.getPlayerList({ order_by: "last_online", desc: true }),
    ]);
    commands.value = asArray(commandResponse.data.value);
    tasks.value = asArray(taskResponse.data.value);
    players.value = asArray(playerResponse.data.value);
    for (const command of commands.value) {
      if (commandContent.value[command.uuid] === undefined) {
        commandContent.value[command.uuid] = "";
      }
    }
  } finally {
    loading.value = false;
  }
};

watch(
  () => props.show,
  (show) => {
    if (show) loadData();
  },
);

const resolveTemplate = (template) =>
  resolveRconTemplate(template, {
    player: selectedPlayer.value,
    item: selectedItem.value,
    pal: selectedPal.value,
  });

const fillCommand = (command) => {
  const result = resolveTemplate(command.placeholder || "");
  commandContent.value[command.uuid] = result.content;
  if (result.missing.length > 0) {
    message.warning(
      t("rconManager.missingSelection", {
        fields: result.missing
          .map((key) => t(`rconManager.placeholder.${key}`))
          .join(", "),
      }),
    );
  } else if (result.unknown.length > 0) {
    message.warning(
      t("rconManager.unknownPlaceholder", {
        fields: result.unknown.join(", "),
      }),
    );
  } else {
    message.success(t("rconManager.filled"));
  }
};

const executeCommand = async (command) => {
  const { data, statusCode } = await api.sendRconCommand({
    uuid: command.uuid,
    content: commandContent.value[command.uuid] || "",
  });
  if (statusCode.value === 200) {
    message.success(data.value?.message || t("rconManager.executed"));
  } else {
    message.error(data.value?.error || t("rconManager.executeFailed"));
  }
};

const commandModal = ref(false);
const editingCommandUUID = ref(null);
const commandForm = ref({ command: "", remark: "", placeholder: "" });
const openCommandModal = (command = null) => {
  editingCommandUUID.value = command?.uuid || null;
  commandForm.value = command
    ? {
        command: command.command,
        remark: command.remark,
        placeholder: command.placeholder,
      }
    : { command: "", remark: "", placeholder: "" };
  commandModal.value = true;
};
const saveCommand = async () => {
  if (!commandForm.value.command.trim() || !commandForm.value.remark.trim()) {
    message.warning(t("rconManager.commandRequired"));
    return;
  }
  const response = editingCommandUUID.value
    ? await api.putRconCommand(editingCommandUUID.value, commandForm.value)
    : await api.addRconCommand(commandForm.value);
  if (response.statusCode.value === 200) {
    message.success(t("rconManager.commandSaved"));
    commandModal.value = false;
    await loadData();
  } else {
    message.error(response.data.value?.error || t("rconManager.saveFailed"));
  }
};
const removeCommand = (command) => {
  dialog.warning({
    title: t("message.warn"),
    content: t("rconManager.removeCommandConfirm", {
      name: command.remark || command.command,
    }),
    positiveText: t("button.confirm"),
    negativeText: t("button.cancel"),
    onPositiveClick: async () => {
      const { data, statusCode } = await api.removeRconCommand(command.uuid);
      if (statusCode.value === 200) {
        message.success(t("rconManager.commandRemoved"));
        await loadData();
      } else {
        message.error(data.value?.error || t("rconManager.removeFailed"));
      }
    },
  });
};
const importFinished = async () => {
  message.success(t("message.importRconSuccess"));
  commandModal.value = false;
  await loadData();
};
const importFailed = () => message.error(t("message.importRconFail", { err: "" }));

const taskModal = ref(false);
const editingTaskUUID = ref(null);
const scheduleType = ref("interval");
const intervalMinutes = ref(15);
const dailyTime = ref("04:00");
const weeklyDay = ref(1);
const weeklyTime = ref("04:00");
const customCron = ref("0 4 * * *");
const taskForm = ref({
  name: "",
  rcon_uuid: null,
  content: "",
  cron: "",
  enabled: true,
});
const scheduleTypeOptions = computed(() => [
  { label: t("rconManager.schedule.interval"), value: "interval" },
  { label: t("rconManager.schedule.daily"), value: "daily" },
  { label: t("rconManager.schedule.weekly"), value: "weekly" },
  { label: t("rconManager.schedule.custom"), value: "custom" },
]);
const weekOptions = computed(() =>
  Array.from({ length: 7 }, (_, index) => ({
    value: index,
    label: t(`rconManager.weekday.${index}`),
  })),
);

const parseTime = (value) => {
  const [hour, minute] = (value || "00:00").split(":").map(Number);
  return { hour: Number.isFinite(hour) ? hour : 0, minute: Number.isFinite(minute) ? minute : 0 };
};
const buildCron = () => {
  if (scheduleType.value === "interval") {
    return `*/${Math.max(1, Number(intervalMinutes.value) || 1)} * * * *`;
  }
  if (scheduleType.value === "daily") {
    const { hour, minute } = parseTime(dailyTime.value);
    return `${minute} ${hour} * * *`;
  }
  if (scheduleType.value === "weekly") {
    const { hour, minute } = parseTime(weeklyTime.value);
    return `${minute} ${hour} * * ${weeklyDay.value}`;
  }
  return customCron.value.trim();
};

const selectTaskCommand = (commandUUID) => {
  const command = commands.value.find((item) => item.uuid === commandUUID);
  if (!command) return;
  const resolved = resolveTemplate(command.placeholder || "");
  taskForm.value.content = resolved.content;
  if (!taskForm.value.name) taskForm.value.name = command.remark || command.command;
};

watch([selectedPlayerUid, selectedItem, selectedPal], () => {
  for (const command of commands.value) {
    if (command.placeholder) {
      commandContent.value[command.uuid] = resolveTemplate(
        command.placeholder
      ).content;
    }
  }
  if (taskModal.value && taskForm.value.rcon_uuid) {
    selectTaskCommand(taskForm.value.rcon_uuid);
  }
});

const openTaskModal = (command = null, existingTask = null) => {
  editingTaskUUID.value = existingTask?.uuid || null;
  if (existingTask) {
    taskForm.value = {
      name: existingTask.name,
      rcon_uuid: existingTask.rcon_uuid,
      content: existingTask.content,
      cron: existingTask.cron,
      enabled: existingTask.enabled,
    };
    const described = describeCron(existingTask.cron);
    const intervalMatch = existingTask.cron.match(/^\*\/(\d+) \* \* \* \*$/);
    const weeklyMatch = existingTask.cron.match(/^(\d+) (\d+) \* \* ([0-6])$/);
    const dailyMatch = existingTask.cron.match(/^(\d+) (\d+) \* \* \*$/);
    if (intervalMatch) {
      scheduleType.value = "interval";
      intervalMinutes.value = Number(intervalMatch[1]);
    } else if (weeklyMatch) {
      scheduleType.value = "weekly";
      weeklyTime.value = `${String(weeklyMatch[2]).padStart(2, "0")}:${String(weeklyMatch[1]).padStart(2, "0")}`;
      weeklyDay.value = Number(weeklyMatch[3]);
    } else if (dailyMatch) {
      scheduleType.value = "daily";
      dailyTime.value = `${String(dailyMatch[2]).padStart(2, "0")}:${String(dailyMatch[1]).padStart(2, "0")}`;
    } else {
      scheduleType.value = described === "custom" ? "custom" : "interval";
      customCron.value = existingTask.cron;
    }
  } else {
    taskForm.value = {
      name: command?.remark || command?.command || "",
      rcon_uuid: command?.uuid || null,
      content: "",
      cron: "",
      enabled: true,
    };
    scheduleType.value = "interval";
    intervalMinutes.value = 15;
    if (command) selectTaskCommand(command.uuid);
  }
  taskModal.value = true;
};

const saveTask = async () => {
  const cron = buildCron();
  if (!taskForm.value.name.trim() || !taskForm.value.rcon_uuid || !cron) {
    message.warning(t("rconManager.taskRequired"));
    return;
  }
  const unresolved = extractRconPlaceholders(taskForm.value.content);
  if (unresolved.length > 0) {
    message.warning(
      t("rconManager.unresolvedTaskPlaceholder", {
        fields: unresolved.map((key) => `{${key}}`).join(", "),
      }),
    );
    return;
  }
  const payload = { ...taskForm.value, cron };
  const response = editingTaskUUID.value
    ? await api.putRconTask(editingTaskUUID.value, payload)
    : await api.addRconTask(payload);
  if (response.statusCode.value === 200) {
    message.success(t("rconManager.taskSaved"));
    taskModal.value = false;
    activeTab.value = "tasks";
    await loadData();
  } else {
    message.error(response.data.value?.error || t("rconManager.saveFailed"));
  }
};

const toggleTask = async (task, enabled) => {
  const response = await api.putRconTask(task.uuid, { ...task, enabled });
  if (response.statusCode.value === 200) {
    message.success(t(enabled ? "rconManager.taskEnabled" : "rconManager.taskPaused"));
    await loadData();
  } else {
    message.error(response.data.value?.error || t("rconManager.saveFailed"));
  }
};
const runTask = (task) => {
  dialog.warning({
    title: t("rconManager.runNow"),
    content: t("rconManager.runConfirm", { name: task.name }),
    positiveText: t("button.execute"),
    negativeText: t("button.cancel"),
    onPositiveClick: async () => {
      const { data, statusCode } = await api.runRconTask(task.uuid);
      if (statusCode.value === 200) {
        message.success(t("rconManager.executed"));
        await loadData();
      } else {
        message.error(data.value?.error || t("rconManager.executeFailed"));
      }
    },
  });
};
const removeTask = (task) => {
  dialog.warning({
    title: t("message.warn"),
    content: t("rconManager.removeTaskConfirm", { name: task.name }),
    positiveText: t("button.confirm"),
    negativeText: t("button.cancel"),
    onPositiveClick: async () => {
      const { data, statusCode } = await api.removeRconTask(task.uuid);
      if (statusCode.value === 200) {
        message.success(t("rconManager.taskRemoved"));
        await loadData();
      } else {
        message.error(data.value?.error || t("rconManager.removeFailed"));
      }
    },
  });
};

const statusType = (status) => {
  if (status === "success") return "success";
  if (status === "failed") return "error";
  return "default";
};
const formatTime = (value) => (value ? dayjs(value).format("YYYY-MM-DD HH:mm:ss") : "—");
</script>

<template>
  <n-drawer v-model:show="visible" :width="drawerWidth" placement="right">
    <n-drawer-content :title="$t('modal.rcon')" closable>
      <n-spin :show="loading">
        <n-card size="small" :title="$t('rconManager.materialTitle')" class="mb-4">
          <n-grid cols="1 700:3" :x-gap="12" :y-gap="12">
            <n-gi>
              <n-select
                v-model:value="selectedPlayerUid"
                filterable
                clearable
                :options="playerOptions"
                :placeholder="$t('input.selectPlayer')"
                aria-label="RCON player"
              />
            </n-gi>
            <n-gi>
              <n-select
                v-model:value="selectedItem"
                filterable
                clearable
                :options="itemOptions"
                :placeholder="$t('input.selectItem')"
                aria-label="RCON item"
              />
            </n-gi>
            <n-gi>
              <n-select
                v-model:value="selectedPal"
                filterable
                clearable
                :options="palOptions"
                :placeholder="$t('input.selectPal')"
                aria-label="RCON pal"
              />
            </n-gi>
          </n-grid>
          <n-flex class="mt-3" align="center">
            <n-text depth="3">{{ $t("rconManager.placeholderHelp") }}</n-text>
            <n-tooltip v-for="placeholder in placeholderLabels" :key="placeholder.key">
              <template #trigger>
                <n-tag size="small" round>{{ placeholder.token }}</n-tag>
              </template>
              {{ placeholder.label }}
            </n-tooltip>
          </n-flex>
        </n-card>

        <n-tabs v-model:value="activeTab" type="segment" animated>
          <n-tab-pane name="commands" :tab="$t('rconManager.commandTab')">
            <n-flex justify="space-between" align="center" class="mb-3">
              <n-text depth="3">{{ $t("rconManager.commandDesc") }}</n-text>
              <n-button type="primary" @click="openCommandModal()">
                {{ $t("button.addRcon") }}
              </n-button>
            </n-flex>
            <n-empty v-if="commands.length === 0" :description="$t('rconManager.noCommands')" />
            <n-space v-else vertical size="medium">
              <n-card v-for="command in commands" :key="command.uuid" size="small">
                <template #header>
                  <n-flex align="center">
                    <n-text strong>{{ command.remark || command.command }}</n-text>
                    <n-code :code="command.command" language="shell" inline />
                  </n-flex>
                </template>
                <template #header-extra>
                  <n-button text type="primary" @click="openCommandModal(command)">
                    {{ $t("rconManager.edit") }}
                  </n-button>
                </template>
                <n-text depth="3" class="block mb-2">
                  {{ $t("rconManager.argumentTemplate") }}:
                  <n-code :code="command.placeholder || '—'" inline />
                </n-text>
                <n-input
                  v-model:value="commandContent[command.uuid]"
                  type="textarea"
                  autosize
                  :placeholder="$t('rconManager.commandContent')"
                  :aria-label="`${command.remark || command.command} arguments`"
                />
                <n-flex class="mt-3" justify="end">
                  <n-button @click="fillCommand(command)">{{ $t("button.fill") }}</n-button>
                  <n-button @click="openTaskModal(command)">
                    {{ $t("rconManager.addTask") }}
                  </n-button>
                  <n-button type="primary" @click="executeCommand(command)">
                    {{ $t("button.execute") }}
                  </n-button>
                  <n-button type="error" secondary @click="removeCommand(command)">
                    {{ $t("button.remove") }}
                  </n-button>
                </n-flex>
              </n-card>
            </n-space>
          </n-tab-pane>

          <n-tab-pane name="tasks" :tab="$t('rconManager.taskTab')">
            <n-flex justify="space-between" align="center" class="mb-3">
              <n-text depth="3">{{ $t("rconManager.taskDesc") }}</n-text>
              <n-button type="primary" :disabled="commands.length === 0" @click="openTaskModal()">
                {{ $t("rconManager.addTask") }}
              </n-button>
            </n-flex>
            <n-empty v-if="tasks.length === 0" :description="$t('rconManager.noTasks')" />
            <n-space v-else vertical size="medium">
              <n-card v-for="task in tasks" :key="task.uuid" size="small">
                <template #header>
                  <n-flex align="center">
                    <n-text strong>{{ task.name }}</n-text>
                    <n-tag size="small" :type="statusType(task.last_status)">
                      {{ $t(`rconManager.status.${task.last_status || 'never'}`) }}
                    </n-tag>
                  </n-flex>
                </template>
                <template #header-extra>
                  <n-switch
                    :value="task.enabled"
                    :aria-label="`${task.name} enabled`"
                    @update:value="(value) => toggleTask(task, value)"
                  />
                </template>
                <n-descriptions label-placement="left" :column="1" size="small">
                  <n-descriptions-item :label="$t('rconManager.boundCommand')">
                    {{ task.rcon_remark || task.rcon_uuid }}
                  </n-descriptions-item>
                  <n-descriptions-item :label="$t('rconManager.argumentContent')">
                    <n-code :code="task.content || '—'" inline />
                  </n-descriptions-item>
                  <n-descriptions-item :label="$t('rconManager.scheduleLabel')">
                    <n-code :code="task.cron" inline />
                  </n-descriptions-item>
                  <n-descriptions-item :label="$t('rconManager.nextRun')">
                    {{ task.enabled ? formatTime(task.next_run_at) : $t("rconManager.paused") }}
                  </n-descriptions-item>
                  <n-descriptions-item :label="$t('rconManager.lastRun')">
                    {{ formatTime(task.last_run_at) }} · {{ $t("rconManager.runCount", { count: task.run_count }) }}
                  </n-descriptions-item>
                  <n-descriptions-item v-if="task.last_result" :label="$t('rconManager.lastResult')">
                    {{ task.last_result }}
                  </n-descriptions-item>
                  <n-descriptions-item v-if="task.last_error" :label="$t('rconManager.lastError')">
                    <n-text type="error">{{ task.last_error }}</n-text>
                  </n-descriptions-item>
                </n-descriptions>
                <n-flex class="mt-3" justify="end">
                  <n-button @click="openTaskModal(null, task)">{{ $t("rconManager.edit") }}</n-button>
                  <n-button type="primary" secondary @click="runTask(task)">
                    {{ $t("rconManager.runNow") }}
                  </n-button>
                  <n-button type="error" secondary @click="removeTask(task)">
                    {{ $t("button.remove") }}
                  </n-button>
                </n-flex>
              </n-card>
            </n-space>
          </n-tab-pane>
        </n-tabs>
      </n-spin>
    </n-drawer-content>
  </n-drawer>

  <n-modal
    v-model:show="commandModal"
    preset="card"
    :title="editingCommandUUID ? $t('rconManager.editCommand') : $t('button.addRcon')"
    style="width: min(92vw, 620px)"
  >
    <n-form label-placement="top">
      <n-form-item :label="$t('input.remark')" required>
        <n-input v-model:value="commandForm.remark" />
      </n-form-item>
      <n-form-item :label="$t('input.rcon')" required>
        <n-input v-model:value="commandForm.command" />
      </n-form-item>
      <n-form-item :label="$t('rconManager.argumentTemplate')">
        <n-input
          v-model:value="commandForm.placeholder"
          type="textarea"
          :placeholder="$t('rconManager.templateExample')"
        />
      </n-form-item>
    </n-form>
    <n-alert type="info" class="mb-3">{{ $t("rconManager.templateTip") }}</n-alert>
    <n-flex justify="space-between">
      <n-upload
        v-if="!editingCommandUUID"
        action="/api/rcon/import"
        name="file"
        accept=".txt"
        :headers="uploadHeaders"
        :show-file-list="false"
        @finish="importFinished"
        @error="importFailed"
      >
        <n-button secondary>{{ $t("button.import") }}</n-button>
      </n-upload>
      <span v-else />
      <n-flex>
        <n-button @click="commandModal = false">{{ $t("button.cancel") }}</n-button>
        <n-button type="primary" @click="saveCommand">{{ $t("button.save") }}</n-button>
      </n-flex>
    </n-flex>
  </n-modal>

  <n-modal
    v-model:show="taskModal"
    preset="card"
    :title="editingTaskUUID ? $t('rconManager.editTask') : $t('rconManager.addTask')"
    style="width: min(92vw, 620px)"
  >
    <n-form label-placement="top">
      <n-form-item :label="$t('rconManager.taskName')" required>
        <n-input v-model:value="taskForm.name" />
      </n-form-item>
      <n-form-item :label="$t('rconManager.boundCommand')" required>
        <n-select
          v-model:value="taskForm.rcon_uuid"
          :options="commandOptions"
          @update:value="selectTaskCommand"
        />
      </n-form-item>
      <n-form-item :label="$t('rconManager.argumentContent')">
        <n-input
          v-model:value="taskForm.content"
          type="textarea"
          autosize
          :placeholder="$t('rconManager.taskContentTip')"
        />
      </n-form-item>
      <n-form-item :label="$t('rconManager.scheduleLabel')" required>
        <n-select v-model:value="scheduleType" :options="scheduleTypeOptions" />
      </n-form-item>
      <n-form-item v-if="scheduleType === 'interval'" :label="$t('rconManager.intervalMinutes')">
        <n-input-number v-model:value="intervalMinutes" :min="1" :max="1440" class="w-full" />
      </n-form-item>
      <n-form-item v-else-if="scheduleType === 'daily'" :label="$t('rconManager.dailyTime')">
        <n-input v-model:value="dailyTime" type="time" />
      </n-form-item>
      <template v-else-if="scheduleType === 'weekly'">
        <n-form-item :label="$t('rconManager.weekdayLabel')">
          <n-select v-model:value="weeklyDay" :options="weekOptions" />
        </n-form-item>
        <n-form-item :label="$t('rconManager.dailyTime')">
          <n-input v-model:value="weeklyTime" type="time" />
        </n-form-item>
      </template>
      <n-form-item v-else :label="$t('rconManager.cronExpression')">
        <n-input v-model:value="customCron" placeholder="0 4 * * *" />
      </n-form-item>
      <n-form-item :label="$t('rconManager.enabled')">
        <n-switch v-model:value="taskForm.enabled" />
      </n-form-item>
    </n-form>
    <n-flex justify="end">
      <n-button @click="taskModal = false">{{ $t("button.cancel") }}</n-button>
      <n-button type="primary" @click="saveTask">{{ $t("button.save") }}</n-button>
    </n-flex>
  </n-modal>
</template>
