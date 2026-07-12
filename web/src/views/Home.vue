<script setup>
import PcHome from "@/views/PcHome/PcHome.vue";
import MobileHome from "@/views/MobileHome/MobileHome.vue";
import pageStore from "@/stores/model/page";
import { onMounted, ref } from "vue";
import userStore from "@/stores/model/user";
import ApiService from "@/service/api";
import FirstRunSetup from "@/components/FirstRunSetup.vue";
import ConfigManager from "@/components/ConfigManager.vue";

const PALWORLD_TOKEN = "palworld_token";
const pageWidth = computed(() => pageStore().getScreenWidth());
const configReady = ref(false);
const initialized = ref(false);
const showConfig = ref(false);

onMounted(async () => {
  let token = localStorage.getItem(PALWORLD_TOKEN);
  if (token) userStore().setIsLogin(true, token);
  const { data, statusCode } = await new ApiService().getConfigStatus();
  initialized.value =
    statusCode.value === 200 && data.value?.initialized === true;
  configReady.value = true;
});

const handleInitialized = () => {
  initialized.value = true;
  showConfig.value = true;
};
</script>

<template>
  <div v-if="configReady">
    <template v-if="initialized">
      <pc-home v-if="pageWidth >= 768" @open-config="showConfig = true" />
      <mobile-home v-else @open-config="showConfig = true" />
    </template>
    <first-run-setup :show="!initialized" @initialized="handleInitialized" />
    <config-manager v-model:show="showConfig" />
  </div>
  <div v-else class="h-screen flex items-center justify-center">
    <n-spin size="large" />
  </div>
</template>
