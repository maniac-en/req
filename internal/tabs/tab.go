package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

var GlobalCollections = []OptionPair{}

// this is what a tab is loosely defined as
type Tab interface {
	Name() string
	Instructions() string
	Init() tea.Cmd
	Update(tea.Msg) (Tab, tea.Cmd)
	View() string
	OnFocus() tea.Cmd
	OnBlur() tea.Cmd
}
