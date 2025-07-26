package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) renderHeader(content string) string {
	headerStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Align(lipgloss.Center)

	return headerStyle.
		Width(m.Width).
		Render(content)
}
