package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pclk/NLPS/internal/ai"
	"github.com/pclk/NLPS/internal/history"
	"github.com/pclk/NLPS/internal/powershell"
	"github.com/pclk/NLPS/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Interactively chat, generate and run PowerShell commands",
	Long:  `This command allows you to have an interactive chat session to generate and run PowerShell commands.`,
	Run:   runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

func runChat(cmd *cobra.Command, args []string) {
	client := openai.NewClient(viper.GetString("openai"))
	ps, err := powershell.InitPowerShell()
	if err != nil {
		log.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	defer ps.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(ui.Info("Enter your request (or 'exit' to quit): "))
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if userInput == "exit" {
			break
		}

		aiGeneratedCommand := ai.GeneratePowerShellCommand(client, userInput)
		if aiGeneratedCommand == "" {
			fmt.Println(ui.Error("Failed to generate a command. Please try again."))
			continue
		}

		fmt.Printf(ui.Command("Generated command: %s\n"), aiGeneratedCommand)
		fmt.Print(ui.Info("Do you want to execute this command? (y/n): "))
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			output, errorOutput := ps.SendCommand(aiGeneratedCommand)
			fmt.Println(ui.Success("Command output:"))
			fmt.Println(ui.Success(output))
			if errorOutput != "" {
				fmt.Println(ui.Error("Error output:"))
				fmt.Println(ui.Error(errorOutput))
			}
			history.AddToHistory(aiGeneratedCommand)
		} else {
			fmt.Println(ui.Info("Command execution cancelled."))
		}
	}
}
