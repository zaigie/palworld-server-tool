<script setup>
import { zhCN, dateZhCN } from "naive-ui";
import pageStore from "@/stores/model/page.js";
import i18n from "@/assets/i18n.js";
import { onMounted } from "vue";

const themeOverrides = {
  common: {
    primaryColor: "#4098fc",
    primaryColorHover: "#4098fc",
  },
};

const locale = ref(null);
const uiLocale = ref(null);
const uiDateLocale = ref(null);

// 移动端适配
// 监听窗口宽度变化
let getScreenWidth = function () {
  let scrollWidth = document.documentElement.clientWidth || window.innerWidth;
  pageStore().setScreenWidth(scrollWidth);
};

onMounted(() => {
  getScreenWidth();
  window.onresize = function () {
    getScreenWidth();
  };

  let localLocale = localStorage.getItem("locale");
  if (localLocale) {
    locale.value = localLocale;
    if (locale.value == "zh") {
      uiLocale.value = zhCN;
      uiDateLocale.value = dateZhCN;
    } else if (locale.value == "en") {
      uiLocale.value = null;
      uiDateLocale.value = null;
    }
  } else {
    localStorage.setItem("locale", "zh");
    locale.value = "zh";
    uiLocale.value = zhCN;
    uiDateLocale.value = dateZhCN;
  }
});
</script>

<template>
  <n-config-provider
    :locale="uiLocale"
    :date-locale="uiDateLocale"
    :theme-overrides="themeOverrides"
  >
    <n-dialog-provider>
      <n-notification-provider>
        <n-message-provider>
          <router-view />
        </n-message-provider>
      </n-notification-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>
