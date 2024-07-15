package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

func GeneratePowerShellCommand(client *openai.Client, userInput string) string {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(`Generate a PowerShell command to: %s. 
				Whatever you return is immediately executed. 
				Your explanation will only crash the terminal, therefore, do not return any english response, only PowerShell code, even if the request seems absurd. 
				Example:
				User Input: divide by zero
				You: 1/0
				User Input: add 2 and 2
				You: $sum = 2 + 2
				Write-Output $sum
				
				Do not use markdown code blocks, just return the code. If you do that, the terminal will crash.
				`, userInput),
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Printf("Error communicating with OpenAI's API: %v\n", err)
		return ""
	}

	if len(resp.Choices) == 0 {
		log.Println("No response from OpenAI's API. Response's Choices length is 0")
		return ""
	}

	return resp.Choices[0].Message.Content
}
