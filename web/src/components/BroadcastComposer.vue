<script setup>
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useMessage } from "naive-ui";
import ApiService from "@/service/api";

const props = defineProps({ show: Boolean });
const emit = defineEmits(["update:show"]);
const { t } = useI18n();
const message = useMessage();
const api = new ApiService();
const content = ref("");
const sending = ref(false);
const templates = computed(() => [
  t("broadcast.templateRestart"),
  t("broadcast.templateMaintenance"),
  t("broadcast.templateWelcome"),
]);

const close = () => emit("update:show", false);
const send = async () => {
  if (!content.value.trim()) return;
  sending.value = true;
  try {
    const { data, statusCode } = await api.sendBroadcast({
      message: content.value.trim(),
    });
    if (statusCode.value === 200) {
      message.success(t("message.broadcastsuccess"));
      content.value = "";
      close();
    } else {
      message.error(t("message.broadcastfail", { err: data.value?.error }));
    }
  } finally {
    sending.value = false;
  }
};

watch(
  () => props.show,
  (show) => {
    if (!show) content.value = "";
  }
);
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    style="width: 92%; max-width: 620px"
    :title="$t('modal.broadcast')"
    @update:show="emit('update:show', $event)"
  >
    <n-space vertical size="large">
      <div>
        <n-text depth="3">{{ $t("broadcast.templates") }}</n-text>
        <n-flex class="mt-2">
          <n-button
            v-for="template in templates"
            :key="template"
            size="small"
            secondary
            @click="content = template"
          >
            {{ template }}
          </n-button>
        </n-flex>
      </div>
      <n-input
        v-model:value="content"
        type="textarea"
        :autosize="{ minRows: 3, maxRows: 7 }"
        :maxlength="500"
        show-count
        :placeholder="$t('broadcast.placeholder')"
        :aria-label="$t('modal.broadcast')"
      />
      <n-alert type="info" :title="$t('broadcast.preview')">
        {{ content || $t("broadcast.previewEmpty") }}
      </n-alert>
    </n-space>
    <template #footer>
      <n-flex justify="end">
        <n-button @click="close">{{ $t("button.cancel") }}</n-button>
        <n-button
          type="primary"
          :disabled="!content.trim()"
          :loading="sending"
          @click="send"
        >
          {{ $t("button.broadcast") }}
        </n-button>
      </n-flex>
    </template>
  </n-modal>
</template>
