package optionsProvider

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type OptionsInput struct {
	input  textinput.Model
	height int
	width  int
	editId int
}

func NewOptionsInput() OptionsInput {
	input := textinput.New()
	input.CharLimit = 100
	input.Placeholder = "Add New Collection..."
	input.Width = 22
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

	i.input, cmd = i.input.Update(msg)
	cmds = append(cmds, cmd)

	return i, tea.Batch(cmds...)
}

func (i OptionsInput) View() string {
	return i.input.View()
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
