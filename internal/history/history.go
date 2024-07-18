package history

import (
	"time"

	"github.com/spf13/viper"
)

const maxHistorySize = 100

type HistoryEntry struct {
	Command string    `json:"string"`
	Time    time.Time `json:"time"`
}

type History []HistoryEntry

func AddToHistory(command string) {
	history := GetHistory()
	newEntry := HistoryEntry{
		Command: command,
		Time:    time.Now(),
	}
	history = append([]HistoryEntry{newEntry}, history...)
	if len(history) > maxHistorySize {
		history = history[:maxHistorySize]
	}
	viper.Set("history", history)
	viper.WriteConfig()
}

func GetHistory() History {
	var history History
	err := viper.UnmarshalKey("history", &history)
	if err != nil {
		return History{}
	}
	return history
}

func ClearHistory() {
	viper.Set("history", []string{})
	viper.WriteConfig()
}
