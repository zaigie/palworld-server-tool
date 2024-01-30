<script setup>
import { zhCN, dateZhCN } from 'naive-ui'
import pageStore from '@/stores/model/page.js'
import { onMounted } from 'vue'

const themeOverrides = {
  common: {
    primaryColor: '#4098fc',
    primaryColorHover: '#4098fc'
  }
}

// 移动端适配
// 监听窗口宽度变化
let getScreenWidth = function () {
  let scrollWidth = document.documentElement.clientWidth || window.innerWidth
  pageStore().setScreenWidth(scrollWidth)
}

onMounted(() => {
  getScreenWidth()
  window.onresize = function () {
    getScreenWidth()
  }
})
</script>

<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-notification-provider>
        <n-message-provider>
          <router-view />
        </n-message-provider>
      </n-notification-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>
