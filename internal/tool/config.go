package tool

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "github.com/spf13/viper"
)

// GetServerConfig reads the PalWorld server configuration file
func GetServerConfig() (string, error) {
    // Determine the path to the configuration file
    savePath := viper.GetString("save.path")
    
    // The configuration file is typically in a subdirectory of the save path
    var configPath string
    
    // Check if we're on Windows or Linux
    windowsConfigPath := filepath.Join(savePath, "Config", "WindowsServer", "PalWorldSettings.ini")
    linuxConfigPath := filepath.Join(savePath, "Config", "LinuxServer", "PalWorldSettings.ini")
    
    // Try Windows path first
    if _, err := os.Stat(windowsConfigPath); err == nil {
        configPath = windowsConfigPath
    } else if _, err := os.Stat(linuxConfigPath); err == nil {
        // Then try Linux path
        configPath = linuxConfigPath
    } else {
        // If neither exists, return an error
        return "", err
    }
    
    // Read the configuration file
    content, err := ioutil.ReadFile(configPath)
    if err != nil {
        return "", err
    }
    
    return string(content), nil
}

// SaveServerConfig writes to the PalWorld server configuration file
func SaveServerConfig(content string) error {
    // Determine the path to the configuration file
    savePath := viper.GetString("save.path")
    
    // The configuration file is typically in a subdirectory of the save path
    var configPath string
    
    // Check if we're on Windows or Linux
    windowsConfigPath := filepath.Join(savePath, "Config", "WindowsServer", "PalWorldSettings.ini")
    linuxConfigPath := filepath.Join(savePath, "Config", "LinuxServer", "PalWorldSettings.ini")
    
    // Try Windows path first
    if _, err := os.Stat(windowsConfigPath); err == nil {
        configPath = windowsConfigPath
    } else if _, err := os.Stat(linuxConfigPath); err == nil {
        // Then try Linux path
        configPath = linuxConfigPath
    } else {
        // If neither exists, return an error
        return err
    }
    
    // Create a backup of the current configuration
    backupPath := configPath + ".bak"
    currentContent, err := ioutil.ReadFile(configPath)
    if err == nil {
        err = ioutil.WriteFile(backupPath, currentContent, 0644)
        if err != nil {
            return err
        }
    }
    
    // Write the new configuration
    err = ioutil.WriteFile(configPath, []byte(content), 0644)
    if err != nil {
        return err
    }
    
    return nil
}
