/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pclk/NLPS/pkg/ai"
	"github.com/pclk/NLPS/pkg/powershell"
	"github.com/pclk/NLPS/pkg/utils"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nlps",
	Short: "A CLI tool for generating and executing PowerShell commands using AI",
	Long: `This application uses AI to generate PowerShell commands based on user input.
It allows users to review and execute the generated commands safely.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize OpenAI client
		client := openai.NewClient(viper.GetString("openai"))

		// Initialize PowerShell
		ps, err := powershell.InitPowerShell()
		if err != nil {
			cmd.PrintErrf("Failed to initialize PowerShell: %v\n", err)
			os.Exit(1)
		}
		defer ps.Close()

		reader := bufio.NewReader(os.Stdin)

		userInput := utils.GetAndStripUserInput(reader)
		if userInput == "exit" {
			os.Exit(1)
		}

		if strings.HasPrefix(userInput, "exe ") {
			powershell.ExecutePowerShellCommand(ps, strings.TrimPrefix(userInput, "exe "))
			os.Exit(1)
		}

		handleAIGeneratedCommand(client, ps, reader, userInput)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Use application-specific directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		os.Exit(1)
	}

	appConfigDir := filepath.Join(configDir, "nlps")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(appConfigDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	rootCmd.PersistentFlags().String("openai", "", "OpenAI API Key")
	viper.BindPFlag("openai", rootCmd.PersistentFlags().Lookup("openai"))

	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))

	// Write config if a new value is set
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Flag("openai").Changed {
			viper.Set("openai", cmd.Flag("openai").Value.String())
			err := viper.WriteConfig()
			if err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					// Config file not found; create it.
					err = viper.SafeWriteConfig()
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error writing config file:", err)
				}
			}
		}
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Fprintln(os.Stderr, "To fix error, please set OpenAI API Key using --openai flag.")
		}
	}
}

func handleAIGeneratedCommand(client *openai.Client, ps *powershell.PowerShell, reader *bufio.Reader, userInput string) {
	aiGeneratedCommand := ai.GeneratePowerShellCommand(client, userInput)
	fmt.Printf("Generated command: %s\n", aiGeneratedCommand)

	if utils.UserConfirmsExecution(reader) {
		powershell.ExecutePowerShellCommand(ps, aiGeneratedCommand)
	}
}
