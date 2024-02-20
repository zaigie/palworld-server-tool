<script setup>
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import skillMap from "@/assets/skill.json";

const { t, locale } = useI18n();

const props = defineProps(["palDetail"]);
const palDetail = computed(() => props.palDetail);

const getPalAvatar = (name) => {
  return new URL(`../../../assets/pal/${name}.png`, import.meta.url).href;
};
const getUnknowPalAvatar = (is_boss = false) => {
  if (is_boss) {
    return new URL("@/assets/pal/BOSS_Unknown.png", import.meta.url).href;
  }
  return new URL("@/assets/pal/Unknown.png", import.meta.url).href;
};
const displayHP = (hp, max_hp) => {
  return (hp / 1000).toFixed(0) + "/" + (max_hp / 1000).toFixed(0);
};
const percentageHP = (hp, max_hp) => {
  if (max_hp === 0) {
    return 0;
  }
  return ((hp / max_hp) * 100).toFixed(2);
};
</script>

<template>
  <div class="pal-detail">
    <n-space class="mb-2" justify="center">
      <n-avatar
        :size="64"
        :src="getPalAvatar(palDetail.type)"
        :fallback-src="getUnknowPalAvatar(palDetail.is_boss)"
      ></n-avatar>
    </n-space>
    <n-space class="mb-2" justify="center">
      <n-tag v-if="palDetail.is_boss" type="success" round>Boss</n-tag>
      <n-tag v-else-if="palDetail.is_lucky" type="warning" round>{{
        $t("pal.lucky")
      }}</n-tag>
      <n-tag v-else-if="palDetail.is_tower" type="error" round>{{
        $t("pal.tower")
      }}</n-tag>
    </n-space>
    <n-space vertical>
      <n-progress
        type="line"
        status="error"
        indicator-placement="inside"
        :percentage="percentageHP(palDetail.hp, palDetail.max_hp)"
        :height="24"
        :border-radius="4"
        :fill-border-radius="0"
        >HP: {{ displayHP(palDetail.hp, palDetail.max_hp) }}</n-progress
      >
      <n-grid cols="4">
        <n-gi>
          <n-statistic :label="$t('pal.ranged')" :value="palDetail.ranged" />
        </n-gi>
        <n-gi>
          <n-statistic :label="$t('pal.defense')" :value="palDetail.defense" />
        </n-gi>
        <n-gi>
          <n-statistic :label="$t('pal.melee')" :value="palDetail.melee" />
        </n-gi>
        <n-gi>
          <n-statistic :label="$t('pal.rank')" :value="palDetail.rank" />
        </n-gi>
      </n-grid>
    </n-space>
    <n-space vertical>
      <div v-for="skill in palDetail.skills" :key="skill">
        <n-tag type="warning">{{
          skillMap[locale][skill] ? skillMap[locale][skill].name : skill
        }}</n-tag>
        :
        {{ skillMap[locale][skill] ? skillMap[locale][skill].desc : "-" }}
      </div>
    </n-space>
  </div>
</template>
