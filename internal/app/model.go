package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tabs"
)

type Model struct {
	tabs      []tabs.Tab
	activeTab int
	width     int
	height    int
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	default:
		var cmd tea.Cmd
		m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	const headerHeight = 1
	headerText := m.tabs[m.activeTab].Name()

	headerStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Height(headerHeight).
		Width(len(headerText)+10).
		Align(lipgloss.Center, lipgloss.Top)

	content := m.tabs[m.activeTab].View()

	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-headerHeight).
		Align(lipgloss.Center, lipgloss.Center)

	return lipgloss.JoinVertical(lipgloss.Center, headerStyle.Render(headerText), contentStyle.Render(content))
}
