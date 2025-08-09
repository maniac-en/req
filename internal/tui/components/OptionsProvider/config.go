package optionsProvider

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListConfig[T, U any] struct {
	OnSelectAction tea.Msg

	ShowPagination bool
	ShowStatusBar  bool
	ShowHelp       bool
	ShowTitle      bool
	Width, Height  int

	FilteringEnabled bool

	Delegate list.ItemDelegate
	KeyMap   list.KeyMap

	ItemMapper func([]T) []list.Item

	GetItemsFunc func(context.Context) ([]T, error)
	// Style    lipgloss.Style
}

type InputConfig struct {
	Prompt      string
	Placeholder string
	CharLimit   int
	Width       int
}
