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

func (l Layout) Header(title string) string {
	return styles.HeaderStyle.
		Width(l.width).
		Render(title)
}

func (l Layout) Footer(instructions string) string {
	return styles.FooterStyle.
		Width(l.width).
		Render(instructions)
}

func (l Layout) Content(content string, headerHeight, footerHeight int) string {
	contentHeight := l.height - headerHeight - footerHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	return styles.ContentStyle.
		Width(l.width).
		Height(contentHeight).
		Render(content)
}

func (l Layout) FullView(title, content, instructions string) string {
	if l.width < 20 || l.height < 10 {
		return content
	}

	// Calculate window dimensions (85% of terminal width, 80% height)
	windowWidth := int(float64(l.width) * 0.85)
	windowHeight := int(float64(l.height) * 0.8)

	// Ensure minimum dimensions
	if windowWidth < 50 {
		windowWidth = 50
	}
	if windowHeight < 15 {
		windowHeight = 15
	}

	// Calculate inner content dimensions (accounting for border)
	innerWidth := windowWidth - 4 // 2 chars for border + padding
	innerHeight := windowHeight - 4

	// Create header and content with simplified, consistent styling
	header := lipgloss.NewStyle().
		Width(innerWidth).
		Padding(1, 2).
		Background(styles.Primary).
		Foreground(styles.TextPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Render(title)

	headerHeight := lipgloss.Height(header)
	contentHeight := innerHeight - headerHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	contentArea := lipgloss.NewStyle().
		Width(innerWidth).
		Height(contentHeight).
		Padding(1, 2).
		Render(content)

	// Join header and content vertically (no footer)
	windowContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		contentArea,
	)

	// Create bordered window
	borderedWindow := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("15")). // White border
		Width(windowWidth).
		Height(windowHeight).
		Render(windowContent)

	// Create elegant app branding banner at top
	brandingText := "Req - Test APIs with Terminal Velocity"
	appBranding := lipgloss.NewStyle().
		Width(l.width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("230")). // Soft cream
		// Background(lipgloss.Color("237")). // Dark gray background
		Bold(true).
		Padding(1, 4).
		Margin(1, 0).
		Render(brandingText)

	// Create footer outside the window
	footer := lipgloss.NewStyle().
		Width(l.width).
		Padding(0, 2).
		Foreground(styles.TextSecondary).
		Align(lipgloss.Center).
		Render(instructions)

	// Calculate vertical position accounting for branding and footer
	brandingHeight := lipgloss.Height(appBranding)
	footerHeight := lipgloss.Height(footer)
	windowPlacementHeight := l.height - brandingHeight - footerHeight - 4 // Extra padding

	centeredWindow := lipgloss.Place(
		l.width, windowPlacementHeight,
		lipgloss.Center, lipgloss.Center,
		borderedWindow,
	)

	// Combine branding, centered window, and footer with proper spacing
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"", // Top padding
		appBranding,
		"", // Extra spacing line
		centeredWindow,
		"", // Reduced spacing before footer
		footer,
	)
}
