package keybinds

import (
	"github.com/charmbracelet/bubbles/key"
)

type Help struct {
	Keys []key.Binding
}

func (h Help) ShortHelp() []key.Binding {
	return h.Keys
}

func (h Help) FullHelp() [][]key.Binding {
	// TODO: Figure how you wanna show this
	return [][]key.Binding{}
}

func (h Help) SetHelp(helpMenu []key.Binding) {
	h.Keys = helpMenu
}
