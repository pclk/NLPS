package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Italic(true)

	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))
)

func Success(s string) string {
	return SuccessStyle.Render(s)
}

func Error(s string) string {
	return ErrorStyle.Render(s)
}

func Info(s string) string {
	return InfoStyle.Render(s)
}

func Command(s string) string {
	return CommandStyle.Render(s)
}

func Value(s string) string {
	return ValueStyle.Render(s)
}
