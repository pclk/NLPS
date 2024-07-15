package main

import (
	"bufio"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int
	}{
		{"Zero length", 0, 0},
		{"Positive length", 10, 10},
		{"Large length", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRandomString(tt.length)
			if len(got) != tt.want {
				t.Errorf("generateRandomString(%d) = %v, want length %v", tt.length, got, tt.want)
			}
		})
	}
}

func TestInitPowerShell(t *testing.T) {
	ps, err := initPowerShell()
	if err != nil {
		t.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	if ps == nil {
		t.Fatal("initPowerShell() returned nil")
	}
	defer ps.close()

	if ps.cmd == nil {
		t.Error("PowerShell cmd is nil")
	}
	if ps.stdin == nil {
		t.Error("PowerShell stdin is nil")
	}
	if ps.stdout == nil {
		t.Error("PowerShell stdout is nil")
	}
}

func TestWrapCommandWithErrorHandling(t *testing.T) {
	ps := &PowerShell{}
	command := "Get-Process"
	marker := "END_MARKER"

	wrapped := ps.wrapCommandWithErrorHandling(command, marker)

	if !strings.Contains(wrapped, command) {
		t.Errorf("Wrapped command does not contain original command")
	}
	if !strings.Contains(wrapped, marker) {
		t.Errorf("Wrapped command does not contain marker")
	}
	if !strings.Contains(wrapped, "try") || !strings.Contains(wrapped, "catch") {
		t.Errorf("Wrapped command does not contain error handling")
	}
}

func TestSendCommand(t *testing.T) {
	ps, err := initPowerShell()
	if err != nil {
		t.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	defer ps.close()

	tests := []struct {
		name           string
		command        string
		expectedOutput string
	}{
		{"Echo command", "echo 'Hello, World!'", "Hello, World!"},
		{"Addition", "$sum = 2 + 2; $sum", "4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, output := ps.sendCommand(tt.command)
			t.Logf("Command: %s", command)

			// Remove PowerShell formatting and empty lines
			output = cleanPowerShellOutput(output)
			t.Logf("Got cleaned: %s", output)

			if tt.expectedOutput != "" && !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("sendCommand(%s) = %v, want %v", tt.command, output, tt.expectedOutput)
			}
			if output == "" {
				t.Errorf("sendCommand(%s) returned empty output", tt.command)
			}

			switch tt.name {
			case "Get current directory":
				if !filepath.IsAbs(output) {
					t.Errorf("sendCommand(%s) returned non-absolute path: %s", tt.command, output)
				}
			case "Addition":
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("sendCommand(%s) = %v, want %v", tt.command, output, tt.expectedOutput)
				}
			}
		})
	}
}

func TestSendCommandError(t *testing.T) {
	ps, err := initPowerShell()
	if err != nil {
		t.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	defer ps.close()

	command := "1/0"
	command, output := ps.sendCommand(command)
	t.Logf("Command: %s", command)
	t.Logf("Got: %s", output)

	// Clean the output
	cleanedOutput := cleanPowerShellOutput(output)

	expectedErrorMessage := "Attempted to divide by zero."
	if !strings.Contains(cleanedOutput, expectedErrorMessage) {
		t.Errorf("Expected error message '%s' not found in output: %s", expectedErrorMessage, cleanedOutput)
	}

	// Check for additional error information
	expectedErrorInfo := []string{
		"RuntimeException",
		"CategoryInfo",
		"FullyQualifiedErrorId",
	}

	for _, info := range expectedErrorInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected error information '%s' not found in output: %s", info, output)
		}
	}
}

func TestReadOutputUntilMarker(t *testing.T) {
	ps, err := initPowerShell()
	if err != nil {
		t.Fatalf("Failed to initialize PowerShell: %v", err)
	}
	defer ps.close()

	exampleStdout := `
4
END_OF_COMMAND_IaBHvxXCSi
`
	ps.stdout = bufio.NewReader(strings.NewReader(exampleStdout))

	result := ps.readOutputUntilMarker("END_OF_COMMAND_IaBHvxXCSi")
	result = cleanPowerShellOutput(result)

	expectedResult := "4"
	if result != expectedResult {
		t.Errorf("Expected result '%s', but got '%s'", expectedResult, result)
	}
}

func cleanPowerShellOutput(output string) string {
	lines := strings.Split(output, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && !strings.HasPrefix(trimmedLine, ">>") {
			cleanedLines = append(cleanedLines, trimmedLine)
		}
	}
	return strings.Join(cleanedLines, "\n")
}
