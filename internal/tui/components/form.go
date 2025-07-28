package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/styles"
)

type Form struct {
	inputs     []TextInput
	focusIndex int
	width      int
	height     int
	title      string
	submitText string
	cancelText string
}

func NewForm(title string, inputs []TextInput) Form {
	if len(inputs) > 0 {
		inputs[0].Focus()
	}

	return Form{
		inputs:     inputs,
		focusIndex: 0,
		title:      title,
		submitText: "Submit",
		cancelText: "Cancel",
	}
}

func (f *Form) SetSize(width, height int) {
	f.width = width
	f.height = height

	for i := range f.inputs {
		f.inputs[i].SetWidth(width - 4)
	}
}

func (f *Form) SetSubmitText(text string) {
	f.submitText = text
}

func (f *Form) SetCancelText(text string) {
	f.cancelText = text
}

func (f Form) GetInput(index int) *TextInput {
	if index >= 0 && index < len(f.inputs) {
		return &f.inputs[index]
	}
	return nil
}

func (f Form) GetValues() []string {
	values := make([]string, len(f.inputs))
	for i, input := range f.inputs {
		values[i] = input.Value()
	}
	return values
}

func (f *Form) Clear() {
	for i := range f.inputs {
		f.inputs[i].Clear()
	}
}

func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.nextInput()
		case "shift+tab", "up":
			f.prevInput()
		}
	}

	if f.focusIndex >= 0 && f.focusIndex < len(f.inputs) {
		f.inputs[f.focusIndex], cmd = f.inputs[f.focusIndex].Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return f, tea.Batch(cmds...)
}

func (f *Form) nextInput() {
	if len(f.inputs) == 0 {
		return
	}

	f.inputs[f.focusIndex].Blur()
	f.focusIndex = (f.focusIndex + 1) % len(f.inputs)
	f.inputs[f.focusIndex].Focus()
}

func (f *Form) prevInput() {
	if len(f.inputs) == 0 {
		return
	}

	f.inputs[f.focusIndex].Blur()
	f.focusIndex--
	if f.focusIndex < 0 {
		f.focusIndex = len(f.inputs) - 1
	}
	f.inputs[f.focusIndex].Focus()
}

func (f Form) View() string {
	var content []string

	for _, input := range f.inputs {
		content = append(content, input.View())
	}

	content = append(content, "")

	buttonStyle := styles.ListItemStyle.Copy().
		Padding(0, 2).
		Background(styles.Primary).
		Foreground(styles.TextPrimary).
		Bold(true)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		buttonStyle.Render(f.submitText+" (enter)"),
		"  ",
		buttonStyle.Copy().
			Background(styles.TextSecondary).
			Render(f.cancelText+" (esc)"),
	)
	content = append(content, buttons)

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}
