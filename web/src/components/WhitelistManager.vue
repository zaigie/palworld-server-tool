<script setup>
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useMessage } from "naive-ui";
import ApiService from "@/service/api";
import whitelistStore from "@/stores/model/whitelist";

const props = defineProps({
  show: Boolean,
  players: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:show", "updated"]);
const { t } = useI18n();
const message = useMessage();
const api = new ApiService();
const rows = ref([]);
const selectedPlayer = ref(null);
const loading = ref(false);
const saving = ref(false);

const playerOptions = computed(() =>
  props.players
    .filter(
      (player) =>
        !rows.value.some(
          (row) =>
            row.player_uid === player.player_uid ||
            (row.steam_id && row.steam_id === player.steam_id)
        )
    )
    .map((player) => ({
      label: `${player.nickname || player.player_uid} · ${player.player_uid}`,
      value: player.player_uid,
    }))
);

const load = async () => {
  loading.value = true;
  try {
    const { data } = await api.getWhitelist();
    const whitelist = Array.isArray(data.value) ? data.value : [];
    rows.value = whitelist.map((row) => ({ ...row }));
    whitelistStore().setWhitelist(rows.value);
  } finally {
    loading.value = false;
  }
};

const addSelectedPlayer = () => {
  const player = props.players.find(
    (candidate) => candidate.player_uid === selectedPlayer.value
  );
  if (!player) return;
  rows.value.unshift({
    name: player.nickname || "",
    player_uid: player.player_uid || "",
    steam_id: player.steam_id || "",
  });
  selectedPlayer.value = null;
};

const save = async () => {
  const validRows = rows.value.filter((row) => row.player_uid || row.steam_id);
  saving.value = true;
  try {
    const { data, statusCode } = await api.putWhitelist(
      JSON.stringify(validRows)
    );
    if (statusCode.value === 200) {
      rows.value = validRows;
      whitelistStore().setWhitelist(validRows);
      message.success(t("message.addwhitesuccess"));
      emit("updated");
      emit("update:show", false);
    } else {
      message.error(t("message.addwhitefail", { err: data.value?.error }));
    }
  } finally {
    saving.value = false;
  }
};

watch(
  () => props.show,
  (show) => show && load(),
  { immediate: true }
);
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    style="width: 94%; max-width: 780px"
    :title="$t('modal.whitelist')"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <n-spin :show="loading">
      <n-space vertical size="large">
        <n-input-group>
          <n-select
            v-model:value="selectedPlayer"
            :options="playerOptions"
            filterable
            clearable
            :placeholder="$t('whitelist.selectExisting')"
          />
          <n-button type="primary" :disabled="!selectedPlayer" @click="addSelectedPlayer">
            {{ $t("button.add") }}
          </n-button>
        </n-input-group>
        <n-empty
          v-if="rows.length === 0"
          :description="$t('whitelist.emptyDescription')"
        >
          <template #extra>{{ $t("whitelist.emptyHint") }}</template>
        </n-empty>
        <n-space v-else vertical>
          <n-card v-for="(row, index) in rows" :key="`${row.player_uid}-${index}`" size="small">
            <n-grid cols="1 640:3" :x-gap="8" :y-gap="8">
              <n-gi><n-input v-model:value="row.name" :placeholder="$t('input.nickname')" /></n-gi>
              <n-gi><n-input v-model:value="row.player_uid" :placeholder="$t('input.player_uid')" /></n-gi>
              <n-gi>
                <n-input-group>
                  <n-input v-model:value="row.steam_id" :placeholder="$t('input.steam_id')" />
                  <n-button type="error" secondary @click="rows.splice(index, 1)">
                    {{ $t("button.remove") }}
                  </n-button>
                </n-input-group>
              </n-gi>
            </n-grid>
          </n-card>
        </n-space>
      </n-space>
    </n-spin>
    <template #footer>
      <n-flex justify="space-between">
        <n-button
          @click="rows.unshift({ name: '', player_uid: '', steam_id: '' })"
        >{{ $t("whitelist.manualAdd") }}</n-button>
        <n-flex>
          <n-button @click="emit('update:show', false)">{{ $t("button.cancel") }}</n-button>
          <n-button type="primary" :loading="saving" @click="save">
            {{ $t("button.save") }}
          </n-button>
        </n-flex>
      </n-flex>
    </template>
  </n-modal>
</template>
