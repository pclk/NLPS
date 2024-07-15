package test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/pclk/NLPS/pkg/utils"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Zero length", 0},
		{"Positive length", 10},
		{"Large length", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.GenerateRandomString(tt.length)
			if len(got) != tt.length {
				t.Errorf("GenerateRandomString(%d) = %v, want length %v", tt.length, got, tt.length)
			}
		})
	}
}

func TestGetAndStripUserInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Simple input", "test input\n", "test input"},
		{"Input with spaces", "  test input  \n", "test input"},
		{"Empty input", "\n", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got := utils.GetAndStripUserInput(reader)
			if got != tt.want {
				t.Errorf("GetAndStripUserInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserConfirmsExecution(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Confirm with Y", "Y\n", true},
		{"Confirm with y", "y\n", true},
		{"Confirm with empty", "\n", true},
		{"Deny with N", "N\n", false},
		{"Deny with n", "n\n", false},
		{"Deny with other input", "no\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got := utils.UserConfirmsExecution(reader)
			if got != tt.want {
				t.Errorf("UserConfirmsExecution() = %v, want %v", got, tt.want)
			}
		})
	}
}
