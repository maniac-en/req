package optionsProvider

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/tui/keybinds"
	"github.com/maniac-en/req/internal/tui/messages"
	"github.com/maniac-en/req/internal/tui/styles"
)

type OptionsInput struct {
	input  textinput.Model
	height int
	width  int
	editId int
}

func NewOptionsInput(config *InputConfig) OptionsInput {
	input := textinput.New()
	input.CharLimit = config.CharLimit
	input.Placeholder = config.Placeholder
	input.Width = config.Width
	input.TextStyle = styles.InputStyle
	input.Prompt = config.Prompt

	return OptionsInput{
		input:  input,
		editId: -1,
	}
}

func (i OptionsInput) Init() tea.Cmd {
	return nil
}

func (i OptionsInput) Update(msg tea.Msg) (OptionsInput, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.Keys.Choose):
			return i, func() tea.Msg { return messages.ItemAdded{Item: i.input.Value()} }
		}
	}

	i.input, cmd = i.input.Update(msg)
	cmds = append(cmds, cmd)

	return i, tea.Batch(cmds...)
}

func (i OptionsInput) View() string {
	return styles.InputStyle.Render(i.input.View())
}

func (i *OptionsInput) SetInput(text string) {
	i.input.SetValue(text)
}

func (i *OptionsInput) OnFocus(id ...int) {
	if len(id) > 0 {
		i.editId = id[0]
	}
	i.input.Focus()
}

func (i *OptionsInput) OnBlur() {
	i.editId = -1
	i.input.Blur()
}
