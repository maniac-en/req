package views

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/database"
	"github.com/maniac-en/req/internal/backend/endpoints"
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
	"github.com/maniac-en/req/internal/tui/keybinds"
	"github.com/maniac-en/req/internal/tui/messages"
)

type EndpointsView struct {
	height     int
	collection optionsProvider.Option
	width      int
	order      int
	list       optionsProvider.OptionsProvider[endpoints.EndpointEntity, database.Endpoint]
	manager    *endpoints.EndpointsManager
}

func (e *EndpointsView) Init() tea.Cmd {
	return nil
}

func (e *EndpointsView) Name() string {
	return "Endpoints"
}

func (e *EndpointsView) Help() []key.Binding {
	return e.list.Help()
}

func (e *EndpointsView) GetFooterSegment() string {
	return fmt.Sprintf("%s/", e.collection.Name)
}

func (e *EndpointsView) Update(msg tea.Msg) (ViewInterface, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.height = msg.Height
		e.width = msg.Width
		e.list, cmd = e.list.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ItemAdded:
		e.manager.CreateEndpoint(context.Background(), endpoints.EndpointData{
			CollectionID: e.collection.ID,
			Name:         msg.Item,
			Method:       "GET",
		})
	case messages.ItemEdited:
		e.manager.UpdateEndpointName(context.Background(), msg.ItemID, msg.Item)
	case messages.DeleteItem:
		e.manager.Delete(context.Background(), msg.ItemID)
		e.list.RefreshItems()
	}

	e.list, cmd = e.list.Update(msg)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e *EndpointsView) View() string {
	return e.list.View()
}

func (e *EndpointsView) OnFocus() {

}

func (e *EndpointsView) SetState(items ...any) error {
	if len(items) == 1 {
		if collection, ok := items[0].(optionsProvider.Option); ok {
			e.collection = collection
			epListFunc := func(ctx context.Context) ([]endpoints.EndpointEntity, error) {
				return e.manager.ListByCollection(ctx, collection.ID)
			}
			e.list.SetGetItemsFunc(epListFunc)
			return nil
		}
	}
	return errors.New("Invalid inputs, this function takes 1 input of type optionsProvider.Options")
}

func (e *EndpointsView) OnBlur() {

}

func (e *EndpointsView) Order() int {
	return e.order
}

func itemMapperEp(items []endpoints.EndpointEntity) []list.Item {
	opts := make([]list.Item, len(items))
	for i, item := range items {
		newOpt := optionsProvider.Option{
			Name:    item.GetName(),
			Subtext: item.Method,
			ID:      item.GetID(),
		}
		opts[i] = newOpt
	}
	return opts
}

func NewEndpointsView(epManager *endpoints.EndpointsManager, order int) *EndpointsView {
	view := &EndpointsView{
		order: order,
		collection: optionsProvider.Option{
			Name:    "",
			Subtext: "",
			ID:      0,
		},
		manager: epManager,
	}

	keybinds := keybinds.NewListKeyMap()
	config := defaultListConfig[endpoints.EndpointEntity, database.Endpoint](keybinds)

	epListFunc := func(ctx context.Context) ([]endpoints.EndpointEntity, error) {
		return epManager.ListByCollection(ctx, view.collection.ID)
	}

	config.GetItemsFunc = epListFunc
	config.ItemMapper = itemMapperEp
	config.AdditionalKeymaps = keybinds
	config.Source = "collections"

	view.list = optionsProvider.NewOptionsProvider(config)

	return view
}
