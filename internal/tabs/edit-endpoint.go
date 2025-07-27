package tabs

import (
	"context"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/endpoints"
	"github.com/maniac-en/req/internal/messages"
)

type EditEndpointTab struct {
	name         string
	inputs       []textinput.Model
	endpointID   string
	focusedInput int
	focused      bool
	state        *global.State
}

func NewEditEndpointTab(globalState *global.State) *EditEndpointTab {
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

	return &EditEndpointTab{
		name: "Add Endpoint",
		inputs: []textinput.Model{
			name,
			method,
			url,
		},
		focusedInput: 0,
		focused:      true,
		state:        globalState,
	}
}

func (e *EditEndpointTab) Name() string {
	return e.name
}
func (e *EditEndpointTab) Instructions() string {
	return "None"
}
func (e *EditEndpointTab) Init() tea.Cmd {
	return textinput.Blink
}
func (e *EditEndpointTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if e.inputs[0].Value() != "" && e.inputs[1].Value() != "" && e.inputs[2].Value() != "" {
				return e.updateEndpoint(e.inputs[0].Value(), e.inputs[1].Value(), e.inputs[2].Value())
			}
			return e, nil
		case "tab":
			e.inputs[e.focusedInput].Blur()
			e.focusedInput = (e.focusedInput + 1) % len(e.inputs)
			e.inputs[e.focusedInput].Focus()
		case "esc":
			return e, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 3}
			}
		}
	}

	e.inputs[e.focusedInput], _ = e.inputs[e.focusedInput].Update(msg)
	return e, nil
}
func (e *EditEndpointTab) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(2)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		MarginBottom(2)

	form := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("Edit Collection"),
		inputStyle.Render(e.inputs[0].View()),
		inputStyle.Render(e.inputs[1].View()),
		inputStyle.Render(e.inputs[2].View()),
	)

	containerStyle := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Align(lipgloss.Center, lipgloss.Center)

	return containerStyle.Render(form)
}
func (e *EditEndpointTab) OnFocus() tea.Cmd {
	e.inputs[0].Focus()
	e.focused = true
	return textinput.Blink
}
func (e *EditEndpointTab) OnBlur() tea.Cmd {
	e.inputs[e.focusedInput].Blur()
	e.focused = false
	return nil
}

func (e *EditEndpointTab) SetEditingEndpoint(msg messages.EditEndpointMsg) {
	e.inputs[0].SetValue(msg.Name)
	e.inputs[1].SetValue(msg.Method)
	e.inputs[2].SetValue(msg.URL)
	e.endpointID = msg.ID
}

func (e *EditEndpointTab) updateEndpoint(newName, newMethod, newURL string) (Tab, tea.Cmd) {
	ctx := global.GetAppContext()
	id, _ := strconv.ParseInt(e.endpointID, 10, 64)
	ctx.Endpoints.UpdateEndpoint(context.Background(), id, endpoints.EndpointData{
		Name:   newName,
		Method: newMethod,
		URL:    newURL,
	})

	return e, func() tea.Msg {
		return messages.SwitchTabMsg{TabIndex: 3}
	}
}
