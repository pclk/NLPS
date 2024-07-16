/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nlps",
	Short: "A CLI tool for generating and executing PowerShell commands using AI",
	Long: `This application uses AI to generate PowerShell commands based on user input.
It allows users to review and execute the generated commands safely.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {
	// Use application-specific directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}

	appConfigDir := filepath.Join(configDir, "nlps")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}
	cobra.OnInitialize(initConfig)

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
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}
	appConfigDir := filepath.Join(configDir, "nlps")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(appConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create it
			configFile := filepath.Join(appConfigDir, "config.yaml")
			if err := viper.SafeWriteConfigAs(configFile); err != nil {
				log.Fatalf("Error creating config file: %v", err)
			}
			log.Println("Created new config file:", configFile)
		} else {
			log.Println("Error reading config file:", err)
		}
	}
}
