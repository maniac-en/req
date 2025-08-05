package views

import (
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
	return "Collections"
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
	return c, tea.Batch(cmds...)
}

func (c CollectionsView) View() string {
	return c.list.View()
}

func (c CollectionsView) OnFocus() {

}

func (c CollectionsView) OnBlur() {

}

func NewCollectionsView(collManager *collections.CollectionsManager) *CollectionsView {
	config := defaultListConfig[string]()
	return &CollectionsView{
		list: optionsProvider.NewOptionsProvider(config),
	}
}
