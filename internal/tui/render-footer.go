package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) renderFooter(content string) string {
	footerStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Padding(1, 2)

	return footerStyle.
		Width(m.Width).
		Render(content)
}
