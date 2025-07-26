package tabs

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/messages"
)

type collectionsOpts struct {
	options []OptionPair
}

type OptionPair struct {
	Label string
	Value string
}

type CollectionsTab struct {
	name     string
	selectUI SelectInput
	loaded   bool
}

func NewCollectionsTab() *CollectionsTab {
	return &CollectionsTab{
		name:     "Collections",
		selectUI: NewSelectInput(),
		loaded:   false,
	}
}

func (c *CollectionsTab) fetchOptions() tea.Cmd {
	// this is here for now to replicate what a db call would look like
	return tea.Tick(time.Millisecond*1000, func(time.Time) tea.Msg {
		return collectionsOpts{
			options: GlobalCollections,
		}
	})
}

func (c *CollectionsTab) Name() string {
	return c.name
}

func (c *CollectionsTab) Init() tea.Cmd {
	c.selectUI.Focus()
	return tea.Batch(
		c.selectUI.Init(),
		c.fetchOptions(),
	)
}

func (c *CollectionsTab) OnFocus() tea.Cmd {
	c.selectUI.Focus()

	return c.fetchOptions()
}

func (c *CollectionsTab) OnBlur() tea.Cmd {
	c.selectUI.Blur()
	return nil
}

func (c *CollectionsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case collectionsOpts:
		c.selectUI.SetOptions(msg.options)
		c.loaded = true

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			return c, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 1}
			}
		case "d":
			if selected := c.selectUI.GetSelected(); selected != "" {
				return c, c.deleteCollection(selected)
			}
		case "e": // Add edit key handling
			if selected := c.selectUI.GetSelected(); selected != "" {
				return c, c.editCollection(selected)
			}
		default:
			c.selectUI, cmd = c.selectUI.Update(msg)
		}

	default:
		c.selectUI, cmd = c.selectUI.Update(msg)
	}

	return c, cmd
}

func (c *CollectionsTab) View() string {

	if c.selectUI.IsLoading() {
		return c.selectUI.View()
	}

	selectContent := c.selectUI.View()

	style := lipgloss.NewStyle().PaddingRight(4)

	if !c.selectUI.IsLoading() && len(c.selectUI.list.Items()) > 0 {
		title := "Select Collection:\n\n"
		instructions := "\n ↑/k - up | ↓/j - down | / - search | + - add collection | enter - select | d - delete collection | e - edit collection"
		return title + style.Render(selectContent) + instructions
	}

	return style.Render(selectContent)

}

func (c *CollectionsTab) deleteCollection(value string) tea.Cmd {
	for i, collection := range GlobalCollections {
		if collection.Value == value {
			GlobalCollections = append(GlobalCollections[:i], GlobalCollections[i+1:]...)
			break
		}
	}
	return c.fetchOptions()
}

func (c *CollectionsTab) editCollection(value string) tea.Cmd {
	var label string
	for _, collection := range GlobalCollections {
		if collection.Value == value {
			label = collection.Label
			break
		}
	}

	return tea.Batch(
		func() tea.Msg {
			return messages.EditCollectionMsg{
				Label: label,
				Value: value,
			}
		},
		func() tea.Msg {
			return messages.SwitchTabMsg{TabIndex: 2}
		},
	)
}
