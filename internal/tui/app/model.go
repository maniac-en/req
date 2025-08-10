package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/tui/keybinds"
	"github.com/maniac-en/req/internal/tui/styles"
	"github.com/maniac-en/req/internal/tui/views"
)

type ViewName string

const (
	Collections ViewName = "collections"
)

type AppModel struct {
	ctx         *Context
	width       int
	height      int
	Views       map[ViewName]views.ViewInterface
	focusedView ViewName
	keys        []key.Binding
	help        help.Model
}

func (a AppModel) Init() tea.Cmd {
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
		a.Views[a.focusedView], cmd = a.Views[a.focusedView].Update(tea.WindowSizeMsg{Height: a.AvailableHeight(), Width: msg.Width})
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.Keys.Quit):
			return a, tea.Quit
		}
	}

	a.Views[a.focusedView], cmd = a.Views[a.focusedView].Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a AppModel) View() string {
	footer := a.Footer()
	header := a.Header()
	view := a.Views[a.focusedView].View()
	help := a.Help()
	return lipgloss.JoinVertical(lipgloss.Top, header, view, help, footer)
}

func (a AppModel) Help() string {
	viewHelp := a.Views[a.focusedView].Help()
	appHelp := append(viewHelp, a.keys...)
	helpStruct := keybinds.Help{
		Keys: appHelp,
	}
	return styles.HelpStyle.Render(a.help.View(helpStruct))
}

func (a *AppModel) AvailableHeight() int {
	footer := a.Footer()
	header := a.Header()
	help := a.Help()
	return a.height - lipgloss.Height(header) - lipgloss.Height(footer) - lipgloss.Height(help)
}

func (a AppModel) Header() string {
	var b strings.Builder

	for key, value := range a.Views {
		if key == a.focusedView {
			b.WriteString(styles.TabHeadingActive.Render(value.Name()))
		} else {
			b.WriteString(styles.TabHeadingInactive.Render(value.Name()))
		}
	}
	b.WriteString(styles.TabHeadingInactive.Render(""))
	return b.String()
}

func (a AppModel) Footer() string {
	name := styles.ApplyGradientToFooter("REQ")
	footerText := styles.FooterSegmentStyle.Render(a.Views[a.focusedView].GetFooterSegment())
	version := styles.FooterVersionStyle.Width(a.width - lipgloss.Width(name) - lipgloss.Width(footerText)).Render("v0.1.0-alpha.2")
	return lipgloss.JoinHorizontal(lipgloss.Left, name, footerText, version)
}

func NewAppModel(ctx *Context) AppModel {
	appKeybinds := []key.Binding{
		keybinds.Keys.Quit,
	}

	model := AppModel{
		focusedView: Collections,
		ctx:         ctx,
		help:        help.New(),
		keys:        appKeybinds,
	}
	model.Views = map[ViewName]views.ViewInterface{
		Collections: views.NewCollectionsView(model.ctx.Collections),
	}
	return model
}
