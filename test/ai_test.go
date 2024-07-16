package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pclk/NLPS/internal/ai"
	"github.com/sashabaranov/go-openai"
)

func TestGeneratePowerShellCommand(t *testing.T) {
	godotenv.Load("../.env.local")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient(apiKey)
	userInput := "Print 'Hello, World!'"

	result := ai.GeneratePowerShellCommand(client, userInput)

	if result == "" {
		t.Errorf("GeneratePowerShellCommand() returned empty string, expected non-empty result")
	}
}
