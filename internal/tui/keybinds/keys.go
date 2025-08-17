package keybinds

import (
	"github.com/charmbracelet/bubbles/key"
)

type Keymaps struct {
	InsertItem           key.Binding
	DeleteItem           key.Binding
	EditItem             key.Binding
	Choose               key.Binding
	Remove               key.Binding
	Back                 key.Binding
	Up                   key.Binding
	Down                 key.Binding
	NextPage             key.Binding
	PrevPage             key.Binding
	Filter               key.Binding
	ClearFilter          key.Binding
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding
	Quit                 key.Binding
}

var Keys = Keymaps{
	Back: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "back"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "back"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("left", "h", "pgup"),
		key.WithHelp("←/h/pgup", "prev page"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("right", "l", "pgdown"),
		key.WithHelp("→/l/pgdn", "next page"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	CancelWhileFiltering: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	AcceptWhileFiltering: key.NewBinding(
		key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
		key.WithHelp("enter", "apply filter"),
	),
	InsertItem: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add item"),
	),
	EditItem: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit item"),
	),
	Choose: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Choose"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x", "backspace"),
		key.WithHelp("x", "delete"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
