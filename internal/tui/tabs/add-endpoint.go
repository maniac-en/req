package tabs

import (
	"context"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/tui/messages"
)

type AddEndpointTab struct {
	name         string
	inputs       []textinput.Model
	focusedInput int
	state        *global.State
	focused      bool
}

func NewAddEndpointTab(globalState *global.State) *AddEndpointTab {
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
		name: "Add Endpoint",
		inputs: []textinput.Model{
			name,
			method,
			url,
		},
		focused: true,
		state:   globalState,
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
			if a.inputs[0].Value() != "" && a.inputs[1].Value() != "" && a.inputs[2].Value() != "" {
				return a.addEndpoint(a.inputs[0].Value(), a.inputs[1].Value(), a.inputs[2].Value())
			}
			return a, nil
		case "tab":
			a.inputs[a.focusedInput].Blur()
			a.focusedInput = (a.focusedInput + 1) % len(a.inputs)
			a.inputs[a.focusedInput].Focus()
		case "esc":
			return a, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 3}
			}
		}
	}

	a.inputs[a.focusedInput], _ = a.inputs[a.focusedInput].Update(msg)
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
		inputStyle.Render(a.inputs[0].View()),
		inputStyle.Render(a.inputs[1].View()),
		inputStyle.Render(a.inputs[2].View()),
	)

	containerStyle := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Align(lipgloss.Center, lipgloss.Center)

	return containerStyle.Render(form)
}

func (a *AddEndpointTab) OnFocus() tea.Cmd {
	a.inputs[a.focusedInput].Focus()
	a.focused = true
	return textinput.Blink
}

func (a *AddEndpointTab) OnBlur() tea.Cmd {
	a.inputs[0].Blur()
	a.inputs[1].Blur()
	a.inputs[2].Blur()
	a.focused = false
	return nil
}

func (a *AddEndpointTab) addEndpoint(name, method, url string) (Tab, tea.Cmd) {
	ctx := global.GetAppContext()

	collectionId := a.state.GetCurrentCollection()
	int64Collection, err := strconv.ParseInt(collectionId, 10, 64)
	if err != nil {
		return a, func() tea.Msg {
			return messages.SwitchTabMsg{TabIndex: 3}
		}
	}

	_, _ = ctx.Endpoints.CreateEndpoint(context.Background(), endpoints.EndpointData{
		Name:         name,
		Method:       method,
		URL:          url,
		CollectionID: int64Collection,
	})

	// newOption := OptionPair{
	// 	Label: collection.GetName(),
	// 	Value: string(collection.GetID()),
	// }

	a.inputs[0].SetValue("")
	a.inputs[1].SetValue("")
	a.inputs[2].SetValue("")

	return a, func() tea.Msg {
		return messages.SwitchTabMsg{TabIndex: 3}
	}
}
