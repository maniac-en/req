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
	ti.CharLimit = 5000 // Allow long content like JSON
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
	// Account for label, colon, spacing, and border padding
	containerWidth := width - 12 - 1 - 2 // 12 for label, 1 for colon, 2 for spacing
	if containerWidth < 15 {
		containerWidth = 15
	}

	// The actual input width inside the container (subtract border and padding)
	inputWidth := containerWidth - 4 // 2 for border, 2 for padding
	if inputWidth < 10 {
		inputWidth = 10
	}

	// Ensure the underlying textinput respects the width
	t.textInput.Width = inputWidth
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
		MarginTop(1).
		Align(lipgloss.Right)

	// Create a fixed-width container for the input to prevent overflow
	containerWidth := t.width - 12 - 1 - 2 // Account for label, colon, spacing
	if containerWidth < 15 {
		containerWidth = 15
	}

	inputContainer := styles.ListItemStyle.Copy().
		Width(containerWidth).
		Height(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(0, 1)

	if t.textInput.Focused() {
		inputContainer = inputContainer.BorderForeground(styles.Primary)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(t.label+":"),
		" ",
		inputContainer.Render(t.textInput.View()),
	)
}
