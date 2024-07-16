package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pclk/NLPS/internal/ui"
	"github.com/spf13/viper"
)

func GetAllConfig() map[string]string {
	config := viper.AllSettings()
	configMap := make(map[string]string)
	for k, v := range config {
		configMap[k] = fmt.Sprintf("%v", v)
	}
	return configMap
}

func GetConfig(key string) string {
	value := viper.GetString(key)
	if value != "" {
		fmt.Printf("%s: %s\n", ui.Command(key), ui.Value(fmt.Sprintf("%v", value)))
	} else {
		fmt.Println(ui.Error(fmt.Sprintf("Configuration '%s' not found.", key)))
	}
	return value
}

func SetConfig(key, value string) {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return
	}
}

func RemoveConfig() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}
	appConfigDir := filepath.Join(configDir, "nlps")
	log.Println(ui.Info("Removing configuration file:"), appConfigDir)
	os.Remove(fmt.Sprintf("%s/config.yaml", appConfigDir))
	err = os.Remove(appConfigDir)
	if err != nil {
		log.Fatalf(ui.Error("Error removing configuration file: %v"), err)
	}
	fmt.Println(ui.Success("Configuration file removed."))
}

func AddAlias(name, command string) error {
	aliases := viper.GetStringMapString("aliases")
	if aliases == nil {
		aliases = make(map[string]string)
	}
	aliases[name] = command
	viper.Set("aliases", aliases)
	return viper.WriteConfig()
}

func RemoveAlias(name string) bool {
	aliases := viper.GetStringMapString("aliases")
	if aliases == nil {
		return false
	}
	if _, exists := aliases[name]; exists {
		delete(aliases, name)
		viper.Set("aliases", aliases)
		viper.WriteConfig()
		return true
	}
	return false
}

func GetAliases() map[string]string {
	return viper.GetStringMapString("aliases")
}
