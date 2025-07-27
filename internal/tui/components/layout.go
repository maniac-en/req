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
	header := l.Header(title)
	footer := l.Footer(instructions)

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)

	contentArea := l.Content(content, headerHeight, footerHeight)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		contentArea,
		footer,
	)
}
