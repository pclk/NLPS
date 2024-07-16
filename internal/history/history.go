package history

import (
	"github.com/spf13/viper"
)

const maxHistorySize = 100

func AddToHistory(command string) {
	history := GetHistory()
	history = append([]string{command}, history...)
	if len(history) > maxHistorySize {
		history = history[:maxHistorySize]
	}
	viper.Set("history", history)
	viper.WriteConfig()
}

func GetHistory() []string {
	return viper.GetStringSlice("history")
}

func ClearHistory() {
	viper.Set("history", []string{})
	viper.WriteConfig()
}
