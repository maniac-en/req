package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	FooterNameStyle    = lipgloss.NewStyle().Bold(true)
	FooterNameBGStyle  = lipgloss.NewStyle().Background(FooterNameBG).Padding(0, 3, 0)
	FooterSegmentStyle = lipgloss.NewStyle().Background(lipgloss.Color("#262626")).PaddingLeft(2).Foreground(lipgloss.Color("#656565"))
	FooterVersionStyle = lipgloss.NewStyle().Background(lipgloss.Color("#262626")).AlignHorizontal(lipgloss.Right).PaddingRight(2).Foreground(lipgloss.Color("#656565"))
	TabHeadingInactive = lipgloss.NewStyle().Width(25).AlignHorizontal(lipgloss.Center).Border(lipgloss.NormalBorder(), false, false, false, true)
	TabHeadingActive   = lipgloss.NewStyle().Background(Accent).Foreground(HeadingForeground).Width(25).AlignHorizontal(lipgloss.Center).Border(lipgloss.NormalBorder(), false, false, false, true)
)
