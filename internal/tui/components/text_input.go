package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/styles"
)

type TextInput struct {
	textInput textinput.Model
	label     string
	width     int
}

func NewTextInput(label, placeholder string) TextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return TextInput{
		textInput: ti,
		label:     label,
		width:     50,
	}
}

func (t *TextInput) SetValue(value string) {
	t.textInput.SetValue(value)
}

func (t TextInput) Value() string {
	return t.textInput.Value()
}

func (t *TextInput) SetWidth(width int) {
	t.width = width
	t.textInput.Width = width - len(t.label) - 4 // Account for label and spacing
}

func (t *TextInput) Focus() {
	t.textInput.Focus()
}

func (t *TextInput) Blur() {
	t.textInput.Blur()
}

func (t *TextInput) Clear() {
	t.textInput.SetValue("")
}

func (t TextInput) Focused() bool {
	return t.textInput.Focused()
}

func (t TextInput) Update(msg tea.Msg) (TextInput, tea.Cmd) {
	var cmd tea.Cmd
	t.textInput, cmd = t.textInput.Update(msg)
	return t, cmd
}

func (t TextInput) View() string {
	labelStyle := styles.TitleStyle.Copy().
		Width(12).
		Align(lipgloss.Right)
	
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(t.label+":"),
		" ",
		t.textInput.View(),
	)
}