package optionsProvider

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListConfig[T, C any] struct {
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

	CrudOps Crud[C]
	// Style    lipgloss.Style
}

type Crud[T any] struct {
	Create func(context.Context, string) (T, error)
	Read   func(context.Context, int64) (T, error)
	Update func(context.Context, int64, string) (T, error)
	Delete func(context.Context, int64) error
	List   func(context.Context) ([]T, error)
}
