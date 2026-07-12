<script setup>
import { ref } from "vue";
import { useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";
import userStore from "@/stores/model/user";

defineProps({ show: Boolean });
const emit = defineEmits(["initialized"]);
const { t } = useI18n();
const message = useMessage();
const password = ref("");
const confirmation = ref("");
const saving = ref(false);

const initialize = async () => {
  if (!password.value.trim()) {
    message.error(t("configuration.passwordRequired"));
    return;
  }
  if (password.value !== confirmation.value) {
    message.error(t("configuration.passwordMismatch"));
    return;
  }
  saving.value = true;
  const { data, statusCode } = await new ApiService().initializeConfig({
    password: password.value,
  });
  saving.value = false;
  if (statusCode.value !== 200 || !data.value?.token) {
    message.error(data.value?.error || t("configuration.initializeFailed"));
    return;
  }
  localStorage.setItem("palworld_token", data.value.token);
  userStore().setIsLogin(true, data.value.token);
  message.success(t("configuration.initialized"));
  emit("initialized");
};
</script>

<template>
  <n-modal
    :show="show"
    :mask-closable="false"
    :close-on-esc="false"
    preset="card"
    class="first-run-card"
    style="width: min(92vw, 520px)"
    :closable="false"
    :title="$t('configuration.firstRunTitle')"
  >
    <n-alert type="info" :bordered="false" class="mb-4">
      <div>{{ $t("configuration.firstRunDescription") }}</div>
      <div class="mt-2 font-medium">
        {{ $t("configuration.panelPasswordNotice") }}
      </div>
    </n-alert>
    <n-form label-placement="top" @submit.prevent="initialize">
      <n-form-item :label="$t('configuration.adminPassword')">
        <n-input
          v-model:value="password"
          type="password"
          show-password-on="click"
          autocomplete="new-password"
          @keyup.enter="initialize"
        />
      </n-form-item>
      <n-form-item :label="$t('configuration.confirmPassword')">
        <n-input
          v-model:value="confirmation"
          type="password"
          show-password-on="click"
          autocomplete="new-password"
          @keyup.enter="initialize"
        />
      </n-form-item>
    </n-form>
    <template #footer>
      <n-space justify="end">
        <n-button type="primary" :loading="saving" @click="initialize">
          {{ $t("configuration.createAdministrator") }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>
