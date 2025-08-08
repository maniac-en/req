package styles

import "github.com/charmbracelet/lipgloss"

var (
	SelectedListStyle = lipgloss.NewStyle().Foreground(accent).PaddingLeft(1).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(accent)
)
