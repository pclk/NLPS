package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/pclk/NLPS/internal/ai"
	"github.com/pclk/NLPS/internal/history"
	"github.com/pclk/NLPS/internal/powershell"
	"github.com/pclk/NLPS/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	silent    bool
	noOutput  bool
	errorOnly bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Generate and run a PowerShell command",
	Long:  `This command generates and runs a PowerShell command based on your input.`,
	Run:   runCommand,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Run the command without asking for confirmation")
	runCmd.Flags().BoolVarP(&noOutput, "no-output", "n", false, "Do not display the command output and errors")
	runCmd.Flags().BoolVarP(&errorOnly, "error-only", "e", false, "Only display the command error")
}

func runCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(ui.Error("Please provide a description of the command you want to run."))
		return
	}

	client := openai.NewClient(viper.GetString("openai"))
	ps, err := powershell.InitPowerShell()
	if err != nil {
		log.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	defer ps.Close()

	userInput := strings.Join(args, " ")
	aiGeneratedCommand := ai.GeneratePowerShellCommand(client, userInput)

	if aiGeneratedCommand == "" {
		fmt.Println(ui.Error("Failed to generate a command. Please try again."))
		return
	}

	fmt.Printf(ui.Command("Generated command: %s\n"), aiGeneratedCommand)

	if !silent {
		fmt.Print(ui.Info("Do you want to execute this command? (y/n): "))
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println(ui.Info("Command execution cancelled."))
			return
		}
	}

	output, errorOutput := ps.SendCommand(aiGeneratedCommand)
	history.AddToHistory(aiGeneratedCommand)

	if !noOutput {
		if !errorOnly {
			fmt.Println(ui.Success("Command output:"))
			fmt.Println(ui.Success(output))
		}
		if errorOutput != "" {
			fmt.Println(ui.Error("Error output:"))
			fmt.Println(ui.Error(errorOutput))
		}
	}
}
