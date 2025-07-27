package tabs

import (
	"context"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/messages"
)

type EndpointsTab struct {
	name        string
	globalState *global.State
	selectUI    SelectInput
	loaded      bool
}

type endpointListOpts struct {
	options []OptionPair
}

func NewEndpointsTab(state *global.State) *EndpointsTab {
	return &EndpointsTab{
		name:        "Endpoints",
		globalState: state,
		loaded:      false,
		selectUI:    NewSelectInput(),
	}
}

func (e *EndpointsTab) IsFiltering() bool {
	return e.selectUI.list.FilterState() == list.Filtering
}

func (e *EndpointsTab) Name() string {
	return e.name
}

func (e *EndpointsTab) Instructions() string {
	return "c - collections page • + - add endpoint • d - delete endpoint"
}

func (e *EndpointsTab) fetchEndpoints(collectionId string, limit, offset int) tea.Cmd {
	ctx := global.GetAppContext()
	collectionIdInt, err := strconv.ParseInt(collectionId, 10, 64)
	if err != nil {
		return func() tea.Msg {
			return endpointListOpts{}
		}
	}
	opts := []OptionPair{}
	endpoints, err := ctx.Endpoints.ListByCollection(context.Background(), collectionIdInt, limit, offset)
	for _, endpoint := range endpoints.Endpoints {
		opts = append(opts, OptionPair{
			Label: endpoint.GetName(),
			Value: strconv.FormatInt(endpoint.GetID(), 10),
		})
	}

	return func() tea.Msg {
		return endpointListOpts{
			options: opts,
		}
	}
}

func (e *EndpointsTab) Init() tea.Cmd {
	e.selectUI.Focus()
	return e.selectUI.Init()
}

func (e *EndpointsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case endpointListOpts:
		e.selectUI.SetOptions(msg.options)
		e.loaded = true
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			return e, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 0}
			}
		case "+":
			return e, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 4}
			}
		case "d":
			if selected := e.selectUI.GetSelected(); selected != "" {
				return e, e.deleteEndpoint(selected)
			}
		}
	default:
		e.selectUI, cmd = e.selectUI.Update(msg)
	}

	return e, cmd
}

func (e *EndpointsTab) View() string {

	if e.selectUI.IsLoading() {
		return e.selectUI.View()
	}
	selectContent := e.selectUI.View()

	style := lipgloss.NewStyle().
		PaddingRight(4)

	if !e.selectUI.IsLoading() && len(e.selectUI.list.Items()) > 0 {
		title := "\n\n\n\n\n\n\nSelect Endpoint:\n\n"
		return title + style.Render(selectContent)
	}

	return style.Render(selectContent)
}

func (e *EndpointsTab) OnFocus() tea.Cmd {
	e.selectUI.Focus()
	return e.fetchEndpoints(e.globalState.GetCurrentCollection(), 5, 0)
}

func (e *EndpointsTab) OnBlur() tea.Cmd {
	e.selectUI.Blur()
	return nil
}

func (e *EndpointsTab) deleteEndpoint(value string) tea.Cmd {
	ctx := global.GetAppContext()
	id, _ := strconv.ParseInt(value, 10, 64)
	err := ctx.Endpoints.Delete(context.Background(), id)
	if err != nil {
		return e.fetchEndpoints(e.globalState.GetCurrentCollection(), 5, 0)
	}
	for i, collection := range GlobalCollections {
		if collection.Value == value {
			GlobalCollections = append(GlobalCollections[:i], GlobalCollections[i+1:]...)
			break
		}
	}
	return e.fetchEndpoints(e.globalState.GetCurrentCollection(), 5, 0)
}
