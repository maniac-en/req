package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/messages"
	"github.com/maniac-en/req/internal/tabs"
)

type Model struct {
	tabs      []tabs.Tab
	activeTab int
	width     int
	height    int

	// Global state for sharing data
	state *global.State
}

func InitialModel() Model {

	globalState := global.NewGlobalState()

	return Model{
		state: globalState,
		tabs: []tabs.Tab{
			tabs.NewCollectionsTab(globalState),
			tabs.NewAddCollectionTab(),
			tabs.NewEditCollectionTab(),
			tabs.NewEndpointsTab(globalState),
			tabs.NewAddEndpointTab(globalState),
		},
	}

}

func (m Model) Init() tea.Cmd {
	return m.tabs[m.activeTab].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case messages.SwitchTabMsg:
		if msg.TabIndex >= 0 && msg.TabIndex < len(m.tabs) {
			m.activeTab = msg.TabIndex
			return m, m.tabs[m.activeTab].OnFocus()
		}
		return m, nil

	case messages.EditCollectionMsg:
		if editTab, ok := m.tabs[2].(*tabs.EditCollectionTab); ok {
			editTab.SetEditingCollection(msg.Label, msg.Value)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		// removed q here because that was causing issues with input fields
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	default:
		m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	const headerFooterHeight = 1
	const padding = 1
	headerText := m.tabs[m.activeTab].Name()
	instructions := m.tabs[m.activeTab].Instructions()

	headerStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Height(headerFooterHeight).
		Width(len(headerText)+10).
		Align(lipgloss.Center, lipgloss.Top)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("45")).
		Width(m.width-50).
		PaddingBottom(1).
		Height(headerFooterHeight).
		Align(lipgloss.Center, lipgloss.Center)

	renderedHeader := headerStyle.Render(headerText)
	renderedFooter := footerStyle.Render(instructions)

	headerHeight := lipgloss.Height(renderedHeader)
	footerHeight := lipgloss.Height(renderedFooter)

	contentHeight := m.height - headerHeight - footerHeight
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		Align(lipgloss.Center, lipgloss.Center)

	content := m.tabs[m.activeTab].View()

	return lipgloss.JoinVertical(lipgloss.Center, renderedHeader, contentStyle.Render(content), renderedFooter)
}
