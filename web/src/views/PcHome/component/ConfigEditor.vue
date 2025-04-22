<script setup>
import { ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useMessage } from "naive-ui";
import ApiService from "@/service/api";

const { t } = useI18n();
const message = useMessage();
const configSections = ref([]);
const originalConfig = ref({});
const loading = ref(false);
const showSaveConfirm = ref(false);

// Load the server configuration
const loadConfig = async () => {
  loading.value = true;
  try {
    const { data, statusCode } = await new ApiService().getServerConfig();
    if (statusCode.value === 200 && data.value) {
      const parsed = parseIniConfig(data.value);
      configSections.value = parsed;
      originalConfig.value = JSON.parse(JSON.stringify(parsed));
    } else {
      message.error(t("config.loadError"));
    }
  } catch (error) {
    message.error(t("config.loadError") + ": " + error.message);
  } finally {
    loading.value = false;
  }
};

// Parse INI config string into structured sections
const parseIniConfig = (iniString) => {
  const sections = [];
  let currentSection = null;
  
  // Check if iniString is a valid string
  if (!iniString || typeof iniString !== 'string') {
    message.error(t("config.invalidFormat"));
    return sections;
  }
  
  try {
    iniString.split('\n').forEach(line => {
      line = line.trim();
      
      // Skip empty lines and comments
      if (!line || line.startsWith(';') || line.startsWith('#')) return;
      
      // Check if line is a section header
      const sectionMatch = line.match(/^\[([^\]]+)\]/);
      if (sectionMatch) {
        currentSection = {
          name: sectionMatch[1],
          settings: []
        };
        sections.push(currentSection);
        return;
      }
      
      // Process key-value pairs
      if (currentSection) {
        const keyValueMatch = line.match(/^([^=]+)=(.*)$/);
        if (keyValueMatch) {
          const key = keyValueMatch[1].trim();
          const value = keyValueMatch[2].trim();
          currentSection.settings.push({
            key,
            value,
            originalValue: value
          });
        }
      }
    });
  } catch (error) {
    message.error(t("config.parseError") + ": " + error.message);
  }
  
  return sections;
};

// Convert structured sections back to INI string
const convertToIniString = (sections) => {
  let iniString = '';
  
  sections.forEach(section => {
    iniString += `[${section.name}]\n`;
    
    section.settings.forEach(setting => {
      iniString += `${setting.key}=${setting.value}\n`;
    });
    
    iniString += '\n';
  });
  
  return iniString;
};

// Save the configuration
const saveConfig = async () => {
  loading.value = true;
  try {
    const iniString = convertToIniString(configSections.value);
    const { data, statusCode } = await new ApiService().saveServerConfig({ config: iniString });
    
    if (statusCode.value === 200) {
      message.success(t("config.saveSuccess"));
      originalConfig.value = JSON.parse(JSON.stringify(configSections.value));
      showSaveConfirm.value = false;
    } else {
      message.error(t("config.saveError") + ": " + data.value?.error);
    }
  } catch (error) {
    message.error(t("config.saveError") + ": " + error.message);
  } finally {
    loading.value = false;
  }
};

// Reset to original values
const resetConfig = () => {
  configSections.value = JSON.parse(JSON.stringify(originalConfig.value));
};

// Check if config has been modified
const hasChanges = () => {
  return JSON.stringify(configSections.value) !== JSON.stringify(originalConfig.value);
};

// Handle save button click
const handleSave = () => {
  if (hasChanges()) {
    showSaveConfirm.value = true;
  } else {
    message.info(t("config.noChanges"));
  }
};

onMounted(() => {
  loadConfig();
});
</script>

<template>
  <div class="config-editor">
    <n-spin :show="loading">
      <n-card :title="$t('config.title')" size="huge">
        <template #header-extra>
          <n-space>
            <n-button @click="resetConfig" :disabled="!hasChanges()">
              {{ $t('button.cancel') }}
            </n-button>
            <n-button type="primary" @click="handleSave">
              {{ $t('button.save') }}
            </n-button>
          </n-space>
        </template>
        
        <n-collapse accordion>
          <n-collapse-item 
            v-for="(section, index) in configSections" 
            :key="index" 
            :title="section.name"
          >
            <n-grid :cols="1" :x-gap="12">
              <n-gi v-for="(setting, settingIndex) in section.settings" :key="settingIndex">
                <n-form-item :label="setting.key">
                  <n-input 
                    v-model:value="setting.value" 
                    :placeholder="setting.key"
                    :status="setting.value !== setting.originalValue ? 'warning' : undefined"
                  />
                </n-form-item>
              </n-gi>
            </n-grid>
          </n-collapse-item>
        </n-collapse>
      </n-card>
    </n-spin>
    
    <!-- Confirmation Dialog -->
    <n-modal v-model:show="showSaveConfirm">
      <n-card
        style="width: 450px"
        :title="$t('config.confirmSave')"
        :bordered="false"
        size="huge"
        role="dialog"
        aria-modal="true"
      >
        <template #header-extra>
          <n-button circle size="small" @click="showSaveConfirm = false">
            <template #icon>
              <n-icon><close /></n-icon>
            </template>
          </n-button>
        </template>
        
        <n-space vertical>
          <span>{{ $t('config.saveWarning') }}</span>
          <n-space justify="end">
            <n-button @click="showSaveConfirm = false">
              {{ $t('button.cancel') }}
            </n-button>
            <n-button type="primary" @click="saveConfig">
              {{ $t('button.confirm') }}
            </n-button>
          </n-space>
        </n-space>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>
.config-editor {
  padding: 16px;
}
</style>
