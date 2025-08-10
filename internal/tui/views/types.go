package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewInterface interface {
	Init() tea.Cmd
	Name() string
	Help() []key.Binding
	GetFooterSegment() string
	Update(tea.Msg) (ViewInterface, tea.Cmd)
	View() string
	OnFocus()
	OnBlur()
}
