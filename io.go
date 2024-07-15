package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func generateRandomString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(randomBytes)
}

func getAndStripUserInput(reader *bufio.Reader) string {
	fmt.Print("Generate a PowerShell command to: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func userConfirmsExecution(reader *bufio.Reader) bool {
	fmt.Print("Execute this command? (Y/n): ")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading user confirmation: %v", err)
		return false
	}
	confirm = strings.TrimSpace(confirm)
	// Return true if the user pressed Enter (empty string) or typed 'y' or 'Y'
	return confirm == "" || strings.ToLower(confirm) == "y"
}

func executePowerShellCommand(ps *PowerShell, command string) {
	command, output := ps.sendCommand(command)
	fmt.Println("Output:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(output)
	fmt.Println(strings.Repeat("-", 40))
}
