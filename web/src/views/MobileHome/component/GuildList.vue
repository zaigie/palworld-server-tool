<script setup>
import { GroupWorkRound } from "@vicons/material";
import { computed } from "vue";

const props = defineProps(["guildList"]);
const guildList = computed(() => props.guildList);

const emits = defineEmits(["onGetInfo"]);

const onGetInfo = (uid) => {
  emits("onGetInfo", uid);
};
</script>
<template>
  <div class="guild-list">
    <n-list hoverable clickable>
      <n-list-item
        v-for="guild in guildList"
        :key="guild.admin_player_uid"
        @click="onGetInfo(guild.admin_player_uid)"
      >
        <template #prefix>
          <n-avatar
            :style="{ color: 'white', backgroundColor: 'darkorange' }"
            round
          >
            <n-icon>
              <GroupWorkRound />
            </n-icon>
          </n-avatar>
        </template>
        <n-tag type="primary" size="small" round>
          Lv.{{ guild.base_camp_level }}
        </n-tag>
        <span class="pl-2 font-bold">{{ guild.name }}</span>
      </n-list-item>
    </n-list>
  </div>
</template>
