package tabs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/messages"
)

type AddEndpointTab struct {
	name        string
	nameInput   textinput.Model
	methodInput textinput.Model
	urlInput    textinput.Model
	focused     bool
}

func NewAddEndpointTab() *AddEndpointTab {
	name := textinput.New()
	name.Placeholder = "Enter your endpoint's name... "
	name.CharLimit = 100
	name.Width = 50

	method := textinput.New()
	method.Placeholder = "Enter your method... "
	method.CharLimit = 100
	method.Width = 50

	url := textinput.New()
	url.Placeholder = "Enter your url..."
	url.CharLimit = 100
	url.Width = 50

	name.Focus()

	return &AddEndpointTab{
		name:        "Add Endpoint",
		nameInput:   name,
		methodInput: method,
		urlInput:    url,
		focused:     true,
	}
}

func (a *AddEndpointTab) Name() string {
	return a.name
}
func (a *AddEndpointTab) Instructions() string {
	return "none"
}
func (a *AddEndpointTab) Init() tea.Cmd {
	return textinput.Blink
}
func (a *AddEndpointTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if a.nameInput.Value() != "" {
				// return a.addCollection(a.nameInput.Value())
			}
		case "esc":
			return a, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 3}
			}
		}
	}

	a.nameInput, _ = a.nameInput.Update(msg)
	return a, nil
}
func (a *AddEndpointTab) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(2)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		MarginBottom(2)

	form := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("Create New Endpoint"),
		inputStyle.Render(a.nameInput.View()),
		inputStyle.Render(a.methodInput.View()),
		inputStyle.Render(a.urlInput.View()),
	)

	containerStyle := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Align(lipgloss.Center, lipgloss.Center)

	return containerStyle.Render(form)
}

func (a *AddEndpointTab) OnFocus() tea.Cmd {
	a.nameInput.Focus()
	a.focused = true
	return textinput.Blink
}

func (a *AddEndpointTab) OnBlur() tea.Cmd {
	a.nameInput.Blur()
	a.focused = false
	return nil
}
