package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pclk/NLPS/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	table         table.Model
	textInput     textinput.Model
	configs       map[string]string
	cursor        int
	editing       bool
	help          help.Model
	tableKeys     tableKeyMap
	textInputKeys textInputKeyMap
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("24"))

type tableKeyMap struct {
	Up   key.Binding
	Down key.Binding
	Edit key.Binding
	Quit key.Binding
}

type textInputKeyMap struct {
	Paste key.Binding
	Save  key.Binding
	Quit  key.Binding
}

func (tk tableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		tk.Up,
		tk.Down,
		tk.Edit,
		tk.Quit,
	}
}

func (tk tableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}, {}}
}

func (tik textInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		tik.Paste,
		tik.Save,
		tik.Quit,
	}
}

func (tik textInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}, {}}
}

var tableKeys = tableKeyMap{
	Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑", "move up")),
	Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓", "down")),
	Edit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter:", "edit text")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q:", "quit program")),
}

var textInputKeys = textInputKeyMap{
	Paste: key.NewBinding(key.WithKeys("just for help"), key.WithHelp("right-click:", "paste from clipboard")),
	Save:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter:", "save text")),
	Quit:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc:", "quit textinput")),
}

func initialModel() model {
	configs := config.GetAllConfig()

	ti := textinput.New()
	ti.Placeholder = "Enter new value"
	ti.Focus()

	return model{
		table:         initialTable(configs),
		textInput:     ti,
		configs:       configs,
		cursor:        0,
		editing:       false,
		help:          help.New(),
		tableKeys:     tableKeys,
		textInputKeys: textInputKeys,
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
		if k == "history" || k == "alias" {
			continue
		}
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
		BorderForeground(lipgloss.Color("24")).
		BorderBottom(true).
		Align(lipgloss.Center).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("24"))
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
	case tea.WindowSizeMsg:
		m.help.Width, m.textInput.Width = msg.Width, msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tableKeys.Quit):
			if !m.editing {
				return m, tea.Quit
			}
		case key.Matches(msg, m.tableKeys.Up, m.tableKeys.Down):
			if !m.editing {
				m.table, cmd = m.table.Update(msg)
			}
		case key.Matches(msg, m.tableKeys.Edit):
			if !m.editing {
				m.editing = true
				m.textInput.SetValue(m.table.SelectedRow()[1])
				m.textInput.Focus()
			} else {
				key := m.table.SelectedRow()[0]
				value := m.textInput.Value()
				m.configs[key] = value
				viper.Set(key, value)
				viper.WriteConfig()
				m.table = initialTable(m.configs)
				m.editing = false
				m.textInput.Blur()
			}
		case key.Matches(msg, m.textInputKeys.Quit):
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
		b.WriteString(m.help.View(m.textInputKeys))
	} else {
		b.WriteString(m.help.View(m.tableKeys))
	}

	return b.String()
}

func init() {
	rootCmd.AddCommand(configCmd)
}
