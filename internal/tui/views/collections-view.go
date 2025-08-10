package views

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
	"github.com/maniac-en/req/internal/tui/keybinds"
	"github.com/maniac-en/req/internal/tui/messages"
)

type CollectionsView struct {
	width   int
	height  int
	list    optionsProvider.OptionsProvider[collections.CollectionEntity, string]
	manager *collections.CollectionsManager
	help    help.Model
	keys    *keybinds.ListKeyMap
	order   int
}

func (c CollectionsView) Init() tea.Cmd {
	return nil
}

func (c CollectionsView) Name() string {
	return "Collections"
}

func (c CollectionsView) Help() []key.Binding {
	return c.list.Help()
}

func (c CollectionsView) GetFooterSegment() string {
	return c.list.GetSelected().Title()
}

func (c CollectionsView) Update(msg tea.Msg) (ViewInterface, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.height = msg.Height
		c.width = msg.Width
		c.list, cmd = c.list.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ItemAdded:
		c.manager.Create(context.Background(), msg.Item)
	case messages.ItemEdited:
		c.manager.Update(context.Background(), msg.ItemID, msg.Item)
	case messages.DeleteItem:
		c.manager.Delete(context.Background(), msg.ItemID)
		c.list.RefreshItems()
	}

	c.list, cmd = c.list.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c CollectionsView) View() string {
	return c.list.View()
}

func (c CollectionsView) Order() int {
	return c.order
}

func (c CollectionsView) OnFocus() {

}

func (c CollectionsView) OnBlur() {

}

func itemMapper(items []collections.CollectionEntity) []list.Item {
	opts := make([]list.Item, len(items))
	for i, item := range items {
		newOpt := optionsProvider.Option{
			Name:    item.GetName(),
			Subtext: fmt.Sprintf("%d endpoints", item.GetEnpointCount()),
			ID:      item.GetID(),
		}
		opts[i] = newOpt
	}
	return opts
}

func NewCollectionsView(collManager *collections.CollectionsManager, order int) *CollectionsView {
	keybinds := keybinds.NewListKeyMap()
	config := defaultListConfig[collections.CollectionEntity, string](keybinds)

	config.GetItemsFunc = collManager.List
	config.ItemMapper = itemMapper
	config.AdditionalKeymaps = keybinds

	return &CollectionsView{
		list:    optionsProvider.NewOptionsProvider(config),
		manager: collManager,
		help:    help.New(),
		keys:    keybinds,
		order:   order,
	}
}
