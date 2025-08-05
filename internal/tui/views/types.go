package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ViewInterface interface {
	Init() tea.Cmd
	Name() string
	Help() string
	GetFooterSegment() string
	Update(tea.Msg) (ViewInterface, tea.Cmd)
	View() string
	OnFocus()
	OnBlur()
}
