package optionsProvider

import (
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
	title       string
	value       string
	description string
}

func (o Option) Title() string       { return o.title }
func (o Option) Description() string { return o.description }
func (o Option) Value() string       { return o.value }
func (o Option) FilterValue() string { return o.title }

func (o OptionsProvider) Init() tea.Cmd {
	return nil
}

func (o OptionsProvider) Update(msg tea.Msg) (OptionsProvider, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		o.height = msg.Height
		o.width = msg.Width
		o.list.SetSize(o.list.Width(), o.height)
	}
	return o, nil
}

func (o OptionsProvider) View() string {
	return o.list.View()
}

func (o OptionsProvider) OnFocus() {

}

func (o OptionsProvider) OnBlur() {

}

func initList[T, C any](config *ListConfig[T, C]) list.Model {

	// items := config.ItemMapper(config.Items)

	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 0)

	// list configuration
	list.SetFilteringEnabled(config.FilteringEnabled)
	list.SetShowStatusBar(config.ShowStatusBar)
	list.SetShowPagination(config.ShowPagination)
	list.SetShowHelp(config.ShowHelp)
	list.SetShowTitle(config.ShowTitle)

	// list.KeyMap = config.KeyMap

	return list
}

func NewOptionsProvider[T, C any](config *ListConfig[T, C]) OptionsProvider {
	return OptionsProvider{
		list: initList(config),
		// onSelectAction: config.OnSelectAction,
	}
}
