<script setup>
import PcHome from "@/views/PcHome/PcHome.vue";
import MobileHome from "@/views/MobileHome/MobileHome.vue";
import pageStore from "@/stores/model/page";
import { onMounted } from "vue";
import userStore from "@/stores/model/user";

const PALWORLD_TOKEN = "palworld_token";
const pageWidth = computed(() => pageStore().getScreenWidth());

onMounted(() => {
  let token = localStorage.getItem(PALWORLD_TOKEN);
  if (token) userStore().setIsLogin(true, token);
});
</script>

<template>
  <div>
    <pc-home v-if="pageWidth >= 768"></pc-home>
    <mobile-home v-else></mobile-home>
  </div>
</template>
