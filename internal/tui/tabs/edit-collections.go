package tabs

import (
	"context"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/tui/messages"
)

type EditCollectionTab struct {
	name          string
	nameInput     textinput.Model
	originalValue string
	focused       bool
}

func NewEditCollectionTab() *EditCollectionTab {
	textInput := textinput.New()
	textInput.Placeholder = "Enter collection name..."
	textInput.CharLimit = 50
	textInput.Width = 30

	return &EditCollectionTab{
		name:      "Edit Collection",
		nameInput: textInput,
		focused:   true,
	}
}

func (e *EditCollectionTab) Name() string {
	return e.name
}

func (e *EditCollectionTab) Instructions() string {
	return "Enter - create • Esc - cancel"
}

func (e *EditCollectionTab) Init() tea.Cmd {
	return textinput.Blink
}

func (e *EditCollectionTab) OnFocus() tea.Cmd {
	e.nameInput.Focus()
	e.focused = true
	return textinput.Blink
}

func (e *EditCollectionTab) OnBlur() tea.Cmd {
	e.nameInput.Blur()
	e.focused = false
	return nil
}

func (e *EditCollectionTab) SetEditingCollection(label, value string) {
	e.nameInput.SetValue(label)
	e.originalValue = value
}

func (e *EditCollectionTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if e.nameInput.Value() != "" {
				return e.updateCollection(e.nameInput.Value())
			}
			return e, nil
		case "esc":
			return e, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 0}
			}
		default:
			e.nameInput, cmd = e.nameInput.Update(msg)
			return e, cmd
		}
	default:
		e.nameInput, cmd = e.nameInput.Update(msg)
		return e, cmd
	}
}

func (e *EditCollectionTab) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(2)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		MarginBottom(2)

	form := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("Edit Collection"),
		inputStyle.Render(e.nameInput.View()),
	)

	containerStyle := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Align(lipgloss.Center, lipgloss.Center)

	return containerStyle.Render(form)
}

func (e *EditCollectionTab) updateCollection(newName string) (Tab, tea.Cmd) {
	ctx := global.GetAppContext()
	id, _ := strconv.Atoi(e.originalValue)
	ctx.Collections.Update(context.Background(), int64(id), newName)
	for i, collection := range GlobalCollections {
		if collection.Value == e.originalValue {
			GlobalCollections[i] = OptionPair{
				Label: newName,
				Value: e.originalValue,
			}
			break
		}
	}

	return e, func() tea.Msg {
		return messages.SwitchTabMsg{TabIndex: 0}
	}
}
