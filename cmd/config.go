package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pclk/NLPS/internal/config"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `View and modify configuration settings for the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := initialModel()
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
		}
	},
}

type model struct {
	table     table.Model
	textInput textinput.Model
	configs   map[string]string
	cursor    int
	editing   bool
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initialModel() model {
	configs := config.GetAllConfig()

	ti := textinput.New()
	ti.Placeholder = "Enter new value"
	ti.Focus()

	return model{
		table:     initialTable(configs),
		textInput: ti,
		configs:   configs,
		cursor:    0,
		editing:   false,
	}
}

func initialTable(configs map[string]string) table.Model {
	columns := []table.Column{
		{Title: "Key", Width: 20},
		{Title: "Value", Width: 40},
	}
	// Sort the keys to maintain a consistent order
	keys := make([]string, 0, len(configs))
	for k := range configs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]table.Row, 0, len(configs))
	for _, k := range keys {
		rows = append(rows, table.Row{k, fmt.Sprintf("%v", configs[k])})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Align(lipgloss.Center).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))
	s.Cell = s.Cell.
		Align(lipgloss.Center).
		Bold(false)
	t.SetStyles(s)

	return t
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "down":
			if !m.editing {
				m.table, cmd = m.table.Update(msg)
			}
		case "enter":
			if !m.editing {
				m.editing = true
				m.textInput.SetValue(m.table.SelectedRow()[1])
				m.textInput.Focus()
			} else {
				key := m.table.SelectedRow()[0]
				value := m.textInput.Value()
				m.configs[key] = value
				config.SetConfig(key, value)
				m.table = initialTable(m.configs)
				m.editing = false
				m.textInput.Blur()
			}
		case "esc":
			if m.editing {
				m.editing = false
				m.textInput.Blur()
			}
		}
	}

	if m.editing {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("> nlps config\n")

	b.WriteString(baseStyle.Render(m.table.View()) + "\n\n")

	if m.editing {
		b.WriteString("Edit value: " + m.textInput.View() + "\n")
	} else {
		b.WriteString("Press 'enter' to edit, 'q' to quit\n")
	}

	return b.String()
}

func init() {
	rootCmd.AddCommand(configCmd)
}
