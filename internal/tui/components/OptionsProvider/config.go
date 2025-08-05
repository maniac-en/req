package optionsProvider

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListConfig[T any] struct {
	Items          []T
	OnSelectAction tea.Msg

	ShowPagination bool
	ShowHelp       bool
	ShowTitle      bool
	Width, Height  int

	FilteringEnabled bool

	Delegate list.ItemDelegate
	KeyMap   list.KeyMap

	ItemMapper func([]T) []list.Item

	// CrudOps        Crud
	// Style    lipgloss.Style
}
