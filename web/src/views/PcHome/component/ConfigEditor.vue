<script setup>
import { ref, onMounted, computed } from "vue";
import { useI18n } from "vue-i18n";
import { useMessage } from "naive-ui";
import ApiService from "@/service/api";
import { Refresh, Close } from "@vicons/ionicons5";

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
    
    // The PalWorld config has a specific format
    // First, look for the OptionSettings part
    const optionsMatch = iniString.match(/OptionSettings=\((.*)\)/);
    
    if (optionsMatch && optionsMatch[1]) {
      const optionsContent = optionsMatch[1];
      console.log("Options content:", optionsContent);
      
      const options = {};
      
      // Split by commas, but be careful with values that might contain commas
      let currentKey = '';
      let currentValue = '';
      let inQuotes = false;
      let keyValuePairs = [];
      
      // Simple split by comma for most cases
      optionsContent.split(',').forEach(option => {
        const parts = option.split('=');
        if (parts.length === 2) {
          const key = parts[0].trim();
          const value = parts[1].trim();
          options[key] = value;
        }
      });
      
      console.log("Parsed options:", options);
      
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
    { 
      key: "ServerName", 
      label: "Server Name", 
      value: options.ServerName || "Default Palworld Server",
      type: "text"
    },
    { 
      key: "ServerDescription", 
      label: "Server Description", 
      value: options.ServerDescription || "",
      type: "textarea"
    },
    { 
      key: "AdminPassword", 
      label: "Admin Password", 
      value: options.AdminPassword || "",
      type: "password"
    },
    { 
      key: "ServerPassword", 
      label: "Server Password", 
      value: options.ServerPassword || "",
      type: "password"
    },
    { 
      key: "PublicIP", 
      label: "Public IP", 
      value: options.PublicIP || "",
      type: "text"
    },
    { 
      key: "PublicPort", 
      label: "Public Port", 
      value: options.PublicPort || "8211",
      type: "number"
    },
    { 
      key: "ServerPlayerMaxNum", 
      label: "Server Player Max Num", 
      value: options.ServerPlayerMaxNum || "32",
      min: 1,
      max: 32,
      type: "slider"
    },
    { 
      key: "bIsUseBackupSaveData", 
      label: "Use Backup Save Data", 
      value: options.bIsUseBackupSaveData === "True",
      type: "switch"
    },
    { 
      key: "Region", 
      label: "Region", 
      value: options.Region || "",
      type: "text"
    },
    { 
      key: "bUseAuth", 
      label: "Use Authentication", 
      value: options.bUseAuth === "True",
      type: "switch"
    },
    { 
      key: "AllowConnectPlatform", 
      label: "Allow Connect Platform", 
      value: options.AllowConnectPlatform || "Steam",
      type: "select",
      options: [
        { label: "Steam Only", value: "Steam" },
        { label: "All Platforms", value: "All" }
      ]
    },
    { 
      key: "bShowPlayerList", 
      label: "Show Player List", 
      value: options.bShowPlayerList === "True",
      type: "switch"
    },
    { 
      key: "RCONEnabled", 
      label: "RCON Enabled", 
      value: options.RCONEnabled === "True",
      type: "switch"
    },
    { 
      key: "RCONPort", 
      label: "RCON Port", 
      value: options.RCONPort || "25575",
      type: "number"
    },
    { 
      key: "RESTAPIEnabled", 
      label: "REST API Enabled", 
      value: options.RESTAPIEnabled === "True",
      type: "switch"
    },
    { 
      key: "RESTAPIPort", 
      label: "REST API Port", 
      value: options.RESTAPIPort || "8212",
      type: "number"
    },
    { 
      key: "BanListURL", 
      label: "Ban List URL", 
      value: options.BanListURL || "https://api.palworldgame.com/api/banlist.txt",
      type: "text"
    }
  ];
  
  // Game settings
  gameSettings.value = [
    { 
      key: "Difficulty", 
      label: "Difficulty", 
      value: options.Difficulty || "None",
      type: "select",
      options: [
        { label: "None", value: "None" },
        { label: "Easy", value: "Easy" },
        { label: "Normal", value: "Normal" },
        { label: "Hard", value: "Hard" }
      ]
    },
    { 
      key: "DayTimeSpeedRate", 
      label: "Day Time Speed Rate", 
      value: options.DayTimeSpeedRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "NightTimeSpeedRate", 
      label: "Night Time Speed Rate", 
      value: options.NightTimeSpeedRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "ExpRate", 
      label: "Exp Rate", 
      value: options.ExpRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "PalCaptureRate", 
      label: "Pal Capture Rate", 
      value: options.PalCaptureRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "PalSpawnNumRate", 
      label: "Pal Spawn Number Rate", 
      value: options.PalSpawnNumRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "PalDamageRateAttack", 
      label: "Pal Damage Rate Attack", 
      value: options.PalDamageRateAttack || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "PalDamageRateDefense", 
      label: "Pal Damage Rate Defense", 
      value: options.PalDamageRateDefense || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "PlayerDamageRateAttack", 
      label: "Player Damage Rate Attack", 
      value: options.PlayerDamageRateAttack || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "PlayerDamageRateDefense", 
      label: "Player Damage Rate Defense", 
      value: options.PlayerDamageRateDefense || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "PlayerStomachDecreaceRate", 
      label: "Player Stomach Decrease Rate", 
      value: options.PlayerStomachDecreaceRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      easier: true,
      type: "slider"
    },
    { 
      key: "PlayerStaminaDecreaceRate", 
      label: "Player Stamina Decrease Rate", 
      value: options.PlayerStaminaDecreaceRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      easier: true,
      type: "slider"
    },
    { 
      key: "PlayerAutoHPRegeneRate", 
      label: "Player Auto HP Regeneration Rate", 
      value: options.PlayerAutoHPRegeneRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "PlayerAutoHpRegeneRateInSleep", 
      label: "Player Auto HP Regeneration Rate In Sleep", 
      value: options.PlayerAutoHpRegeneRateInSleep || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "PalStomachDecreaceRate", 
      label: "Pal Stomach Decrease Rate", 
      value: options.PalStomachDecreaceRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      easier: true,
      type: "slider"
    },
    { 
      key: "PalStaminaDecreaceRate", 
      label: "Pal Stamina Decrease Rate", 
      value: options.PalStaminaDecreaceRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      easier: true,
      type: "slider"
    },
    { 
      key: "PalAutoHPRegeneRate", 
      label: "Pal Auto HP Regeneration Rate", 
      value: options.PalAutoHPRegeneRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "PalAutoHpRegeneRateInSleep", 
      label: "Pal Auto HP Regeneration Rate In Sleep", 
      value: options.PalAutoHpRegeneRateInSleep || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "BuildObjectDamageRate", 
      label: "Build Object Damage Rate", 
      value: options.BuildObjectDamageRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "BuildObjectDeteriorationDamageRate", 
      label: "Build Object Deterioration Damage Rate", 
      value: options.BuildObjectDeteriorationDamageRate || "1.000000",
      min: 0,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "CollectionDropRate", 
      label: "Collection Drop Rate", 
      value: options.CollectionDropRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "CollectionObjectHpRate", 
      label: "Collection Object HP Rate", 
      value: options.CollectionObjectHpRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "CollectionObjectRespawnSpeedRate", 
      label: "Collection Object Respawn Speed Rate", 
      value: options.CollectionObjectRespawnSpeedRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "EnemyDropItemRate", 
      label: "Enemy Drop Item Rate", 
      value: options.EnemyDropItemRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    },
    { 
      key: "DeathPenalty", 
      label: "Death Penalty", 
      value: options.DeathPenalty || "All",
      type: "select",
      options: [
        { label: "All", value: "All" },
        { label: "Item", value: "Item" },
        { label: "ItemAndEquipment", value: "ItemAndEquipment" },
        { label: "None", value: "None" }
      ]
    },
    { 
      key: "PalEggDefaultHatchingTime", 
      label: "Pal Egg Default Hatching Time", 
      value: options.PalEggDefaultHatchingTime || "72.000000",
      min: 0.1,
      max: 72,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "WorkSpeedRate", 
      label: "Work Speed Rate", 
      value: options.WorkSpeedRate || "1.000000",
      min: 0.1,
      max: 5,
      step: 0.1,
      harder: true,
      type: "slider"
    }
  ];
  
  // Advanced settings
  advancedSettings.value = [
    { 
      key: "bEnablePlayerToPlayerDamage", 
      label: "Enable Player To Player Damage", 
      value: options.bEnablePlayerToPlayerDamage === "True",
      type: "switch"
    },
    { 
      key: "bEn ableFriendlyFire", 
      label: "Enable Friendly Fire", 
      value: options.bEnableFriendlyFire === "True",
      type: "switch"
    },
    { 
      key: "bEnableInvaderEnemy", 
      label: "Enable Invader Enemy", 
      value: options.bEnableInvaderEnemy === "True",
      type: "switch"
    },
    { 
      key: "bActiveUNKO", 
      label: "Active UNKO (Pal will poop)", 
      value: options.bActiveUNKO === "True",
      type: "switch"
    },
    { 
      key: "bEnableAimAssistPad", 
      label: "Enable Aim Assist Controller", 
      value: options.bEnableAimAssistPad === "True",
      type: "switch"
    },
    { 
      key: "bEnableAimAssistKeyboard", 
      label: "Enable Aim Assist Keyboard", 
      value: options.bEnableAimAssistKeyboard === "True",
      type: "switch"
    },
    { 
      key: "DropItemMaxNum", 
      label: "Drop Item Max Num", 
      value: options.DropItemMaxNum || "3000",
      min: 100,
      max: 5000,
      type: "slider"
    },
    { 
      key: "DropItemMaxNum_UNKO", 
      label: "Drop Item Max Num (UNKO)", 
      value: options.DropItemMaxNum_UNKO || "100",
      min: 10,
      max: 1000,
      type: "slider"
    },
    { 
      key: "BaseCampMaxNum", 
      label: "Base Camp Max Num", 
      value: options.BaseCampMaxNum || "128",
      min: 1,
      max: 256,
      type: "slider"
    },
    { 
      key: "BaseCampWorkerMaxNum", 
      label: "Base Camp Worker Max Num", 
      value: options.BaseCampWorkerMaxNum || "15",
      min: 1,
      max: 50,
      type: "slider"
    },
    { 
      key: "DropItemAliveMaxHours", 
      label: "Drop Item Alive Max Hours", 
      value: options.DropItemAliveMaxHours || "1.000000",
      min: 0.1,
      max: 24,
      step: 0.1,
      type: "slider"
    },
    { 
      key: "bAutoResetGuildNoOnlinePlayers", 
      label: "Auto Reset Guild No Online Players", 
      value: options.bAutoResetGuildNoOnlinePlayers === "True",
      type: "switch"
    },
    { 
      key: "AutoResetGuildTimeNoOnlinePlayers", 
      label: "Auto Reset Guild Time No Online Players", 
      value: options.AutoResetGuildTimeNoOnlinePlayers || "72.000000",
      min: 1,
      max: 168,
      step: 1,
      type: "slider"
    },
    { 
      key: "GuildPlayerMaxNum", 
      label: "Guild Player Max Num", 
      value: options.GuildPlayerMaxNum || "20",
      min: 1,
      max: 32,
      type: "slider"
    },
    { 
      key: "bIsMultiplay", 
      label: "Is Multiplay", 
      value: options.bIsMultiplay === "True",
      type: "switch"
    },
    { 
      key: "bIsPvP", 
      label: "Is PvP", 
      value: options.bIsPvP === "True",
      type: "switch"
    },
    { 
      key: "bCanPickupOtherGuildDeathPenaltyDrop", 
      label: "Can Pickup Other Guild Death Penalty Drop", 
      value: options.bCanPickupOtherGuildDeathPenaltyDrop === "True",
      type: "switch"
    },
    { 
      key: "bEnableNonLoginPenalty", 
      label: "Enable Non Login Penalty", 
      value: options.bEnableNonLoginPenalty === "True",
      type: "switch"
    },
    { 
      key: "bEnableFastTravel", 
      label: "Enable Fast Travel", 
      value: options.bEnableFastTravel === "True",
      type: "switch"
    },
    { 
      key: "bIsStartLocationSelectByMap", 
      label: "Is Start Location Select By Map", 
      value: options.bIsStartLocationSelectByMap === "True",
      type: "switch"
    },
    { 
      key: "bExistPlayerAfterLogout", 
      label: "Exist Player After Logout", 
      value: options.bExistPlayerAfterLogout === "True",
      type: "switch"
    },
    { 
      key: "bEnableDefenseOtherGuildPlayer", 
      label: "Enable Defense Other Guild Player", 
      value: options.bEnableDefenseOtherGuildPlayer === "True",
      type: "switch"
    },
    { 
      key: "CoopPlayerMaxNum", 
      label: "Coop Player Max Num", 
      value: options.CoopPlayerMaxNum || "4",
      min: 1,
      max: 8,
      type: "slider"
    }
  ];
};

// Convert settings back to INI format
const generateIniConfig = () => {
  // Combine all settings
  const allSettings = {};
  
  // Process server settings
  serverSettings.value.forEach(setting => {
    if (setting.type === "switch") {
      allSettings[setting.key] = setting.value ? "True" : "False";
    } else {
      allSettings[setting.key] = setting.value.toString();
    }
  });
  
  // Process game settings
  gameSettings.value.forEach(setting => {
    if (setting.type === "switch") {
      allSettings[setting.key] = setting.value ? "True" : "False";
    } else {
      allSettings[setting.key] = setting.value.toString();
    }
  });
  
  // Process advanced settings
  advancedSettings.value.forEach(setting => {
    if (setting.type === "switch") {
      allSettings[setting.key] = setting.value ? "True" : "False";
    } else {
      allSettings[setting.key] = setting.value.toString();
    }
  });
  
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
        <!-- Debug info -->
        <div v-if="serverSettings.length === 0 && gameSettings.length === 0 && advancedSettings.length === 0">
          <p>No settings loaded. Check console for errors.</p>
        </div>
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
            <n-space vertical size="large" style="margin-top: 20px">
              <n-form-item 
                v-for="setting in serverSettings" 
                :key="setting.key" 
                :label="setting.label"
              >
                <!-- Text input -->
                <n-input 
                  v-if="setting.type === 'text'" 
                  v-model:value="setting.value" 
                  :placeholder="setting.label"
                />
                
                <!-- Password input -->
                <n-input 
                  v-else-if="setting.type === 'password'" 
                  v-model:value="setting.value" 
                  type="password"
                  :placeholder="setting.label"
                  show-password-on="click"
                />
                
                <!-- Textarea input -->
                <n-input 
                  v-else-if="setting.type === 'textarea'" 
                  v-model:value="setting.value" 
                  type="textarea"
                  :placeholder="setting.label"
                  :autosize="{ minRows: 3, maxRows: 5 }"
                />
                
                <!-- Number input -->
                <n-input-number 
                  v-else-if="setting.type === 'number'" 
                  v-model:value="setting.value" 
                  :min="setting.min"
                  :max="setting.max"
                />
                
                <!-- Select input -->
                <n-select 
                  v-else-if="setting.type === 'select'" 
                  v-model:value="setting.value" 
                  :options="setting.options"
                />
                
                <!-- Slider -->
                <div v-else-if="setting.type === 'slider'" class="slider-container">
                  <n-input-number 
                    v-model:value="setting.value" 
                    :min="setting.min"
                    :max="setting.max"
                    :step="setting.step || 1"
                    style="width: 100px"
                  />
                  <n-slider 
                    v-model:value="setting.value" 
                    :min="setting.min"
                    :max="setting.max"
                    :step="setting.step || 1"
                    style="margin: 0 10px; flex-grow: 1"
                  />
                  <n-button circle size="small">
                    <template #icon>
                      <n-icon><refresh /></n-icon>
                    </template>
                  </n-button>
                </div>
                
                <!-- Switch -->
                <n-switch v-else-if="setting.type === 'switch'" v-model:value="setting.value" />
              </n-form-item>
            </n-space>
          </n-tab-pane>
          
          <n-tab-pane name="game" tab="In-Game Settings">
            <n-space vertical size="large" style="margin-top: 20px">
              <n-form-item 
                v-for="setting in gameSettings" 
                :key="setting.key" 
                :label="setting.label"
              >
                <!-- Select input -->
                <n-select 
                  v-if="setting.type === 'select'" 
                  v-model:value="setting.value" 
                  :options="setting.options"
                />
                
                <!-- Slider with harder/easier indicators -->
                <div v-else-if="setting.type === 'slider'" class="slider-container">
                  <n-input-number 
                    v-model:value="setting.value" 
                    :min="setting.min"
                    :max="setting.max"
                    :step="setting.step || 1"
                    style="width: 100px"
                  />
                  <div style="display: flex; flex-direction: column; flex-grow: 1; margin: 0 10px;">
                    <div v-if="setting.harder || setting.easier" style="display: flex; justify-content: space-between; margin-bottom: 5px">
                      <span v-if="setting.harder" style="color: #ff4d4f">⬤ Harder</span>
                      <span v-else>&nbsp;</span>
                      <span v-if="setting.easier" style="color: #52c41a">Easier ⬤</span>
                      <span v-else>&nbsp;</span>
                    </div>
                    <n-slider 
                      v-model:value="setting.value" 
                      :min="setting.min"
                      :max="setting.max"
                      :step="setting.step || 1"
                    />
                  </div>
                  <n-button circle size="small">
                    <template #icon>
                      <n-icon><refresh /></n-icon>
                    </template>
                  </n-button>
                </div>
              </n-form-item>
            </n-space>
          </n-tab-pane>
          
          <n-tab-pane name="advanced" tab="Advanced Settings">
            <n-space vertical size="large" style="margin-top: 20px">
              <n-form-item 
                v-for="setting in advancedSettings" 
                :key="setting.key" 
                :label="setting.label"
              >
                <!-- Select input -->
                <n-select 
                  v-if="setting.type === 'select'" 
                  v-model:value="setting.value" 
                  :options="setting.options"
                />
                
                <!-- Switch -->
                <n-switch v-else-if="setting.type === 'switch'" v-model:value="setting.value" />
                
                <!-- Slider -->
                <div v-else-if="setting.type === 'slider'" class="slider-container">
                  <n-input-number 
                    v-model:value="setting.value" 
                    :min="setting.min"
                    :max="setting.max"
                    :step="setting.step || 1"
                    style="width: 100px"
                  />
                  <n-slider 
                    v-model:value="setting.value" 
                    :min="setting.min"
                    :max="setting.max"
                    :step="setting.step || 1"
                    style="margin: 0 10px; flex-grow: 1"
                  />
                  <n-button circle size="small">
                    <template #icon>
                      <n-icon><refresh /></n-icon>
                    </template>
                  </n-button>
                </div>
              </n-form-item>
            </n-space>
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
.slider-container {
  display: flex;
  align-items: center;
  width: 100%;
}
</style>
