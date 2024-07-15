package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func GenerateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[int(rand.Int63()%int64(len(letters)))]
	}
	return string(b)
}

func GetAndStripUserInput(reader *bufio.Reader) string {
	fmt.Print("Enter a description for the PowerShell command you want to generate (or 'exit' to quit): ")
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading user input:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(userInput)
}

func UserConfirmsExecution(reader *bufio.Reader) bool {
	fmt.Print("Do you want to execute this command? (Y/n): ")
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading user input:", err)
		os.Exit(1)
	}
	userInput = strings.TrimSpace(userInput)
	if userInput == "Y" || userInput == "y" || userInput == "" {
		return true
	}
	return false
}
