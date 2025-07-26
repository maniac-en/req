package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tabs"
)

type Model struct {
	tabs      []tabs.Tab
	activeTab int
}

func InitialModel() Model {
	tabList := []tabs.Tab{
		tabs.NewCollectionsTab(),
	}

	return Model{
		tabs:      tabList,
		activeTab: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.tabs[m.activeTab].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
			return m, cmd
		}
	default:
		var cmd tea.Cmd
		m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	var tabHeaders strings.Builder
	for i, tab := range m.tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.activeTab {
			style = style.Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
		} else {
			style = style.Background(lipgloss.Color("240")).Foreground(lipgloss.Color("255"))
		}
		tabHeaders.WriteString(style.Render(tab.Name()))
	}

	content := m.tabs[m.activeTab].View()

	return fmt.Sprintf("%s\n\n%s",
		tabHeaders.String(),
		content,
	)
}
