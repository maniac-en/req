package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EndpointsView struct {
	height int
	width  int
	order  int
}

func (e EndpointsView) Init() tea.Cmd {
	return nil
}

func (e EndpointsView) Name() string {
	return "Endpoints"
}

func (e EndpointsView) Help() []key.Binding {
	return []key.Binding{}
}

func (e EndpointsView) GetFooterSegment() string {
	return ""
}

func (e EndpointsView) Update(msg tea.Msg) (ViewInterface, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.height = msg.Height
		e.width = msg.Width
	}
	return e, nil
}

func (e EndpointsView) View() string {
	return lipgloss.NewStyle().Height(e.height).Width(e.width).Align(lipgloss.Center, lipgloss.Center).Render("Endpoints View!")
}

func (e EndpointsView) OnFocus() {

}

func (e EndpointsView) OnBlur() {

}

func (e EndpointsView) Order() int {
	return e.order
}

func NewEndpointsView(order int) *EndpointsView {
	return &EndpointsView{
		order: order,
	}
}
