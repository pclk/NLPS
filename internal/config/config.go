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
