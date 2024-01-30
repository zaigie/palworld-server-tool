<script setup>
import router from '@/router'
import { ChevronLeft12Regular, ArrowRight12Filled } from '@vicons/fluent'

const login = () => {
  console.log('登录')
}

const navigateToReturn = () => {
  router.go(-1)
}
</script>
<template>
  <div class="login-page w-screen h-screen flex justify-center items-center">
    <div class="w-90% max-w-100 px-4 pt-6 rounded-lg bg-#fff drop-shadow-md">
      <span class="block pb-6 text-6 text-center font-bold text-#666">PalWorld 服务器管理后台</span>
      <n-input class="mb-6" type="password" show-password-on="click" placeholder="密码"></n-input>
      <div class="flex flex-wrap justify-center">
        <n-button strong circle type="info" size="large" @click="login">
          <template #icon>
            <n-icon :size="26"><ArrowRight12Filled /></n-icon>
          </template>
        </n-button>
        <span class="block w-full text-center text-xs text-#999 pt-1">登录进入</span>
      </div>
      <div class="flex border-t border-t-dashed border-t-#ddd mt-3">
        <span
          class="pt-2 pb-3 text-#666 flex justify-left items-center cursor-pointer"
          @click="navigateToReturn"
        >
          <n-icon>
            <ChevronLeft12Regular />
          </n-icon>
          <span>返回首页</span>
        </span>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.login-page {
  background-image: linear-gradient(to right bottom, #fafeff, #edf9fd);
}
</style>
