<script setup>
import { ref, watch } from "vue";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import { useMessage, useDialog } from "naive-ui";
import ApiService from "@/service/api";

const props = defineProps({ show: Boolean });
const emit = defineEmits(["update:show"]);
const { t } = useI18n();
const message = useMessage();
const dialog = useDialog();
const api = new ApiService();
const backups = ref([]);
const loading = ref(false);
const busyId = ref("");

const load = async () => {
  loading.value = true;
  try {
    const { data } = await api.getBackupList({});
    backups.value = Array.isArray(data.value) ? data.value : [];
  } finally {
    loading.value = false;
  }
};

const download = async (backup) => {
  busyId.value = backup.backup_id;
  try {
    const { data, execute } = await api.downloadBackup(backup.backup_id);
    await execute();
    const url = URL.createObjectURL(data.value);
    const link = document.createElement("a");
    link.href = url;
    link.download = backup.path;
    link.click();
    URL.revokeObjectURL(url);
    message.success(t("message.downloadsuccess"));
  } finally {
    busyId.value = "";
  }
};

const remove = (backup) => {
  dialog.warning({
    title: t("backup.removeTitle"),
    content: t("backup.removeConfirm"),
    positiveText: t("button.remove"),
    negativeText: t("button.cancel"),
    onPositiveClick: async () => {
      busyId.value = backup.backup_id;
      const { data, statusCode } = await api.removeBackup(backup.backup_id);
      busyId.value = "";
      if (statusCode.value === 200) {
        message.success(t("message.removebackupsuccess"));
        await load();
      } else {
        message.error(t("message.removebackupfail", { err: data.value?.error }));
      }
    },
  });
};

watch(
  () => props.show,
  (show) => show && load(),
  { immediate: true }
);
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    style="width: 94%; max-width: 760px"
    :title="$t('modal.backup')"
    @update:show="emit('update:show', $event)"
  >
    <n-spin :show="loading">
      <n-empty
        v-if="backups.length === 0"
        :description="$t('backup.emptyDescription')"
      >
        <template #extra>
          <n-text depth="3">{{ $t("backup.emptyHint") }}</n-text>
        </template>
      </n-empty>
      <n-list v-else hoverable bordered>
        <n-list-item v-for="backup in backups" :key="backup.backup_id">
          <n-thing :title="dayjs(backup.save_time).format('YYYY-MM-DD HH:mm:ss')">
            <template #description>{{ backup.path }}</template>
          </n-thing>
          <template #suffix>
            <n-flex>
              <n-button
                size="small"
                type="primary"
                secondary
                :loading="busyId === backup.backup_id"
                @click="download(backup)"
              >{{ $t("button.download") }}</n-button>
              <n-button size="small" type="error" secondary @click="remove(backup)">
                {{ $t("button.remove") }}
              </n-button>
            </n-flex>
          </template>
        </n-list-item>
      </n-list>
    </n-spin>
  </n-modal>
</template>
