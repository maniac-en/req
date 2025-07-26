package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

var GlobalCollections = []OptionPair{
	{Label: "Collection 1", Value: "1"},
	{Label: "Collection 2", Value: "2"},
	{Label: "Collection 3", Value: "3"},
	{Label: "Collection 4", Value: "4"},
	{Label: "Collection 5", Value: "5"},
}

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
