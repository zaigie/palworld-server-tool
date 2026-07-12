<script setup>
import { nextTick, ref, watch } from "vue";
import { useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import userStore from "@/stores/model/user";
import DirectoryPicker from "@/components/DirectoryPicker.vue";

const props = defineProps({ show: Boolean });
const emit = defineEmits(["update:show"]);
const { t } = useI18n();
const message = useMessage();
const loading = ref(false);
const saving = ref(false);
const showDirectoryPicker = ref(false);
const newPassword = ref("");
const passwordConfirmation = ref("");
const saveTesting = ref(false);
const rconTesting = ref(false);
const saveStatus = ref({ status: "unconfigured", message: "" });
const rconStatus = ref({ status: "unconfigured", message: "" });
const sourcePaths = ref({ directory: "", agent: "" });

const statusType = (status) =>
  ({ normal: "success", error: "error", unconfigured: "warning" })[status] ||
  "default";
const statusLabel = (status) =>
  t(`configuration.connectionStatus.${status || "unconfigured"}`);
const restartFieldLabel = (field) =>
  ({
    "web.port": t("configuration.webPort"),
    "web.tls": "TLS",
    "web.cert_path": t("configuration.certPath"),
    "web.key_path": t("configuration.keyPath"),
    "web.public_url": t("configuration.publicUrl"),
    "task.sync_interval": t("configuration.playerSyncInterval"),
    "save.sync_interval": t("configuration.saveSyncInterval"),
    "save.backup_interval": t("configuration.backupInterval"),
  })[field] || field;

const emptySettings = () => ({
  web: { port: 8080, tls: false, cert_path: "", key_path: "", public_url: "" },
  task: {
    sync_interval: 60,
    player_logging: false,
    player_login_message: "",
    player_logout_message: "",
  },
  rcon: { address: "", password: "", use_base64: false, timeout: 5 },
  rest: { address: "", username: "admin", password: "", timeout: 5 },
  save: {
    source_mode: "directory",
    path: "",
    decode_path: "",
    sync_interval: 120,
    backup_interval: 14400,
    backup_keep_days: 7,
  },
  manage: { kick_non_whitelist: false },
});
const settings = ref(emptySettings());

const checkSaveSource = async () => {
  saveTesting.value = true;
  const { data, statusCode } = await new ApiService().testSaveConfig(
    settings.value.save,
  );
  saveTesting.value = false;
  if (statusCode.value !== 200) {
    saveStatus.value = {
      status: "error",
      message: data.value?.error || t("configuration.connectionTestFailed"),
    };
    return;
  }
  saveStatus.value = data.value;
};

const testRcon = async () => {
  rconTesting.value = true;
  const { data, statusCode } = await new ApiService().testRconConfig(
    settings.value.rcon,
  );
  rconTesting.value = false;
  if (statusCode.value !== 200) {
    rconStatus.value = {
      status: "error",
      message: data.value?.error || t("configuration.connectionTestFailed"),
    };
    return;
  }
  rconStatus.value = data.value;
};

const markRconDirty = () => {
  rconStatus.value = {
    status: "unconfigured",
    message: t("configuration.retestRequired"),
  };
};

const changeSourceMode = async (mode) => {
  const previousMode = settings.value.save.source_mode;
  sourcePaths.value[previousMode] = settings.value.save.path;
  settings.value.save.source_mode = mode;
  settings.value.save.path = sourcePaths.value[mode] || "";
  await nextTick();
  checkSaveSource();
};

const load = async () => {
  loading.value = true;
  const { data, statusCode } = await new ApiService().getConfig();
  loading.value = false;
  if (statusCode.value !== 200) {
    message.error(data.value?.error || t("configuration.loadFailed"));
    emit("update:show", false);
    return;
  }
  settings.value = data.value;
  sourcePaths.value[settings.value.save.source_mode] = settings.value.save.path;
  newPassword.value = "";
  passwordConfirmation.value = "";
  await Promise.all([checkSaveSource(), testRcon()]);
};

const save = async () => {
  if (newPassword.value !== passwordConfirmation.value) {
    message.error(t("configuration.passwordMismatch"));
    return;
  }
  saving.value = true;
  const { data, statusCode } = await new ApiService().updateConfig({
    settings: settings.value,
    new_password: newPassword.value,
  });
  saving.value = false;
  if (statusCode.value !== 200) {
    message.error(data.value?.error || t("configuration.saveFailed"));
    return;
  }
  if (data.value?.token) {
    localStorage.setItem("palworld_token", data.value.token);
    userStore().setIsLogin(true, data.value.token);
  }
  if (data.value?.restart_required) {
    const fields = (data.value.restart_fields || [])
      .map(restartFieldLabel)
      .join("、");
    message.warning(t("configuration.savedWithRestart", { fields }));
  } else {
    message.success(t("configuration.savedImmediately"));
  }
  emit("update:show", false);
};

const selectDirectory = (path) => {
  settings.value.save.path = path;
  sourcePaths.value.directory = path;
  showDirectoryPicker.value = false;
  checkSaveSource();
};

watch(
  () => props.show,
  (show) => {
    if (show) load();
  },
);
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    class="config-card"
    style="width: min(94vw, 860px); max-height: calc(100vh - 64px)"
    content-style="overflow: hidden"
    :title="$t('configuration.title')"
    @update:show="emit('update:show', $event)"
  >
    <n-spin :show="loading">
      <n-alert type="warning" :bordered="false" class="mb-3">
        {{ $t("configuration.migrationWarning") }}
      </n-alert>
      <n-alert type="info" :bordered="false" class="mb-4">
        {{ $t("configuration.restartWarning") }}
      </n-alert>

      <n-scrollbar
        class="config-scroll"
        style="max-height: min(52vh, 560px); padding-right: 10px"
      >
        <n-collapse :default-expanded-names="['save', 'rcon', 'rest']">
          <n-collapse-item :title="$t('configuration.saveSection')" name="save">
            <template #header-extra>
              <n-spin v-if="saveTesting" size="small" />
              <n-tooltip v-else>
                <template #trigger>
                  <n-tag
                    size="small"
                    round
                    :type="statusType(saveStatus.status)"
                  >
                    {{ statusLabel(saveStatus.status) }}
                  </n-tag>
                </template>
                {{ saveStatus.message || $t("configuration.noStatusDetails") }}
              </n-tooltip>
            </template>
            <n-form label-placement="top">
              <n-form-item :label="$t('configuration.sourceMode')">
                <n-radio-group
                  :value="settings.save.source_mode"
                  @update:value="changeSourceMode"
                >
                  <n-space>
                    <n-radio value="directory">{{
                      $t("configuration.directoryMode")
                    }}</n-radio>
                    <n-radio value="agent">{{
                      $t("configuration.agentMode")
                    }}</n-radio>
                  </n-space>
                </n-radio-group>
              </n-form-item>
              <n-form-item
                :label="
                  settings.save.source_mode === 'agent'
                    ? $t('configuration.agentUrl')
                    : $t('configuration.saveDirectory')
                "
              >
                <n-input-group>
                  <n-input
                    v-model:value="settings.save.path"
                    @blur="checkSaveSource"
                    :placeholder="
                      settings.save.source_mode === 'agent'
                        ? 'http://game-server:8081/sync'
                        : $t('configuration.saveDirectoryPlaceholder')
                    "
                  />
                  <n-button
                    v-if="settings.save.source_mode === 'directory'"
                    @click="showDirectoryPicker = true"
                  >
                    {{ $t("configuration.browse") }}
                  </n-button>
                </n-input-group>
              </n-form-item>
              <div class="form-grid">
                <n-form-item :label="$t('configuration.decodePath')">
                  <n-input
                    v-model:value="settings.save.decode_path"
                    :placeholder="$t('configuration.autoDetect')"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.saveSyncInterval')">
                  <n-input-number
                    v-model:value="settings.save.sync_interval"
                    :min="0"
                    class="full-width"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.backupInterval')">
                  <n-input-number
                    v-model:value="settings.save.backup_interval"
                    :min="0"
                    class="full-width"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.backupKeepDays')">
                  <n-input-number
                    v-model:value="settings.save.backup_keep_days"
                    :min="0"
                    class="full-width"
                  />
                </n-form-item>
              </div>
            </n-form>
          </n-collapse-item>

          <n-collapse-item title="RCON" name="rcon">
            <template #header-extra>
              <n-space size="small" align="center" @click.stop>
                <n-tooltip>
                  <template #trigger>
                    <n-tag
                      size="small"
                      round
                      :type="statusType(rconStatus.status)"
                    >
                      {{ statusLabel(rconStatus.status) }}
                    </n-tag>
                  </template>
                  {{
                    rconStatus.message || $t("configuration.noStatusDetails")
                  }}
                </n-tooltip>
                <n-button
                  size="tiny"
                  secondary
                  :loading="rconTesting"
                  @click.stop="testRcon"
                >
                  {{ $t("configuration.testConnection") }}
                </n-button>
              </n-space>
            </template>
            <n-form label-placement="top">
              <div class="form-grid">
                <n-form-item :label="$t('configuration.address')">
                  <n-input
                    v-model:value="settings.rcon.address"
                    placeholder="127.0.0.1:25575"
                    @update:value="markRconDirty"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.password')">
                  <n-input
                    v-model:value="settings.rcon.password"
                    type="password"
                    show-password-on="click"
                    @update:value="markRconDirty"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.timeout')">
                  <n-input-number
                    v-model:value="settings.rcon.timeout"
                    :min="0"
                    class="full-width"
                    @update:value="markRconDirty"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.base64')">
                  <n-switch
                    v-model:value="settings.rcon.use_base64"
                    @update:value="markRconDirty"
                  />
                </n-form-item>
              </div>
            </n-form>
          </n-collapse-item>

          <n-collapse-item title="REST API" name="rest">
            <n-form label-placement="top">
              <div class="form-grid">
                <n-form-item :label="$t('configuration.address')">
                  <n-input
                    v-model:value="settings.rest.address"
                    placeholder="http://127.0.0.1:8212"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.username')">
                  <n-input v-model:value="settings.rest.username" />
                </n-form-item>
                <n-form-item :label="$t('configuration.password')">
                  <n-input
                    v-model:value="settings.rest.password"
                    type="password"
                    show-password-on="click"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.timeout')">
                  <n-input-number
                    v-model:value="settings.rest.timeout"
                    :min="0"
                    class="full-width"
                  />
                </n-form-item>
              </div>
            </n-form>
          </n-collapse-item>

          <n-collapse-item
            :title="$t('configuration.tasksSection')"
            name="tasks"
          >
            <n-form label-placement="top">
              <div class="form-grid">
                <n-form-item :label="$t('configuration.playerSyncInterval')">
                  <n-input-number
                    v-model:value="settings.task.sync_interval"
                    :min="0"
                    class="full-width"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.playerLogging')">
                  <n-switch v-model:value="settings.task.player_logging" />
                </n-form-item>
              </div>
              <n-form-item :label="$t('configuration.loginMessage')">
                <n-input
                  v-model:value="settings.task.player_login_message"
                  type="textarea"
                />
              </n-form-item>
              <n-form-item :label="$t('configuration.logoutMessage')">
                <n-input
                  v-model:value="settings.task.player_logout_message"
                  type="textarea"
                />
              </n-form-item>
              <n-checkbox v-model:checked="settings.manage.kick_non_whitelist">
                {{ $t("configuration.kickNonWhitelist") }}
              </n-checkbox>
            </n-form>
          </n-collapse-item>

          <n-collapse-item :title="$t('configuration.webSection')" name="web">
            <n-form label-placement="top">
              <div class="form-grid">
                <n-form-item :label="$t('configuration.webPort')">
                  <n-input-number
                    v-model:value="settings.web.port"
                    :min="1"
                    :max="65535"
                    class="full-width"
                  />
                </n-form-item>
                <n-form-item label="TLS">
                  <n-switch v-model:value="settings.web.tls" />
                </n-form-item>
                <n-form-item :label="$t('configuration.certPath')">
                  <n-input v-model:value="settings.web.cert_path" />
                </n-form-item>
                <n-form-item :label="$t('configuration.keyPath')">
                  <n-input v-model:value="settings.web.key_path" />
                </n-form-item>
              </div>
              <n-form-item :label="$t('configuration.publicUrl')">
                <n-input
                  v-model:value="settings.web.public_url"
                  placeholder="https://pst.example.com"
                />
              </n-form-item>
            </n-form>
          </n-collapse-item>

          <n-collapse-item
            :title="$t('configuration.securitySection')"
            name="security"
          >
            <n-alert type="info" :bordered="false" class="mb-3">
              {{ $t("configuration.panelPasswordNotice") }}
            </n-alert>
            <n-form label-placement="top">
              <div class="form-grid">
                <n-form-item :label="$t('configuration.newAdminPassword')">
                  <n-input
                    v-model:value="newPassword"
                    type="password"
                    show-password-on="click"
                    autocomplete="new-password"
                  />
                </n-form-item>
                <n-form-item :label="$t('configuration.confirmPassword')">
                  <n-input
                    v-model:value="passwordConfirmation"
                    type="password"
                    show-password-on="click"
                    autocomplete="new-password"
                  />
                </n-form-item>
              </div>
            </n-form>
          </n-collapse-item>
        </n-collapse>
      </n-scrollbar>
    </n-spin>

    <template #footer>
      <n-space justify="end">
        <n-button @click="emit('update:show', false)">{{
          $t("button.cancel")
        }}</n-button>
        <n-button type="primary" :loading="saving" @click="save">{{
          $t("button.save")
        }}</n-button>
      </n-space>
    </template>
  </n-modal>

  <directory-picker
    v-model:show="showDirectoryPicker"
    :initial-path="
      settings.save.source_mode === 'directory' ? settings.save.path : ''
    "
    @select="selectDirectory"
  />
</template>

<style scoped>
.config-scroll {
  min-height: 220px;
}
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 16px;
}
.full-width {
  width: 100%;
}
@media (max-width: 700px) {
  .config-scroll {
    min-height: 160px;
  }
  .form-grid {
    grid-template-columns: 1fr;
    gap: 0;
  }
}
</style>
