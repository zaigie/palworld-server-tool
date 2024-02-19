<script setup>
import PcHome from "@/views/PcHome/PcHome.vue";
import MobileHome from "@/views/MobileHome/MobileHome.vue";
import pageStore from "@/stores/model/page";

const pageWidth = computed(() => pageStore().getScreenWidth());
</script>

<template>
  <div>
    <pc-home v-if="pageWidth >= 768"></pc-home>
    <mobile-home v-else></mobile-home>
  </div>
</template>
