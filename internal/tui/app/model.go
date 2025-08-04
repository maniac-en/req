package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/styles"
)

type AppModel struct {
	width  int
	height int
}

func (a AppModel) Init() tea.Cmd {
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a AppModel) View() string {
	footer := a.Footer()
	header := a.Header()
	availableHeight := a.height - lipgloss.Height(header) - lipgloss.Height(footer)
	view := lipgloss.NewStyle().Height(availableHeight).Width(a.width).Align(lipgloss.Center, lipgloss.Center).Render("Hello world!")
	return lipgloss.JoinVertical(lipgloss.Top, header, view, footer)
}

func (a AppModel) Header() string {
	return "Hello World"
}

func (a AppModel) Footer() string {
	name := styles.GradientText("REQ", styles.FooterNameFGFrom, styles.FooterNameFGTo, styles.FooterNameStyle, styles.FooterNameBGStyle)
	version := styles.FooterVersionStyle.Width(a.width - lipgloss.Width(name)).Render("v0.1.0-alpha2")
	return lipgloss.JoinHorizontal(lipgloss.Left, name, version)
}

func NewAppModel() AppModel {
	return AppModel{}
}
