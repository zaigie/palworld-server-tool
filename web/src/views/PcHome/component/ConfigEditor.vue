<script setup>
import { ref, onMounted, computed } from "vue";
import { useI18n } from "vue-i18n";
import { useMessage } from "naive-ui";
import ApiService from "@/service/api";

const { t } = useI18n();
const message = useMessage();
const loading = ref(false);
const originalConfig = ref({});
const showSaveConfirm = ref(false);
const activeTab = ref("server"); // Default to server settings tab

// Settings categories
const serverSettings = ref([]);
const gameSettings = ref([]);
const advancedSettings = ref([]);

// Load the server configuration
const loadConfig = async () => {
  loading.value = true;
  try {
    const { data, statusCode } = await new ApiService().getServerConfig();
    if (statusCode.value === 200 && data.value) {
      if (data.value.config) {
        parseConfig(data.value.config);
        // Store original config for comparison
        originalConfig.value = {
          server: JSON.parse(JSON.stringify(serverSettings.value)),
          game: JSON.parse(JSON.stringify(gameSettings.value)),
          advanced: JSON.parse(JSON.stringify(advancedSettings.value))
        };
      } else {
        message.error(t("config.invalidResponse"));
      }
    } else {
      message.error(t("config.loadError"));
    }
  } catch (error) {
    message.error(t("config.loadError") + ": " + error.message);
  } finally {
    loading.value = false;
  }
};

// Parse the INI configuration into structured settings
const parseConfig = (iniString) => {
  try {
    console.log("Raw config string:", iniString);
    
    // Extract options from the configuration string
    const optionsMatch = iniString.match(/OptionSettings=\((.*)\)/);
    
    if (optionsMatch && optionsMatch[1]) {
      const optionsContent = optionsMatch[1];
      console.log("Options content:", optionsContent);
      
      const options = {};
      
      // Parse the options content
      let currentKey = '';
      let currentValue = '';
      let inQuotes = false;
      let buffer = '';
      
      for (let i = 0; i < optionsContent.length; i++) {
        const char = optionsContent[i];
        
        if (char === '"') {
          inQuotes = !inQuotes;
          buffer += char;
        } else if (char === '=' && !inQuotes && currentKey === '') {
          currentKey = buffer.trim();
          buffer = '';
        } else if (char === ',' && !inQuotes) {
          currentValue = buffer.trim();
          if (currentKey && currentValue) {
            // Remove quotes if present
            if (currentValue.startsWith('"') && currentValue.endsWith('"')) {
              currentValue = currentValue.substring(1, currentValue.length - 1);
            }
            options[currentKey] = currentValue;
          }
          currentKey = '';
          currentValue = '';
          buffer = '';
        } else {
          buffer += char;
        }
      }
      
      // Handle the last key-value pair
      if (currentKey && buffer) {
        currentValue = buffer.trim();
        if (currentValue.startsWith('"') && currentValue.endsWith('"')) {
          currentValue = currentValue.substring(1, currentValue.length - 1);
        }
        options[currentKey] = currentValue;
      }
      
      console.log("All parsed options:", options);
      
      // Categorize settings
      categorizeSettings(options);
    } else {
      console.error("Could not extract options content from:", iniString);
      message.error(t("config.invalidFormat"));
    }
  } catch (error) {
    console.error("Error parsing config:", error);
    message.error(t("config.parseError") + ": " + error.message);
  }
};

// Categorize settings into server, game, and advanced tabs
const categorizeSettings = (options) => {
  // Server settings
  serverSettings.value = [
    { key: "ServerName", label: "Server Name", value: options.ServerName || "", type: "text" },
    { key: "ServerDescription", label: "Server Description", value: options.ServerDescription || "", type: "text" },
    { key: "AdminPassword", label: "Admin Password", value: options.AdminPassword || "", type: "text" },
    { key: "ServerPassword", label: "Server Password", value: options.ServerPassword || "", type: "text" },
    { key: "PublicIP", label: "Public IP", value: options.PublicIP || "", type: "text" },
    { key: "PublicPort", label: "Public Port", value: options.PublicPort || "8211", type: "text" },
    { key: "RCONEnabled", label: "RCON Enabled", value: options.RCONEnabled || "False", type: "text" },
    { key: "RCONPort", label: "RCON Port", value: options.RCONPort || "25575", type: "text" },
    { key: "RESTAPIEnabled", label: "REST API Enabled", value: options.RESTAPIEnabled || "False", type: "text" },
    { key: "RESTAPIPort", label: "REST API Port", value: options.RESTAPIPort || "8212", type: "text" },
    { key: "Region", label: "Region", value: options.Region || "", type: "text" },
    { key: "bUseAuth", label: "Use Authentication", value: options.bUseAuth || "True", type: "text" },
    { key: "BanListURL", label: "Ban List URL", value: options.BanListURL || "", type: "text" },
    { key: "ServerPlayerMaxNum", label: "Server Player Max Num", value: options.ServerPlayerMaxNum || "32", type: "text" },
    { key: "bIsUseBackupSaveData", label: "Use Backup Save Data", value: options.bIsUseBackupSaveData || "True", type: "text" },
    { key: "AutoSaveSpan", label: "Auto Save Span", value: options.AutoSaveSpan || "30.000000", type: "text" },
    { key: "LogFormatType", label: "Log Format Type", value: options.LogFormatType || "Text", type: "text" },
    { key: "bShowPlayerList", label: "Show Player List", value: options.bShowPlayerList || "False", type: "text" },
    { key: "CrossplayPlatforms", label: "Crossplay Platforms", value: options.CrossplayPlatforms || "(Steam,Xbox,PS5,Mac)", type: "text" },
  ];
  
  // Game settings
  gameSettings.value = [
    { key: "Difficulty", label: "Difficulty", value: options.Difficulty || "None", type: "text" },
    { key: "RandomizerType", label: "Randomizer Type", value: options.RandomizerType || "None", type: "text" },
    { key: "RandomizerSeed", label: "Randomizer Seed", value: options.RandomizerSeed || "", type: "text" },
    { key: "bIsRandomizerPalLevelRandom", label: "Randomizer Pal Level Random", value: options.bIsRandomizerPalLevelRandom || "False", type: "text" },
    { key: "DayTimeSpeedRate", label: "Day Time Speed Rate", value: options.DayTimeSpeedRate || "1.000000", type: "text" },
    { key: "NightTimeSpeedRate", label: "Night Time Speed Rate", value: options.NightTimeSpeedRate || "1.000000", type: "text" },
    { key: "ExpRate", label: "Exp Rate", value: options.ExpRate || "1.000000", type: "text" },
    { key: "PalCaptureRate", label: "Pal Capture Rate", value: options.PalCaptureRate || "1.000000", type: "text" },
    { key: "PalSpawnNumRate", label: "Pal Spawn Number Rate", value: options.PalSpawnNumRate || "1.000000", type: "text" },
    { key: "PalDamageRateAttack", label: "Pal Damage Rate Attack", value: options.PalDamageRateAttack || "1.000000", type: "text" },
    { key: "PalDamageRateDefense", label: "Pal Damage Rate Defense", value: options.PalDamageRateDefense || "1.000000", type: "text" },
    { key: "PlayerDamageRateAttack", label: "Player Damage Rate Attack", value: options.PlayerDamageRateAttack || "1.000000", type: "text" },
    { key: "PlayerDamageRateDefense", label: "Player Damage Rate Defense", value: options.PlayerDamageRateDefense || "1.000000", type: "text" },
    { key: "PlayerStomachDecreaceRate", label: "Player Stomach Decrease Rate", value: options.PlayerStomachDecreaceRate || "1.000000", type: "text" },
    { key: "PlayerStaminaDecreaceRate", label: "Player Stamina Decrease Rate", value: options.PlayerStaminaDecreaceRate || "1.000000", type: "text" },
    { key: "PlayerAutoHPRegeneRate", label: "Player Auto HP Regeneration Rate", value: options.PlayerAutoHPRegeneRate || "1.000000", type: "text" },
    { key: "PlayerAutoHpRegeneRateInSleep", label: "Player Auto HP Regeneration Rate In Sleep", value: options.PlayerAutoHpRegeneRateInSleep || "1.000000", type: "text" },
    { key: "PalStomachDecreaceRate", label: "Pal Stomach Decrease Rate", value: options.PalStomachDecreaceRate || "1.000000", type: "text" },
    { key: "PalStaminaDecreaceRate", label: "Pal Stamina Decrease Rate", value: options.PalStaminaDecreaceRate || "1.000000", type: "text" },
    { key: "PalAutoHPRegeneRate", label: "Pal Auto HP Regeneration Rate", value: options.PalAutoHPRegeneRate || "1.000000", type: "text" },
    { key: "PalAutoHpRegeneRateInSleep", label: "Pal Auto HP Regeneration Rate In Sleep", value: options.PalAutoHpRegeneRateInSleep || "1.000000", type: "text" },
    { key: "BuildObjectHpRate", label: "Build Object HP Rate", value: options.BuildObjectHpRate || "1.000000", type: "text" },
    { key: "BuildObjectDamageRate", label: "Build Object Damage Rate", value: options.BuildObjectDamageRate || "1.000000", type: "text" },
    { key: "BuildObjectDeteriorationDamageRate", label: "Build Object Deterioration Damage Rate", value: options.BuildObjectDeteriorationDamageRate || "1.000000", type: "text" },
    { key: "CollectionDropRate", label: "Collection Drop Rate", value: options.CollectionDropRate || "1.000000", type: "text" },
    { key: "CollectionObjectHpRate", label: "Collection Object HP Rate", value: options.CollectionObjectHpRate || "1.000000", type: "text" },
    { key: "CollectionObjectRespawnSpeedRate", label: "Collection Object Respawn Speed Rate", value: options.CollectionObjectRespawnSpeedRate || "1.000000", type: "text" },
    { key: "EnemyDropItemRate", label: "Enemy Drop Item Rate", value: options.EnemyDropItemRate || "1.000000", type: "text" },
    { key: "DeathPenalty", label: "Death Penalty", value: options.DeathPenalty || "All", type: "text" },
    { key: "PalEggDefaultHatchingTime", label: "Pal Egg Default Hatching Time", value: options.PalEggDefaultHatchingTime || "72.000000", type: "text" },
    { key: "WorkSpeedRate", label: "Work Speed Rate", value: options.WorkSpeedRate || "1.000000", type: "text" },
    { key: "ItemWeightRate", label: "Item Weight Rate", value: options.ItemWeightRate || "1.000000", type: "text" },
  ];
  
  // Advanced settings
  advancedSettings.value = [
    { key: "bEnablePlayerToPlayerDamage", label: "Enable Player To Player Damage", value: options.bEnablePlayerToPlayerDamage || "False", type: "text" },
    { key: "bEnableFriendlyFire", label: "Enable Friendly Fire", value: options.bEnableFriendlyFire || "False", type: "text" },
    { key: "bEnableInvaderEnemy", label: "Enable Invader Enemy", value: options.bEnableInvaderEnemy || "True", type: "text" },
    { key: "bActiveUNKO", label: "Active UNKO (Pal will poop)", value: options.bActiveUNKO || "False", type: "text" },
    { key: "bEnableAimAssistPad", label: "Enable Aim Assist Controller", value: options.bEnableAimAssistPad || "True", type: "text" },
    { key: "bEnableAimAssistKeyboard", label: "Enable Aim Assist Keyboard", value: options.bEnableAimAssistKeyboard || "False", type: "text" },
    { key: "DropItemMaxNum", label: "Drop Item Max Num", value: options.DropItemMaxNum || "3000", type: "text" },
    { key: "DropItemMaxNum_UNKO", label: "Drop Item Max Num (UNKO)", value: options.DropItemMaxNum_UNKO || "100", type: "text" },
    { key: "BaseCampMaxNum", label: "Base Camp Max Num", value: options.BaseCampMaxNum || "128", type: "text" },
    { key: "BaseCampWorkerMaxNum", label: "Base Camp Worker Max Num", value: options.BaseCampWorkerMaxNum || "15", type: "text" },
    { key: "DropItemAliveMaxHours", label: "Drop Item Alive Max Hours", value: options.DropItemAliveMaxHours || "1.000000", type: "text" },
    { key: "bAutoResetGuildNoOnlinePlayers", label: "Auto Reset Guild No Online Players", value: options.bAutoResetGuildNoOnlinePlayers || "False", type: "text" },
    { key: "AutoResetGuildTimeNoOnlinePlayers", label: "Auto Reset Guild Time No Online Players", value: options.AutoResetGuildTimeNoOnlinePlayers || "72.000000", type: "text" },
    { key: "GuildPlayerMaxNum", label: "Guild Player Max Num", value: options.GuildPlayerMaxNum || "20", type: "text" },
    { key: "BaseCampMaxNumInGuild", label: "Base Camp Max Num In Guild", value: options.BaseCampMaxNumInGuild || "4", type: "text" },
    { key: "bIsMultiplay", label: "Is Multiplay", value: options.bIsMultiplay || "False", type: "text" },
    { key: "bIsPvP", label: "Is Pv P", value: options.bIsPvP || "False", type: "text" },
    { key: "bHardcore", label: "Hardcore Mode", value: options.bHardcore || "False", type: "text" },
    { key: "bPalLost", label: "Pal Lost Mode", value: options.bPalLost || "False", type: "text" },
    { key: "bCharacterRecreateInHardcore", label: "Character Recreate In Hardcore", value: options.bCharacterRecreateInHardcore || "False", type: "text" },
    { key: "bCanPickupOtherGuildDeathPenaltyDrop", label: "Can Pickup Other Guild Death Penalty Drop", value: options.bCanPickupOtherGuildDeathPenaltyDrop || "False", type: "text" },
    { key: "bEnableNonLoginPenalty", label: "Enable Non Login Penalty", value: options.bEnableNonLoginPenalty || "True", type: "text" },
    { key: "bEnableFastTravel", label: "Enable Fast Travel", value: options.bEnableFastTravel || "True", type: "text" },
    { key: "bIsStartLocationSelectByMap", label: "Is Start Location Select By Map", value: options.bIsStartLocationSelectByMap || "True", type: "text" },
    { key: "bExistPlayerAfterLogout", label: "Exist Player After Logout", value: options.bExistPlayerAfterLogout || "False", type: "text" },
    { key: "bEnableDefenseOtherGuildPlayer", label: "Enable Defense Other Guild Player", value: options.bEnableDefenseOtherGuildPlayer || "False", type: "text" },
    { key: "bInvisibleOtherGuildBaseCampAreaFX", label: "Invisible Other Guild Base Camp Area FX", value: options.bInvisibleOtherGuildBaseCampAreaFX || "False", type: "text" },
    { key: "bBuildAreaLimit", label: "Build Area Limit", value: options.bBuildAreaLimit || "False", type: "text" },
    { key: "CoopPlayerMaxNum", label: "Coop Player Max Num", value: options.CoopPlayerMaxNum || "4", type: "text" },
    { key: "ChatPostLimitPerMinute", label: "Chat Post Limit Per Minute", value: options.ChatPostLimitPerMinute || "10", type: "text" },
    { key: "SupplyDropSpan", label: "Supply Drop Span", value: options.SupplyDropSpan || "180", type: "text" },
    { key: "EnablePredatorBossPal", label: "Enable Predator Boss Pal", value: options.EnablePredatorBossPal || "True", type: "text" },
    { key: "MaxBuildingLimitNum", label: "Max Building Limit Num", value: options.MaxBuildingLimitNum || "0", type: "text" },
    { key: "ServerReplicatePawnCullDistance", label: "Server Replicate Pawn Cull Distance", value: options.ServerReplicatePawnCullDistance || "15000.000000", type: "text" },
    { key: "bAllowGlobalPalboxExport", label: "Allow Global Palbox Export", value: options.bAllowGlobalPalboxExport || "True", type: "text" },
    { key: "bAllowGlobalPalboxImport", label: "Allow Global Palbox Import", value: options.bAllowGlobalPalboxImport || "False", type: "text" },
  ];
};

// Convert settings back to INI format
const generateIniConfig = () => {
  // Combine all settings
  const allSettings = {};
  
  // Process all settings
  const processSettings = (settings) => {
    settings.forEach(setting => {
      let value = setting.value.toString();
      
      // Add quotes for text fields that need them
      if (!value.startsWith("True") && 
          !value.startsWith("False") && 
          !value.startsWith("(") && 
          value !== "None" && 
          value !== "All" && 
          value !== "Text" &&
          !value.match(/^[0-9.]+$/)) {
        value = `"${value}"`;
      }
      
      allSettings[setting.key] = value;
    });
  };
  
  // Process all setting categories
  processSettings(serverSettings.value);
  processSettings(gameSettings.value);
  processSettings(advancedSettings.value);
  
  // Build the OptionSettings string
  const optionParts = [];
  for (const [key, value] of Object.entries(allSettings)) {
    optionParts.push(`${key}=${value}`);
  }
  
  // Construct the INI content
  let iniContent = "[/Script/Pal.PalGameWorldSettings]\n";
  iniContent += `OptionSettings=(${optionParts.join(',')})\n`;
  
  return iniContent;
};

// Save the configuration
const saveConfig = async () => {
  loading.value = true;
  try {
    const iniString = generateIniConfig();
    const { data, statusCode } = await new ApiService().saveServerConfig({ config: iniString });
    
    if (statusCode.value === 200) {
      message.success(t("config.saveSuccess"));
      // Update original config to match current state
      originalConfig.value = {
        server: JSON.parse(JSON.stringify(serverSettings.value)),
        game: JSON.parse(JSON.stringify(gameSettings.value)),
        advanced: JSON.parse(JSON.stringify(advancedSettings.value))
      };
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

// Check if config has been modified
const hasChanges = computed(() => {
  const originalServer = JSON.stringify(originalConfig.value.server || []);
  const originalGame = JSON.stringify(originalConfig.value.game || []);
  const originalAdvanced = JSON.stringify(originalConfig.value.advanced || []);
  
  const currentServer = JSON.stringify(serverSettings.value);
  const currentGame = JSON.stringify(gameSettings.value);
  const currentAdvanced = JSON.stringify(advancedSettings.value);
  
  return originalServer !== currentServer || 
         originalGame !== currentGame || 
         originalAdvanced !== currentAdvanced;
});

// Handle save button click
const handleSave = () => {
  if (hasChanges.value) {
    showSaveConfirm.value = true;
  } else {
    message.info(t("config.noChanges"));
  }
};

// Reset to original values
const resetConfig = () => {
  if (originalConfig.value.server) {
    serverSettings.value = JSON.parse(JSON.stringify(originalConfig.value.server));
  }
  if (originalConfig.value.game) {
    gameSettings.value = JSON.parse(JSON.stringify(originalConfig.value.game));
  }
  if (originalConfig.value.advanced) {
    advancedSettings.value = JSON.parse(JSON.stringify(originalConfig.value.advanced));
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
            <n-button @click="resetConfig" :disabled="!hasChanges">
              {{ $t('button.cancel') }}
            </n-button>
            <n-button type="primary" @click="handleSave">
              {{ $t('button.save') }}
            </n-button>
          </n-space>
        </template>
        
        <n-tabs v-model:value="activeTab" type="segment" animated>
          <n-tab-pane name="server" tab="Server Settings">
            <div class="tab-content">
              <n-form-item 
                v-for="setting in serverSettings" 
                :key="setting.key" 
                :label="setting.label"
              >
                <n-input 
                  v-model:value="setting.value" 
                  :placeholder="setting.label"
                />
              </n-form-item>
            </div>
          </n-tab-pane>
          
          <n-tab-pane name="game" tab="In-Game Settings">
            <div class="tab-content">
              <n-form-item 
                v-for="setting in gameSettings" 
                :key="setting.key" 
                :label="setting.label"
              >
                <n-input 
                  v-model:value="setting.value" 
                  :placeholder="setting.label"
                />
              </n-form-item>
            </div>
          </n-tab-pane>
          
          <n-tab-pane name="advanced" tab="Advanced Settings">
            <div class="tab-content">
              <n-form-item 
                v-for="setting in advancedSettings" 
                :key="setting.key" 
                :label="setting.label"
              >
                <n-input 
                  v-model:value="setting.value" 
                  :placeholder="setting.label"
                />
              </n-form-item>
            </div>
          </n-tab-pane>
        </n-tabs>
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
  height: 100%;
}

.tab-content {
  height: calc(100vh - 200px);
  overflow-y: auto;
  padding: 20px;
  padding-bottom: 60px;
}

:deep(.n-form-item) {
  margin-bottom: 16px;
}

:deep(.n-input) {
  width: 100%;
}
</style>
