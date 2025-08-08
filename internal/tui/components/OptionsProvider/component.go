package optionsProvider

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/keybinds"
)

type focusedComp string

const (
	listComponent = "list"
	textComponent = "text"
)

type OptionsProvider struct {
	list           list.Model
	input          OptionsInput
	onSelectAction tea.Msg
	width          int
	height         int
	focused        focusedComp
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
		switch o.focused {
		case listComponent:
			if !o.IsFiltering() {
				switch {
				case key.Matches(msg, keybinds.Keys.InsertItem):
					o.list.SetSize(o.list.Width(), o.height-lipgloss.Height(o.input.View()))
					o.input.OnFocus()
					o.focused = textComponent
					return o, tea.Batch(cmds...)
				}
			}
		}

	}
	switch o.focused {
	case listComponent:
		o.list, cmd = o.list.Update(msg)
	case textComponent:
		o.input, cmd = o.input.Update(msg)
	}
	cmds = append(cmds, cmd)
	return o, tea.Batch(cmds...)
}

func (o OptionsProvider) View() string {
	if o.focused == textComponent {
		return lipgloss.JoinVertical(lipgloss.Left, o.list.View(), o.input.View())
	}
	return o.list.View()
}

func (o *OptionsProvider) OnFocus() {
}

func (o OptionsProvider) OnBlur() {

}

func (o OptionsProvider) GetSelected() Option {
	return o.list.SelectedItem().(Option)
}

func (o OptionsProvider) IsFiltering() bool {
	return o.list.FilterState() == list.Filtering
}

func initList[T, U any](config *ListConfig[T, U]) list.Model {

	rawItems, err := config.CrudOps.List(context.Background())

	if err != nil {
		rawItems = []T{}
	}

	items := config.ItemMapper(rawItems)

	list := list.New(items, config.Delegate, 30, 30)

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
		list:    initList(config),
		focused: listComponent,
		input:   NewOptionsInput(),
		// onSelectAction: config.OnSelectAction,
	}
}
