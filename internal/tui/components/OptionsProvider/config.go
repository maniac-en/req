package optionsProvider

import (
	"context"
	"github.com/charmbracelet/bubbles/key"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/tui/keybinds"
)

type ListConfig[T, U any] struct {
	OnSelectAction tea.Msg

	ShowPagination bool
	ShowStatusBar  bool
	ShowHelp       bool
	ShowTitle      bool
	Width, Height  int

	FilteringEnabled bool

	Delegate          list.ItemDelegate
	KeyMap            list.KeyMap
	AdditionalKeymaps *keybinds.ListKeyMap

	ItemMapper func([]T) []list.Item

	GetItemsFunc func(context.Context) ([]T, error)
	// Style    lipgloss.Style
}

type InputConfig struct {
	Prompt      string
	Placeholder string
	CharLimit   int
	Width       int
	KeyMap      InputKeyMaps
}

type InputKeyMaps struct {
	Accept key.Binding
	Back   key.Binding
}
