package tabs

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/messages"
)

type AddCollectionTab struct {
	name      string
	nameInput textinput.Model
	focused   bool
}

func NewAddCollectionTab() *AddCollectionTab {
	textInput := textinput.New()
	textInput.Placeholder = "Enter your collection's name... "
	textInput.Focus()

	textInput.CharLimit = 100
	textInput.Width = 50

	return &AddCollectionTab{
		name:      "Add Collection",
		nameInput: textInput,
		focused:   true,
	}
}

func (a *AddCollectionTab) Name() string {
	return a.name
}

func (a *AddCollectionTab) Instructions() string {
	return "Enter - create • Esc - cancel"
}

func (a *AddCollectionTab) Init() tea.Cmd {
	return textinput.Blink
}

func (a *AddCollectionTab) OnFocus() tea.Cmd {
	a.nameInput.Focus()
	a.focused = true
	return textinput.Blink
}

func (a *AddCollectionTab) OnBlur() tea.Cmd {
	a.nameInput.Blur()
	a.focused = false
	return nil
}

func (a *AddCollectionTab) Update(msg tea.Msg) (Tab, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if a.nameInput.Value() != "" {
				return a.addCollection(a.nameInput.Value())
			}
		case "esc":
			return a, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 0}
			}
		}
	}

	a.nameInput, _ = a.nameInput.Update(msg)
	return a, nil
}

func (a *AddCollectionTab) View() string {
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
		titleStyle.Render("Create New Collection"),
		inputStyle.Render(a.nameInput.View()),
	)

	containerStyle := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Align(lipgloss.Center, lipgloss.Center)

	return containerStyle.Render(form)
}

func (a *AddCollectionTab) addCollection(name string) (Tab, tea.Cmd) {
	ctx := global.GetAppContext()
	collection, _ := ctx.Collections.Create(context.Background(), name)
	newOption := OptionPair{
		Label: collection.GetName(),
		Value: string(collection.GetID()),
	}

	GlobalCollections = append(GlobalCollections, newOption)

	a.nameInput.SetValue("")

	return a, func() tea.Msg {
		return messages.SwitchTabMsg{TabIndex: 0}
	}
}
