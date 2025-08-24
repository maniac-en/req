package app

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/log"
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
	"github.com/maniac-en/req/internal/tui/keybinds"
	"github.com/maniac-en/req/internal/tui/messages"
	"github.com/maniac-en/req/internal/tui/styles"
	"github.com/maniac-en/req/internal/tui/views"
)

type ViewName string

const (
	Collections ViewName = "collections"
	Endpoints   ViewName = "endpoints"
)

type Heading struct {
	name  string
	order int
}

type AppModel struct {
	ctx         *Context
	width       int
	height      int
	Views       map[ViewName]views.ViewInterface
	focusedView ViewName
	keys        []key.Binding
	help        help.Model
	errorMsg    string
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
		for key, _ := range a.Views {
			a.Views[key], cmd = a.Views[key].Update(tea.WindowSizeMsg{Height: a.AvailableHeight(), Width: msg.Width})
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	case messages.ChooseItem[optionsProvider.Option]:
		switch msg.Source {
		case "collections":
			return a, func() tea.Msg {
				return messages.NavigateToView{
					ViewName: string(Endpoints),
					Data:     msg.Item,
				}
			}
		}
	case messages.NavigateToView:
		a.Views[a.focusedView].OnBlur()

		if msg.Data != nil {
			err := a.Views[ViewName(msg.ViewName)].SetState(msg.Data)
			if err != nil {
				log.Error("failed to set view state during navigation", "target_view", msg.ViewName, "error", err)
				return a, nil
			}
		}

		a.focusedView = ViewName(msg.ViewName)
		a.Views[a.focusedView].OnFocus()
		return a, nil
	case messages.ShowError:
		log.Error("user operation failed", "error", msg.Message)
		a.errorMsg = msg.Message
		return a, nil
	case tea.KeyMsg:
		a.errorMsg = ""
		switch {
		case key.Matches(msg, keybinds.Keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, keybinds.Keys.Back):
			if a.focusedView == Endpoints {
				return a, func() tea.Msg {
					return messages.NavigateToView{
						ViewName: string(Collections),
						Data:     nil,
					}
				}
			}
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

	if a.errorMsg != "" {
		errorBar := styles.ErrorBarStyle.Width(a.width).Render("Error: " + a.errorMsg)
		return lipgloss.JoinVertical(lipgloss.Top, header, view, errorBar, help, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Top, header, view, help, footer)
}

func (a AppModel) Help() string {
	viewHelp := a.Views[a.focusedView].Help()

	var appHelp []key.Binding
	appHelp = append(appHelp, a.keys...)

	if a.focusedView == Endpoints {
		appHelp = append(appHelp, keybinds.Keys.Back)
	}

	allHelp := append(viewHelp, appHelp...)
	helpStruct := keybinds.Help{
		Keys: allHelp,
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

	// INFO: this might be a bit messy, could be a nice idea to look into OrderedMaps maybe?
	views := []Heading{}
	for key := range a.Views {
		views = append(views, Heading{
			name:  a.Views[key].Name(),
			order: a.Views[key].Order(),
		})
	}
	sort.Slice(views, func(i, j int) bool {
		return views[i].order < views[j].order
	})

	for _, item := range views {
		if item.name == a.Views[a.focusedView].Name() {
			b.WriteString(styles.TabHeadingActive.Render(item.name))
		} else {
			b.WriteString(styles.TabHeadingInactive.Render(item.name))
		}
	}

	b.WriteString(styles.TabHeadingInactive.Render(""))

	return b.String()
}

func (a AppModel) Footer() string {
	name := styles.ApplyGradientToFooter("REQ")
	footerText := styles.FooterSegmentStyle.Render(a.Views[a.focusedView].GetFooterSegment())
	version := styles.FooterVersionStyle.Width(a.width - lipgloss.Width(name) - lipgloss.Width(footerText)).Render(a.ctx.Version)
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
		Collections: views.NewCollectionsView(model.ctx.Collections, model.ctx.Endpoints, 1),
		Endpoints:   views.NewEndpointsView(model.ctx.Endpoints, 2),
	}
	return model
}
