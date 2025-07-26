package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// TODO: Rework the `ActiveTab` logic so it doesnt use generic numbers
	case tea.KeyMsg:
		switch msg.String() {
		case m.Keybinds.Quit, m.Keybinds.KeyboardInterrupt:
			return m, tea.Quit
		case m.Keybinds.Collections:
			m.ActiveTab = 1
		case m.Keybinds.Endpoints:
			m.ActiveTab = 2
		case m.Keybinds.Environments:
			m.ActiveTab = 3
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	return m, nil
}
