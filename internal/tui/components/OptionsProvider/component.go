package optionsProvider

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/log"
	input "github.com/maniac-en/req/internal/tui/components/Input"
	"github.com/maniac-en/req/internal/tui/keybinds"
	"github.com/maniac-en/req/internal/tui/messages"
)

type focusedComp string

const (
	listComponent = "list"
	textComponent = "text"
)

type OptionsProvider[T, U any] struct {
	list           list.Model
	input          input.OptionsInput
	onSelectAction tea.Msg
	keys           *keybinds.ListKeyMap
	width          int
	height         int
	focused        focusedComp
	getItems       func(context.Context) ([]T, error)
	itemMapper     func([]T) []list.Item
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

func (o OptionsProvider[T, U]) Init() tea.Cmd {
	return nil
}

func (o OptionsProvider[T, U]) Update(msg tea.Msg) (OptionsProvider[T, U], tea.Cmd) {
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
				case key.Matches(msg, o.keys.AddItem):
					o.list.SetSize(o.list.Width(), o.height-lipgloss.Height(o.input.View()))
					o.input.OnFocus()
					o.focused = textComponent
					return o, tea.Batch(cmds...)
				case key.Matches(msg, o.keys.DeleteItem):
					return o, func() tea.Msg { return messages.DeleteItem{ItemID: int64(o.GetSelected().ID)} }
				case key.Matches(msg, o.keys.Choose):
					return o, func() tea.Msg { return messages.ChooseCollection{} }
				case key.Matches(msg, o.keys.EditItem):
					o.list.SetSize(o.list.Width(), o.height-lipgloss.Height(o.input.View()))
					o.input.SetInput(o.GetSelected().Name)
					o.input.OnFocus(o.GetSelected().ID)
					o.focused = textComponent
					return o, tea.Batch(cmds...)
				}
			}
		}
	case messages.ItemAdded, messages.ItemEdited:
		o.input.OnBlur()
		o.focused = listComponent
		o.list.SetSize(o.list.Width(), o.height)
		o.RefreshItems()
	case messages.DeactivateView:
		o.input.OnBlur()
		o.focused = listComponent
		o.list.SetSize(o.list.Width(), o.height)
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

func (o OptionsProvider[T, U]) View() string {
	if o.focused == textComponent {
		return lipgloss.JoinVertical(lipgloss.Left, o.list.View(), o.input.View())
	}
	return o.list.View()
}

func (o *OptionsProvider[T, U]) OnFocus() {
}

func (o OptionsProvider[T, U]) OnBlur() {

}

func (o OptionsProvider[T, U]) GetSelected() Option {
	if o.IsFiltering() {
		return Option{
			Name:    "Filtering....",
			ID:      -1,
			Subtext: "",
		}
	}
	return o.list.SelectedItem().(Option)
}

func (o OptionsProvider[T, U]) IsFiltering() bool {
	return o.list.FilterState() == list.Filtering
}

func (o *OptionsProvider[T, U]) RefreshItems() {
	newItems, err := o.getItems(context.Background())
	if err != nil {
		log.Warn("Fetching items failed")
		return
	}
	o.list.SetItems(o.itemMapper(newItems))
}

func (o *OptionsProvider[T, U]) Help() []key.Binding {
	var binds []key.Binding
	switch o.focused {
	case listComponent:
		if o.IsFiltering() {
			binds = []key.Binding{
				o.keys.AcceptWhileFiltering,
				o.keys.CancelWhileFiltering,
				o.keys.ClearFilter,
			}
		} else {
			binds = []key.Binding{
				o.keys.CursorUp,
				o.keys.CursorDown,
				o.keys.NextPage,
				o.keys.PrevPage,
				o.keys.Filter,
				o.keys.AddItem,
				o.keys.EditItem,
				o.keys.DeleteItem,
			}
		}
	case textComponent:
		binds = o.input.Help()
	default:
		binds = []key.Binding{}
	}
	return binds
}

func initList[T, U any](config *ListConfig[T, U]) list.Model {

	rawItems, err := config.GetItemsFunc(context.Background())

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
	list.KeyMap = config.KeyMap

	return list
}

func NewOptionsProvider[T, U any](config *ListConfig[T, U]) OptionsProvider[T, U] {

	inputConfig := input.InputConfig{
		CharLimit:   100,
		Placeholder: "Add A New Collection...",
		Width:       22,
		Prompt:      "",
		KeyMap: input.InputKeyMaps{
			Accept: config.AdditionalKeymaps.Accept,
			Back:   config.AdditionalKeymaps.Back,
		},
	}

	return OptionsProvider[T, U]{
		list:       initList(config),
		focused:    listComponent,
		input:      input.NewOptionsInput(&inputConfig),
		getItems:   config.GetItemsFunc,
		itemMapper: config.ItemMapper,
		keys:       config.AdditionalKeymaps,
		// onSelectAction: config.OnSelectAction,
	}
}
