package styles

import "github.com/charmbracelet/lipgloss"

var (
	FooterNameStyle    = lipgloss.NewStyle().Bold(true)
	FooterNameBGStyle  = lipgloss.NewStyle().Background(FooterNameBG).Padding(0, 3, 0)
	FooterVersionStyle = lipgloss.NewStyle().Background(lipgloss.Color("#262626")).AlignHorizontal(lipgloss.Right).PaddingRight(2).Foreground(lipgloss.Color("#656565"))
)
