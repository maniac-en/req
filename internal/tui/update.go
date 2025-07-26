package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case m.Keybinds.Quit, m.Keybinds.KeyboardInterrupt:
			return m, tea.Quit
		case m.Keybinds.Collections:
			m.ActiveTab = &m.Tabs.Collections
		case m.Keybinds.Endpoints:
			m.ActiveTab = &m.Tabs.Endpoints
		case m.Keybinds.Environments:
			m.ActiveTab = &m.Tabs.Environments
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}
