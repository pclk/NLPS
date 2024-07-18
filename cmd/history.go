package cmd

import (
	"fmt"
	"time"

	"github.com/pclk/NLPS/internal/history"
	"github.com/pclk/NLPS/internal/ui"
	"github.com/spf13/cobra"
)

var clearHistory bool

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Display and manage previously executed commands",
	Long:  `This command allows you to view and manage the history of executed commands.`,
	Run:   runHistory,
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().BoolVarP(&clearHistory, "clear", "c", false, "Clear the command history")
}

func runHistory(cmd *cobra.Command, args []string) {
	if clearHistory {
		history.ClearHistory()
		fmt.Println(ui.Success("Command history cleared."))
		return
	}

	commands := history.GetHistory()
	if len(commands) == 0 {
		fmt.Println(ui.Info("No command history found."))
		return
	}

	fmt.Println(ui.Info("Command History:"))
	for i, command := range commands {
		fmt.Printf("%d. %s at %s\n", i+1, ui.Command(command.Command), command.Time.Format(time.RFC3339))
	}
}
