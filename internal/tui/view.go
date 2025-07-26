package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {

	header := m.renderHeader(m.Tabs[m.ActiveTab].Title)
	footer := m.renderFooter("Some keybinds")

	availableHeight := m.Height - lipgloss.Height(header) - lipgloss.Height(footer)

	content := m.Tabs[m.ActiveTab].Content(availableHeight, m.Width)

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}
