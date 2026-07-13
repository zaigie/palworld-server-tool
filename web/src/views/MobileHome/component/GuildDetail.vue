<script setup>
import GuildOverview from "@/components/GuildOverview.vue";

defineProps({
  guildInfo: { type: Object, default: () => ({}) },
});

const emit = defineEmits(["view-player"]);
</script>

<template>
  <div class="guild-detail">
    <guild-overview
      :guild-info="guildInfo"
      compact
      @view-player="emit('view-player', $event)"
    />
  </div>
</template>
