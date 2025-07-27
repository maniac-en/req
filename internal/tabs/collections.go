package tabs

import (
	"context"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/messages"
)

type collectionsOpts struct {
	options    []OptionPair
	totalItems int
}

type OptionPair struct {
	Label string
	Value string
}

type CollectionsTab struct {
	name             string
	selectUI         SelectInput
	loaded           bool
	currentPage      int
	itemsPerPage     int
	totalCollections int
	globalState      *global.State
}

func NewCollectionsTab(state *global.State) *CollectionsTab {
	itemsPerPage := 5
	return &CollectionsTab{
		name:         "Collections",
		selectUI:     NewSelectInput(),
		loaded:       false,
		currentPage:  0,
		itemsPerPage: itemsPerPage,
		globalState:  state,
	}
}

func (c *CollectionsTab) IsFiltering() bool {
	return c.selectUI.list.FilterState() == list.Filtering
}

func (c *CollectionsTab) fetchOptions(limit, offset int) tea.Cmd {
	ctx := global.GetAppContext()
	paginatedCollections, err := ctx.Collections.ListPaginated(context.Background(), limit, offset)
	if err != nil {
	}
	options := []OptionPair{}
	for i, _ := range paginatedCollections.Collections {
		options = append(options, OptionPair{
			Label: paginatedCollections.Collections[i].GetName(),
			Value: strconv.FormatInt(paginatedCollections.Collections[i].GetID(), 10),
		})
	}
	GlobalCollections = options
	c.totalCollections = int(paginatedCollections.Total)
	return func() tea.Msg {
		return collectionsOpts{
			options:    GlobalCollections,
			totalItems: c.totalCollections,
		}
	}
}

func (c *CollectionsTab) Name() string {
	return c.name
}

func (c *CollectionsTab) Instructions() string {
	return "\n k - up • j - down • / - search • + - add collection • enter - select • d - delete collection • e - edit collection • h - prev page • l - next page"
}

func (c *CollectionsTab) Init() tea.Cmd {
	c.selectUI.Focus()
	return tea.Batch(
		c.selectUI.Init(),
		c.fetchOptions(c.itemsPerPage, 0),
	)
}

func (c *CollectionsTab) OnFocus() tea.Cmd {
	c.selectUI.Focus()

	return c.fetchOptions(c.itemsPerPage, c.currentPage*c.itemsPerPage)
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
		// Check if list is filtering otherwise the keybinds wouldn't let us type
		if c.IsFiltering() {
			c.selectUI, cmd = c.selectUI.Update(msg)
			return c, cmd
		}

		switch msg.String() {
		case "+":
			return c, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 1}
			}
		case "d":
			if selected := c.selectUI.GetSelected(); selected != "" {
				return c, c.deleteCollection(selected)
			}
		case "e":
			if selected := c.selectUI.GetSelected(); selected != "" {
				return c, c.editCollection(selected)
			}
		case "h":
			if c.currentPage > 0 {
				c.currentPage--
				newOffset := c.currentPage * c.itemsPerPage
				return c, c.fetchOptions(c.itemsPerPage, newOffset)
			}
		case "l":
			totalPages := (c.totalCollections + c.itemsPerPage - 1) / c.itemsPerPage
			if c.currentPage < totalPages-1 {
				c.currentPage++
				newOffset := c.currentPage * c.itemsPerPage
				return c, c.fetchOptions(c.itemsPerPage, newOffset)
			}
		case "enter":
			c.globalState.SetCurrentCollection(c.selectUI.GetSelected())
			return c, func() tea.Msg {
				return messages.SwitchTabMsg{TabIndex: 3}
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

	style := lipgloss.NewStyle().
		PaddingRight(4)

	if !c.selectUI.IsLoading() && len(c.selectUI.list.Items()) > 0 {
		title := "\n\n\n\n\n\n\nSelect Collection:\n\n"
		return title + style.Render(selectContent)
	}

	return style.Render(selectContent)

}

func (c *CollectionsTab) deleteCollection(value string) tea.Cmd {
	ctx := global.GetAppContext()
	id, _ := strconv.Atoi(value)
	err := ctx.Collections.Delete(context.Background(), int64(id))
	if err != nil {
		return c.fetchOptions(c.itemsPerPage, c.currentPage*c.itemsPerPage)
	}
	for i, collection := range GlobalCollections {
		if collection.Value == value {
			GlobalCollections = append(GlobalCollections[:i], GlobalCollections[i+1:]...)
			break
		}
	}
	return c.fetchOptions(c.itemsPerPage, c.currentPage*c.itemsPerPage)
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
