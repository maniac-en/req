package keybinds

import (
	"github.com/charmbracelet/bubbles/key"
)

type Keymaps struct {
	InsertItem key.Binding
	DeleteItem key.Binding
	EditItem   key.Binding
	Choose     key.Binding
	Remove     key.Binding
	Back       key.Binding
}

var Keys = Keymaps{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
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
}
