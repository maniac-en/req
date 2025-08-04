package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	width  int
	height int
}

func (a AppModel) Init() tea.Cmd {
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a AppModel) View() string {
	return "Hello World"
}

func NewAppModel() AppModel {
	return AppModel{}
}
