<script setup>
import { LogInOutlined, FaceRetouchingOffFilled } from '@vicons/material'
import { RefreshSharp } from '@vicons/ionicons5'
import { computed, onMounted, ref } from 'vue'
import ApiService from '@/service/api'
import { useMessage } from 'naive-ui'
import pageStore from '@/stores/model/page.js'

const message = useMessage()

const pageWidth = computed(() => pageStore().getScreenWidth())
const smallScreen = computed(() => pageWidth.value < 1024)
const avatar =
  'data:image/webp;base64,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'

const loading = ref(false)
const serverInfo = ref({}) // 服务器信息
const playerList = ref([]) // 玩家列表

// 获取数据
const getServerInfo = async () => {
  const { data, error } = await new ApiService().getServerInfo()
  console.log(data.value, error)
  serverInfo.value = data.value
}
const getPlayerList = async () => {
  const { data, error } = await new ApiService().getPlayerList({ update: false })
  console.log(data.value, error)
  playerList.value = []
  data.value.forEach((item) => {
    playerList.value.push({
      ...item,
      select: false
    })
  })
}

// 操作
const currentHandleType = ref(null)
const showModal = ref(false)
const password = ref(null)
// 游戏广播
const clickBroadcast = () => {
  currentHandleType.value = 'broadcast'
  showModal.value = true
}
// 踢出
const clickKick = () => {
  currentHandleType.value = 'kick'
  showModal.value = true
}

// 封禁
const clickBan = () => {
  currentHandleType.value = 'ban'
  showModal.value = true
}

// 刷新
const clickRefresh = () => {
  currentHandleType.value = 'refresh'
  showModal.value = true
}
const goToRefresh = async () => {
  loading.value = true
  await getPlayerList()
  loading.value = false
  password.value = null
  message.success('刷新成功')
}

// 最终操作
const handleConfirm = () => {
  if (checkoutPassord()) {
    showModal.value = false
    if (currentHandleType.value === 'refresh') {
      goToRefresh()
    }
  } else {
    message.error('密码错误')
  }
}
const handleCancel = () => {
  showModal.value = false
  password.value = null
}
const checkoutPassord = async () => {
  return true
}

onMounted(async () => {
  loading.value = true
  await getServerInfo()
  await getPlayerList()
  loading.value = false
})
</script>

<template>
  <div class="admin-system p-3">
    <div class="flex justify-between items-center">
      <span class="line-clamp-1" :class="smallScreen ? 'text-xl' : 'text-2xl'"
        >PalWorld 服务器管理</span
      >
      <n-button
        type="success"
        round
        :size="smallScreen ? 'medium' : 'large'"
        @click="clickBroadcast"
        >游戏广播</n-button
      >
    </div>
    <div class="h-6"></div>
    <div class="text-center" :class="smallScreen ? 'text-lg' : 'text-xl'">
      <n-select
        v-model:value="value"
        :options="[]"
        placeholder="更换服务器"
        class="w-90% max-w-80 mx-auto mb-4"
        disabled
      />
      <span class="block">
        服务器名称: <strong>{{ serverInfo.name ?? '- -' }}</strong>
      </span>
      <span class="block">
        版本: <strong>{{ serverInfo.version ?? '- -' }}</strong>
      </span>
    </div>
    <div
      class="w-full flex items-center mt-5"
      :class="smallScreen ? 'justify-between' : 'justify-center'"
    >
      <n-button
        type="error"
        round
        class="w-30% max-w-50"
        :class="smallScreen ? '' : 'mx-2'"
        @click="clickKick"
      >
        <template #icon>
          <n-icon><LogInOutlined /></n-icon>
        </template>
        踢出
      </n-button>
      <n-button
        type="warning"
        round
        class="w-30% max-w-50"
        :class="smallScreen ? '' : 'mx-2'"
        @click="clickBan"
      >
        <template #icon>
          <n-icon><FaceRetouchingOffFilled /></n-icon>
        </template>
        封禁
      </n-button>
      <n-button
        type="info"
        round
        class="w-30% max-w-50"
        :class="smallScreen ? '' : 'mx-2'"
        @click="clickRefresh"
      >
        <template #icon>
          <n-icon><RefreshSharp /></n-icon>
        </template>
        刷新
      </n-button>
    </div>
    <div class="mt-4 rounded-lg bg-#F2F4F5 p-4" v-if="!loading && playerList.length > 0">
      <n-list class="bg-transparent">
        <n-list-item v-for="player in playerList" :key="player.playeruid">
          <template #prefix>
            <n-avatar :src="avatar" round></n-avatar>
          </template>
          <template #suffix>
            <n-checkbox v-model:checked="player.select"></n-checkbox>
          </template>
          <div>
            <div>
              <n-tag
                :bordered="false"
                size="small"
                :type="player.online ? 'success' : 'error'"
                round
              >
                {{ player.online ? '在线' : '离线' }}
              </n-tag>
              <span class="pl-3 font-bold">{{ player.name }}</span>
            </div>
            <span class="inline-block mt-1 rounded-full bg-#ddd text-xs px-2 py-0.5"
              >最近：{{ player.last_online }}</span
            >
          </div>
        </n-list-item>
      </n-list>
    </div>
    <div
      v-if="!loading && playerList.length === 0"
      class="w-full h-25 flex justify-center items-center"
    >
      <n-empty description="什么都没有"> </n-empty>
    </div>
    <div v-if="loading" class="w-full h-25 flex justify-center items-center">
      <n-spin size="small" />
    </div>
  </div>
  <!-- 输入密码模态框 -->
  <n-modal
    v-model:show="showModal"
    class="custom-card"
    preset="card"
    :mask-closable="false"
    style="width: 90%; max-width: 400px; border-radius: 12px"
    title="需要密码"
    size="huge"
    :bordered="false"
    :segmented="segmented"
    header-style="padding:16px;"
    content-style="padding: 16px; padding-top: 0;"
  >
    <div class="pb-1">请输入 AdminPassword 以继续</div>
    <n-input v-model:value="password" type="password" show-password-on="click" placeholder="密码" />
    <div class="mt-4 flex justify-end">
      <n-button class="mr-4" type="warning" @click="handleCancel">取消</n-button>
      <n-button type="info" @click="handleConfirm">继续</n-button>
    </div>
  </n-modal>
</template>
