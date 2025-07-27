package tabs

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/messages"
)

type EndpointsTab struct {
	name        string
	globalState *global.State
}

func NewEndpointsTab(state *global.State) *EndpointsTab {
	return &EndpointsTab{
		name:        "Endpoints",
		globalState: state,
	}
}

func (e *EndpointsTab) Name() string {
	return e.name
}

func (e *EndpointsTab) Instructions() string {
	return "c - collections page"
}

func (e *EndpointsTab) Init() tea.Cmd {
	return nil
}

func (e *EndpointsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			return e, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 0}
			}
		}
	}

	return e, nil
}

func (e *EndpointsTab) View() string {
	return fmt.Sprintf("This is the endpoints page\n %s", e.globalState.GetCurrentCollection())
}

func (e *EndpointsTab) OnFocus() tea.Cmd {
	return nil
}

func (e *EndpointsTab) OnBlur() tea.Cmd {
	return nil
}
