<script setup>
import { ref, watch } from "vue";
import { useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";
import ApiService from "@/service/api";

const props = defineProps({
  show: Boolean,
  initialPath: { type: String, default: "" },
});
const emit = defineEmits(["update:show", "select"]);
const { t } = useI18n();
const message = useMessage();
const loading = ref(false);
const current = ref("");
const parent = ref("");
const entries = ref([]);
const roots = ref([]);

const load = async (path = "") => {
  loading.value = true;
  const { data, statusCode } = await new ApiService().listDirectories(path);
  loading.value = false;
  if (statusCode.value !== 200) {
    message.error(data.value?.error || t("configuration.directoryLoadFailed"));
    return;
  }
  current.value = data.value.current;
  parent.value = data.value.parent;
  entries.value = data.value.entries || [];
  roots.value = data.value.roots || [];
};

watch(
  () => props.show,
  (show) => {
    if (show) load(props.initialPath);
  },
);
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    class="directory-card"
    style="width: min(92vw, 680px); max-height: calc(100vh - 64px)"
    :title="$t('configuration.chooseDirectory')"
    @update:show="emit('update:show', $event)"
  >
    <n-input-group class="mb-3">
      <n-input v-model:value="current" @keyup.enter="load(current)" />
      <n-button :loading="loading" @click="load(current)">{{
        $t("configuration.go")
      }}</n-button>
    </n-input-group>
    <n-space class="mb-2">
      <n-button
        size="small"
        :disabled="parent === current"
        @click="load(parent)"
      >
        {{ $t("configuration.parentDirectory") }}
      </n-button>
      <n-button size="small" type="primary" @click="emit('select', current)">
        {{ $t("configuration.selectCurrent") }}
      </n-button>
      <n-button
        v-for="root in roots"
        :key="root"
        size="small"
        quaternary
        @click="load(root)"
      >
        {{ root }}
      </n-button>
    </n-space>
    <n-spin :show="loading">
      <n-scrollbar style="max-height: min(52vh, 420px)">
        <n-empty
          v-if="!entries.length"
          :description="$t('configuration.noSubdirectories')"
        />
        <n-list v-else hoverable clickable>
          <n-list-item
            v-for="entry in entries"
            :key="entry.path"
            @click="load(entry.path)"
          >
            📁 {{ entry.name }}
          </n-list-item>
        </n-list>
      </n-scrollbar>
    </n-spin>
  </n-modal>
</template>
