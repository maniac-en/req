package tabs

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type collectionsOpts struct {
	options []string
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
			options: []string{
				"Collection 1",
				"Collection 2",
				"Collection 3",
				"Collection 4",
			},
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
	if !c.loaded {
		return c.fetchOptions()
	}
	return nil
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
	default:
		c.selectUI, cmd = c.selectUI.Update(msg)
	}
	return c, cmd
}

func (c *CollectionsTab) View() string {
	content := "Select Collection:\n\n" + c.selectUI.View()

	return content
}
