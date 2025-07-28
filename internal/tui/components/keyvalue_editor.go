package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/styles"
)

type KeyValuePair struct {
	Key     string
	Value   string
	Enabled bool
}

type KeyValueEditor struct {
	label      string
	pairs      []KeyValuePair
	width      int
	height     int
	focused    bool
	focusIndex int // Which pair is focused
	fieldIndex int // 0=key, 1=value, 2=enabled
}

func NewKeyValueEditor(label string) KeyValueEditor {
	return KeyValueEditor{
		label:      label,
		pairs:      []KeyValuePair{{"", "", true}}, // Start with one empty pair
		width:      50,
		height:     6,
		focused:    false,
		focusIndex: 0,
		fieldIndex: 0,
	}
}

func (kv *KeyValueEditor) SetSize(width, height int) {
	kv.width = width
	kv.height = height
}

func (kv *KeyValueEditor) Focus() {
	kv.focused = true
}

func (kv *KeyValueEditor) Blur() {
	kv.focused = false
}

func (kv KeyValueEditor) Focused() bool {
	return kv.focused
}

func (kv *KeyValueEditor) SetPairs(pairs []KeyValuePair) {
	if len(pairs) == 0 {
		kv.pairs = []KeyValuePair{{"", "", true}}
	} else {
		kv.pairs = pairs
	}
	// Ensure focus is within bounds
	if kv.focusIndex >= len(kv.pairs) {
		kv.focusIndex = len(kv.pairs) - 1
	}
}

func (kv KeyValueEditor) GetPairs() []KeyValuePair {
	return kv.pairs
}

func (kv KeyValueEditor) GetEnabledPairsAsMap() map[string]string {
	result := make(map[string]string)
	for _, pair := range kv.pairs {
		if pair.Enabled && pair.Key != "" {
			result[pair.Key] = pair.Value
		}
	}
	return result
}

func (kv KeyValueEditor) Update(msg tea.Msg) (KeyValueEditor, tea.Cmd) {
	if !kv.focused {
		return kv, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Move to next field
			kv.fieldIndex++
			if kv.fieldIndex > 2 { // key, value, enabled
				kv.fieldIndex = 0
				kv.focusIndex++
				if kv.focusIndex >= len(kv.pairs) {
					kv.focusIndex = 0
				}
			}
		case "shift+tab":
			// Move to previous field
			kv.fieldIndex--
			if kv.fieldIndex < 0 {
				kv.fieldIndex = 2
				kv.focusIndex--
				if kv.focusIndex < 0 {
					kv.focusIndex = len(kv.pairs) - 1
				}
			}
		case "up":
			if kv.focusIndex > 0 {
				kv.focusIndex--
			}
		case "down":
			if kv.focusIndex < len(kv.pairs)-1 {
				kv.focusIndex++
			}
		case "ctrl+n":
			// Add new pair
			kv.pairs = append(kv.pairs, KeyValuePair{"", "", true})
		case "ctrl+d":
			// Delete current pair (but keep at least one)
			if len(kv.pairs) > 1 {
				kv.pairs = append(kv.pairs[:kv.focusIndex], kv.pairs[kv.focusIndex+1:]...)
				if kv.focusIndex >= len(kv.pairs) {
					kv.focusIndex = len(kv.pairs) - 1
				}
			}
		case " ":
			// Toggle enabled state when on enabled field
			if kv.fieldIndex == 2 {
				kv.pairs[kv.focusIndex].Enabled = !kv.pairs[kv.focusIndex].Enabled
			}
		case "backspace":
			// Delete character from current field
			if kv.fieldIndex == 0 && len(kv.pairs[kv.focusIndex].Key) > 0 {
				kv.pairs[kv.focusIndex].Key = kv.pairs[kv.focusIndex].Key[:len(kv.pairs[kv.focusIndex].Key)-1]
			} else if kv.fieldIndex == 1 && len(kv.pairs[kv.focusIndex].Value) > 0 {
				kv.pairs[kv.focusIndex].Value = kv.pairs[kv.focusIndex].Value[:len(kv.pairs[kv.focusIndex].Value)-1]
			}
		default:
			// Add printable characters
			if len(msg.String()) == 1 && msg.String() >= " " {
				char := msg.String()
				if kv.fieldIndex == 0 {
					kv.pairs[kv.focusIndex].Key += char
				} else if kv.fieldIndex == 1 {
					kv.pairs[kv.focusIndex].Value += char
				}
			}
		}
	}

	return kv, nil
}

func (kv KeyValueEditor) View() string {
	// Calculate container dimensions (use full width like textarea)
	containerWidth := kv.width - 4 // Just account for padding
	if containerWidth < 30 {
		containerWidth = 30
	}

	container := styles.ListItemStyle.Copy().
		Width(containerWidth).
		Height(kv.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(1, 1)

	if kv.focused {
		container = container.BorderForeground(styles.Primary)
	}

	// Build content
	var lines []string
	visibleHeight := kv.height - 2 // Account for border

	// Header - better column proportions
	headerStyle := styles.ListItemStyle.Copy().Bold(true)
	availableWidth := containerWidth - 8      // Account for padding and separators
	keyWidth := availableWidth * 40 / 100     // 40% for key
	valueWidth := availableWidth * 50 / 100   // 50% for value
	enabledWidth := availableWidth * 10 / 100 // 10% for enabled

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Copy().Width(keyWidth).Render("Key"),
		"  ",
		headerStyle.Copy().Width(valueWidth).Render("Value"),
		"  ",
		headerStyle.Copy().Width(enabledWidth).Align(lipgloss.Center).Render("On"),
	)
	lines = append(lines, header)

	// Show pairs (limit to visible height)
	maxPairs := visibleHeight - 2 // Reserve space for header and instructions
	if maxPairs < 1 {
		maxPairs = 1
	}

	for i := 0; i < maxPairs && i < len(kv.pairs); i++ {
		pair := kv.pairs[i]

		// Style fields based on focus
		keyStyle := styles.ListItemStyle.Copy().Width(keyWidth)
		valueStyle := styles.ListItemStyle.Copy().Width(valueWidth)
		enabledStyle := styles.ListItemStyle.Copy().Width(enabledWidth).Align(lipgloss.Center)

		if kv.focused && i == kv.focusIndex {
			if kv.fieldIndex == 0 {
				keyStyle = keyStyle.Background(styles.Primary).Foreground(styles.TextPrimary)
			} else if kv.fieldIndex == 1 {
				valueStyle = valueStyle.Background(styles.Primary).Foreground(styles.TextPrimary)
			} else if kv.fieldIndex == 2 {
				enabledStyle = enabledStyle.Background(styles.Primary).Foreground(styles.TextPrimary)
			}
		}

		// Truncate long text
		keyText := pair.Key
		if len(keyText) > keyWidth-2 {
			keyText = keyText[:keyWidth-2]
		}
		valueText := pair.Value
		if len(valueText) > valueWidth-2 {
			valueText = valueText[:valueWidth-2]
		}

		checkbox := "☐"
		if pair.Enabled {
			checkbox = "☑"
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			keyStyle.Render(keyText),
			"  ",
			valueStyle.Render(valueText),
			"  ",
			enabledStyle.Render(checkbox),
		)
		lines = append(lines, row)
	}

	// Add instructions at bottom
	if len(lines) < visibleHeight-1 {
		instructions := "tab: next field • ↑↓: navigate rows • space: toggle"
		instrStyle := styles.ListItemStyle.Copy().Foreground(styles.TextMuted)
		lines = append(lines, "", instrStyle.Render(instructions))
	}

	// Fill remaining space
	for len(lines) < visibleHeight {
		lines = append(lines, "")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	containerView := container.Render(content)

	return containerView
}
