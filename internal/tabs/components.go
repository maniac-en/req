package tabs

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// for now the focus is just on the collections tab
// we'll see how we can change this around to accommodate
// more tabs
type collection string

type renderMethod struct{}

func (c collection) FilterValue() string { return string(c) }

func (r renderMethod) Height() int {
	return 1
}

func (r renderMethod) Spacing() int {
	return 0
}

func (r renderMethod) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (r renderMethod) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(collection)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true).
				PaddingLeft(2).
				Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type SelectInput struct {
	list    list.Model
	loading bool
	focused bool
	spinner spinner.Model
}

func NewSelectInput() SelectInput {
	l := list.New([]list.Item{}, renderMethod{}, 50, 14)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.SetShowTitle(false)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return SelectInput{
		list:    l,
		loading: true,
		focused: false,
		spinner: s,
	}
}
func (s SelectInput) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s SelectInput) Update(msg tea.Msg) (SelectInput, tea.Cmd) {
	var cmd tea.Cmd

	if s.loading {
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}

	if s.focused && !s.loading {
		s.list, cmd = s.list.Update(msg)
	}

	return s, cmd
}

func (s SelectInput) View() string {
	if s.loading {
		return fmt.Sprintf("%s Loading options...", s.spinner.View())
	}
	return s.list.View()
}

func (s SelectInput) Focused() bool   { return s.focused }
func (s *SelectInput) Focus()         { s.focused = true }
func (s *SelectInput) Blur()          { s.focused = false }
func (s SelectInput) IsLoading() bool { return s.loading }

func (s *SelectInput) SetOptions(options []string) {
	items := make([]list.Item, len(options))
	for i, option := range options {
		items[i] = collection(option)
	}
	s.list.SetItems(items)
	s.loading = false
}

func (s SelectInput) GetSelected() string {
	if s.loading || len(s.list.Items()) == 0 {
		return ""
	}
	if selectedItem := s.list.SelectedItem(); selectedItem != nil {
		return string(selectedItem.(collection))
	}
	return ""
}
