package tui

import "github.com/charmbracelet/lipgloss"
import "github.com/charmbracelet/lipgloss/table"

func renderHome(contentHeight, contentWidth int) string {
	const logo = `
	▗▄▄▖ ▗▄▄▄▖▗▄▄▄▖
	▐▌ ▐▌▐▌   ▐▌ ▐▌
	▐▛▀▚▖▐▛▀▀▘▐▌ ▐▌
	▐▌ ▐▌▐▙▄▄▖▐▙▄▟▙▖
	`
	logoHeight := contentHeight / 2
	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle()
			if col == 0 {
				style = style.
					PaddingRight(5)
			}
			if row == 0 {
				style = style.
					PaddingBottom(1)
			}
			return style
		}).
		Height(contentHeight-logoHeight).
		Rows(
			[]string{"Collections", "c"},
			[]string{"Environment", "n"},
		).
		Border(lipgloss.Border{})

	tableStyled := lipgloss.NewStyle().
		Height(contentHeight - logoHeight).
		Width(contentWidth).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Center).
		Render(t.Render())

	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(logoHeight).
		PaddingTop(5).
		PaddingRight(5).
		Align(lipgloss.Center).
		Render(logo)

	return lipgloss.JoinVertical(lipgloss.Left, contentStyle, tableStyled)
}
