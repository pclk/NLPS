package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type PowerShell struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Reader
}

func initPowerShell() (*PowerShell, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start PowerShell: %v", err)
	}

	ps := &PowerShell{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdin),
		stdout: bufio.NewReader(stdout),
	}
	ps.sendCommand("$OutputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8; Clear-Host")

	return ps, nil
}

func (ps *PowerShell) sendCommand(command string) (string, string) {
	uniqueMarker := "END_OF_COMMAND_" + generateRandomString(10)
	wrappedCommand := ps.wrapCommandWithErrorHandling(command, uniqueMarker)

	_, err := ps.stdin.WriteString(wrappedCommand + "\n")
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return "", ""
	}
	err = ps.stdin.Flush()
	if err != nil {
		log.Printf("Error sending command while flushing: %v", err)
		return "", ""
	}

	output := ps.readOutputUntilMarker(uniqueMarker)
	if output == "" {
		log.Printf("Warning: No output received for command: %s", command)
	}
	return command, output
}

func (ps *PowerShell) wrapCommandWithErrorHandling(command, marker string) string {
	return fmt.Sprintf(`
try {
	$ErrorActionPreference = "Stop"
	$result = @(
		%s
	)
	$result | ForEach-Object { Write-Output $_ }
} catch {
	$errorMessage = "Error: " + $_.Exception.Message
	$errorMessage += "\\nCategoryInfo: " + $_.CategoryInfo
	$errorMessage += "\\nFullyQualifiedErrorId: " + $_.FullyQualifiedErrorId
	Write-Output $errorMessage
} finally {
	Write-Output ""
	Write-Output '%s'
}
`, command, marker)
}

func (ps *PowerShell) readOutputUntilMarker(marker string) string {
	var output strings.Builder
	var foundResult bool

	for {
		line, err := ps.stdout.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading output: %v", err)
			}
			break
		}

		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == marker {
			break
		}

		if foundResult {
			output.WriteString(trimmedLine)
		} else if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "PS ") && !strings.HasPrefix(trimmedLine, ">>") {
			foundResult = true
			output.WriteString(trimmedLine)
		}

		if err == io.EOF {
			log.Printf("Warning: Marker '%s' not found in the output", marker)
			break
		}
	}

	return strings.TrimSpace(output.String())
}

func (ps *PowerShell) close() {
	if err := ps.cmd.Process.Kill(); err != nil {
		log.Fatal(err)
	}
}