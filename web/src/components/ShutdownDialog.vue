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
const seconds = ref(60);
const shutdownMessage = ref("");
const confirmation = ref("");
const submitting = ref(false);
const canSubmit = computed(
  () => shutdownMessage.value.trim() && confirmation.value === "SHUTDOWN"
);

const close = () => emit("update:show", false);
const submit = async () => {
  if (!canSubmit.value) return;
  submitting.value = true;
  try {
    const { data, statusCode } = await api.shutdownServer({
      seconds: seconds.value,
      message: shutdownMessage.value.trim(),
    });
    if (statusCode.value === 200) {
      message.success(t("message.shutdownsuccess"));
      close();
    } else {
      message.error(t("message.shutdownfail", { err: data.value?.error }));
    }
  } finally {
    submitting.value = false;
  }
};

watch(
  () => props.show,
  (show) => {
    if (show) {
      seconds.value = 60;
      shutdownMessage.value = t("danger.defaultMessage", { seconds: 60 });
      confirmation.value = "";
    }
  }
);
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    style="width: 92%; max-width: 560px"
    :title="$t('danger.title')"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <n-space vertical size="large">
      <n-alert type="error" :title="$t('danger.warning')">
        {{ $t("danger.description") }}
      </n-alert>
      <n-form label-placement="top">
        <n-form-item :label="$t('danger.countdown')">
          <n-input-number
            v-model:value="seconds"
            :min="10"
            :max="3600"
            class="w-full"
            @update:value="shutdownMessage = $t('danger.defaultMessage', { seconds: $event })"
          />
        </n-form-item>
        <n-form-item :label="$t('danger.broadcastMessage')">
          <n-input v-model:value="shutdownMessage" type="textarea" />
        </n-form-item>
        <n-form-item :label="$t('danger.confirmLabel')">
          <n-input
            v-model:value="confirmation"
            :placeholder="$t('danger.confirmPlaceholder')"
            autocomplete="off"
          />
        </n-form-item>
      </n-form>
    </n-space>
    <template #footer>
      <n-flex justify="end">
        <n-button @click="close">{{ $t("button.cancel") }}</n-button>
        <n-button
          type="error"
          :disabled="!canSubmit"
          :loading="submitting"
          @click="submit"
        >
          {{ $t("button.shutdown") }}
        </n-button>
      </n-flex>
    </template>
  </n-modal>
</template>
