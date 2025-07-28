package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/styles"
)

type Layout struct {
	width  int
	height int
}

func NewLayout() Layout {
	return Layout{}
}

func (l *Layout) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l Layout) FullView(title, content, instructions string) string {
	if l.width < 20 || l.height < 10 {
		return content
	}

	windowWidth := int(float64(l.width) * 0.85)
	windowHeight := int(float64(l.height) * 0.8)

	if windowWidth < 50 {
		windowWidth = 50
	}
	if windowHeight < 15 {
		windowHeight = 15
	}

	innerWidth := windowWidth - 4
	innerHeight := windowHeight - 4
	header := styles.WindowHeaderStyle.Copy().
		Width(innerWidth).
		Render(title)

	headerHeight := lipgloss.Height(header)
	contentHeight := innerHeight - headerHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	contentArea := styles.WindowContentStyle.Copy().
		Width(innerWidth).
		Height(contentHeight).
		Render(content)

	windowContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		contentArea,
	)

	borderedWindow := styles.WindowBorderStyle.Copy().
		Width(windowWidth).
		Height(windowHeight).
		Render(windowContent)

	brandingText := "Req - Test APIs with Terminal Velocity"
	appBranding := styles.AppBrandingStyle.Copy().
		Width(l.width).
		Render(brandingText)

	footer := styles.WindowFooterStyle.Copy().
		Width(l.width).
		Render(instructions)

	brandingHeight := lipgloss.Height(appBranding)
	footerHeight := lipgloss.Height(footer)
	windowPlacementHeight := l.height - brandingHeight - footerHeight - 4

	centeredWindow := lipgloss.Place(
		l.width, windowPlacementHeight,
		lipgloss.Center, lipgloss.Center,
		borderedWindow,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		appBranding,
		"",
		centeredWindow,
		"",
		footer,
	)
}
