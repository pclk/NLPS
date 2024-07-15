package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	// godotenv is used to load environment variables from a file
	"github.com/joho/godotenv"
	// go-openai is a client library for interacting with the OpenAI API
	"github.com/sashabaranov/go-openai"
)

func main() {
	// Load environment variables and initialize OpenAI client
	if err := godotenv.Load("../.env.local"); err != nil {
		log.Fatalf("Error loading .env.local file: %v", err)
	}
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Initialize PowerShell
	ps, err := initPowerShell()
	if err != nil {
		log.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	defer ps.close()

	reader := bufio.NewReader(os.Stdin)

	for {
		userInput := getAndStripUserInput(reader)
		if userInput == "exit" {
			break
		}

		if strings.HasPrefix(userInput, "exe ") {
			executePowerShellCommand(ps, strings.TrimPrefix(userInput, "exe "))
			continue
		}

		handleAIGeneratedCommand(client, ps, reader, userInput)
	}
}

func handleAIGeneratedCommand(client *openai.Client, ps *PowerShell, reader *bufio.Reader, userInput string) {
	aiGeneratedCommand := AIGeneratesPowerShellCommand(client, userInput)
	if aiGeneratedCommand == "" {
		return
	}

	if !userConfirmsExecution(reader) {
		return
	}

	executePowerShellCommand(ps, aiGeneratedCommand)
}
