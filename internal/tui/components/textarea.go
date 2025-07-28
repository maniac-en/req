package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/styles"
)

type Textarea struct {
	label    string
	content  string
	width    int
	height   int
	focused  bool
	cursor   int
	lines    []string
	cursorRow int
	cursorCol int
	scrollOffset int
}

func NewTextarea(label, placeholder string) Textarea {
	return Textarea{
		label:        label,
		content:      "",
		width:        50,
		height:       6,
		focused:      false,
		cursor:       0,
		lines:        []string{""},
		cursorRow:    0,
		cursorCol:    0,
		scrollOffset: 0,
	}
}

func (t *Textarea) SetValue(value string) {
	t.content = value
	rawLines := strings.Split(value, "\n")
	if len(rawLines) == 0 {
		rawLines = []string{""}
	}
	
	// Wrap long lines to fit within the textarea width
	t.lines = []string{}
	contentWidth := t.getContentWidth()
	
	for _, line := range rawLines {
		if len(line) <= contentWidth {
			t.lines = append(t.lines, line)
		} else {
			// Wrap long lines
			wrapped := t.wrapLine(line, contentWidth)
			t.lines = append(t.lines, wrapped...)
		}
	}
	
	if len(t.lines) == 0 {
		t.lines = []string{""}
	}
	
	// Set cursor to end
	t.cursorRow = len(t.lines) - 1
	t.cursorCol = len(t.lines[t.cursorRow])
}

func (t Textarea) Value() string {
	return strings.Join(t.lines, "\n")
}

func (t *Textarea) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Textarea) Focus() {
	t.focused = true
}

func (t *Textarea) Blur() {
	t.focused = false
}

func (t Textarea) Focused() bool {
	return t.focused
}

func (t *Textarea) moveCursor(row, col int) {
	// Ensure row is in bounds
	if row < 0 {
		row = 0
	}
	if row >= len(t.lines) {
		row = len(t.lines) - 1
	}
	
	// Ensure col is in bounds for the row
	if col < 0 {
		col = 0
	}
	if col > len(t.lines[row]) {
		col = len(t.lines[row])
	}
	
	t.cursorRow = row
	t.cursorCol = col
}


func (t Textarea) Update(msg tea.Msg) (Textarea, tea.Cmd) {
	if !t.focused {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Insert new line
			currentLine := t.lines[t.cursorRow]
			beforeCursor := currentLine[:t.cursorCol]
			afterCursor := currentLine[t.cursorCol:]
			
			t.lines[t.cursorRow] = beforeCursor
			newLines := make([]string, len(t.lines)+1)
			copy(newLines[:t.cursorRow+1], t.lines[:t.cursorRow+1])
			newLines[t.cursorRow+1] = afterCursor
			copy(newLines[t.cursorRow+2:], t.lines[t.cursorRow+1:])
			t.lines = newLines
			
			t.cursorRow++
			t.cursorCol = 0
			
		case "tab":
			// Insert 2 spaces for indentation
			currentLine := t.lines[t.cursorRow]
			t.lines[t.cursorRow] = currentLine[:t.cursorCol] + "  " + currentLine[t.cursorCol:]
			t.cursorCol += 2
			
		case "backspace":
			if t.cursorCol > 0 {
				// Remove character
				currentLine := t.lines[t.cursorRow]
				t.lines[t.cursorRow] = currentLine[:t.cursorCol-1] + currentLine[t.cursorCol:]
				t.cursorCol--
			} else if t.cursorRow > 0 {
				// Join with previous line
				prevLine := t.lines[t.cursorRow-1]
				currentLine := t.lines[t.cursorRow]
				t.lines[t.cursorRow-1] = prevLine + currentLine
				
				newLines := make([]string, len(t.lines)-1)
				copy(newLines[:t.cursorRow], t.lines[:t.cursorRow])
				copy(newLines[t.cursorRow:], t.lines[t.cursorRow+1:])
				t.lines = newLines
				
				t.cursorRow--
				t.cursorCol = len(prevLine)
			}
			
		case "delete":
			if t.cursorCol < len(t.lines[t.cursorRow]) {
				// Remove character
				currentLine := t.lines[t.cursorRow]
				t.lines[t.cursorRow] = currentLine[:t.cursorCol] + currentLine[t.cursorCol+1:]
			} else if t.cursorRow < len(t.lines)-1 {
				// Join with next line
				currentLine := t.lines[t.cursorRow]
				nextLine := t.lines[t.cursorRow+1]
				t.lines[t.cursorRow] = currentLine + nextLine
				
				newLines := make([]string, len(t.lines)-1)
				copy(newLines[:t.cursorRow+1], t.lines[:t.cursorRow+1])
				copy(newLines[t.cursorRow+1:], t.lines[t.cursorRow+2:])
				t.lines = newLines
			}
			
		case "up":
			t.moveCursor(t.cursorRow-1, t.cursorCol)
		case "down":
			t.moveCursor(t.cursorRow+1, t.cursorCol)
		case "left":
			if t.cursorCol > 0 {
				t.cursorCol--
			} else if t.cursorRow > 0 {
				t.cursorRow--
				t.cursorCol = len(t.lines[t.cursorRow])
			}
		case "right":
			if t.cursorCol < len(t.lines[t.cursorRow]) {
				t.cursorCol++
			} else if t.cursorRow < len(t.lines)-1 {
				t.cursorRow++
				t.cursorCol = 0
			}
		case "home":
			t.cursorCol = 0
		case "end":
			t.cursorCol = len(t.lines[t.cursorRow])
			
		default:
			// Insert printable characters
			if len(msg.String()) == 1 && msg.String() >= " " {
				char := msg.String()
				currentLine := t.lines[t.cursorRow]
				t.lines[t.cursorRow] = currentLine[:t.cursorCol] + char + currentLine[t.cursorCol:]
				t.cursorCol++
			}
		}
	}

	return t, nil
}

func (t Textarea) View() string {
	// Use full width since we don't need label space
	containerWidth := t.width - 4 // Just account for padding
	if containerWidth < 20 {
		containerWidth = 20
	}

	// Create the textarea container
	containerHeight := t.height
	if containerHeight < 3 {
		containerHeight = 3
	}

	container := styles.ListItemStyle.Copy().
		Width(containerWidth).
		Height(containerHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(0, 1)

	if t.focused {
		container = container.BorderForeground(styles.Primary)
	}

	// Prepare visible lines with cursor
	visibleLines := make([]string, containerHeight-2) // Account for border
	for i := 0; i < len(visibleLines); i++ {
		lineIndex := i // No scrolling for now
		if lineIndex < len(t.lines) {
			line := t.lines[lineIndex]
			
			// Add cursor if this is the cursor row and textarea is focused
			if t.focused && lineIndex == t.cursorRow {
				if t.cursorCol <= len(line) {
					line = line[:t.cursorCol] + "│" + line[t.cursorCol:]
				}
			}
			
			// Lines should already be wrapped, no need to truncate
			
			visibleLines[i] = line
		} else {
			visibleLines[i] = ""
		}
	}

	content := strings.Join(visibleLines, "\n")
	textareaView := container.Render(content)

	return textareaView
}

func (t Textarea) getContentWidth() int {
	// Calculate content width (no label needed)
	containerWidth := t.width - 4 // Just account for padding
	if containerWidth < 20 {
		containerWidth = 20
	}
	contentWidth := containerWidth - 4 // border + padding
	if contentWidth < 10 {
		contentWidth = 10
	}
	return contentWidth
}

func (t Textarea) wrapLine(line string, maxWidth int) []string {
	if len(line) <= maxWidth {
		return []string{line}
	}
	
	var wrapped []string
	for len(line) > maxWidth {
		// Find the best place to break (prefer spaces)
		breakPoint := maxWidth
		for i := maxWidth - 1; i >= maxWidth-20 && i >= 0; i-- {
			if line[i] == ' ' {
				breakPoint = i
				break
			}
		}
		
		wrapped = append(wrapped, line[:breakPoint])
		line = line[breakPoint:]
		
		// Skip leading space on continuation lines
		if len(line) > 0 && line[0] == ' ' {
			line = line[1:]
		}
	}
	
	if len(line) > 0 {
		wrapped = append(wrapped, line)
	}
	
	return wrapped
}

