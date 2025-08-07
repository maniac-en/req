package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
)

type CollectionsView struct {
	width  int
	height int
	list   optionsProvider.OptionsProvider
}

func (c CollectionsView) Init() tea.Cmd {
	return nil
}

func (c CollectionsView) Name() string {
	return "Collections"
}

func (c CollectionsView) Help() string {
	return ""
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
	}

	c.list, cmd = c.list.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c CollectionsView) View() string {
	return c.list.View()
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
			Subtext: "Sample",
			ID:      item.GetID(),
		}
		opts[i] = newOpt
	}
	return opts
}

func NewCollectionsView(collManager *collections.CollectionsManager) *CollectionsView {
	config := defaultListConfig[collections.CollectionEntity, string]()
	config.CrudOps = optionsProvider.Crud[collections.CollectionEntity, string]{
		Create: collManager.Create,
		Read:   collManager.Read,
		Update: collManager.Update,
		Delete: collManager.Delete,
		List:   collManager.List,
	}
	config.ItemMapper = itemMapper
	return &CollectionsView{
		list: optionsProvider.NewOptionsProvider(config),
	}
}
