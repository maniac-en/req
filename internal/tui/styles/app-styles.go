package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	footerNameStyle    = lipgloss.NewStyle().Bold(true)
	footerNameBGStyle  = lipgloss.NewStyle().Background(footerNameBG).Padding(0, 3, 0)
	FooterSegmentStyle = lipgloss.NewStyle().Background(footerSegmentBG).PaddingLeft(2).Foreground(footerSegmentFG)
	FooterVersionStyle = lipgloss.NewStyle().Background(footerSegmentBG).AlignHorizontal(lipgloss.Right).PaddingRight(2).Foreground(footerSegmentFG)
	TabHeadingInactive = lipgloss.NewStyle().Width(25).AlignHorizontal(lipgloss.Center).Border(lipgloss.NormalBorder(), false, false, false, true)
	TabHeadingActive   = lipgloss.NewStyle().Background(accent).Foreground(headingForeground).Width(25).AlignHorizontal(lipgloss.Center).Border(lipgloss.NormalBorder(), false, false, false, true)
)
