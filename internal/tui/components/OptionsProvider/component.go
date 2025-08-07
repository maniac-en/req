package optionsProvider

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type OptionsProvider struct {
	list           list.Model
	onSelectAction tea.Msg
	width          int
	height         int
}

type Option struct {
	Name    string
	ID      int64
	Subtext string
}

func (o Option) Title() string       { return o.Name }
func (o Option) Description() string { return o.Subtext }
func (o Option) Value() int64        { return o.ID }
func (o Option) FilterValue() string { return o.Name }

func (o OptionsProvider) Init() tea.Cmd {
	return nil
}

func (o OptionsProvider) Update(msg tea.Msg) (OptionsProvider, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		o.height = msg.Height
		o.width = msg.Width
		o.list.SetSize(o.list.Width(), o.height)
	case tea.KeyMsg:

		o.list, cmd = o.list.Update(msg)
		cmds = append(cmds, cmd)

		switch msg.String() {

		default:
		}

	}
	return o, tea.Batch(cmds...)
}

func (o OptionsProvider) View() string {
	return o.list.View()
}

func (o OptionsProvider) OnFocus() {

}

func (o OptionsProvider) OnBlur() {

}

func (o OptionsProvider) GetSelected() Option {
	return o.list.SelectedItem().(Option)
}

func initList[T, U any](config *ListConfig[T, U]) list.Model {

	rawItems, err := config.CrudOps.List(context.Background())

	if err != nil {
		rawItems = []T{}
	}

	items := config.ItemMapper(rawItems)

	list := list.New(items, list.NewDefaultDelegate(), 30, 30)

	// list configuration
	list.SetFilteringEnabled(config.FilteringEnabled)
	list.SetShowStatusBar(config.ShowStatusBar)
	list.SetShowPagination(config.ShowPagination)
	list.SetShowHelp(config.ShowHelp)
	list.SetShowTitle(config.ShowTitle)

	// list.KeyMap = config.KeyMap

	return list
}

func NewOptionsProvider[T, U any](config *ListConfig[T, U]) OptionsProvider {
	return OptionsProvider{
		list: initList(config),
		// onSelectAction: config.OnSelectAction,
	}
}
