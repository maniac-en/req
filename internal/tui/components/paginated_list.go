package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/tui/styles"
)

type ListItem interface {
	list.Item
	GetID() string
	GetTitle() string
	GetDescription() string
}

type PaginatedList struct {
	list   list.Model
	width  int
	height int
}

func NewPaginatedList(items []ListItem, title string) PaginatedList {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	const defaultWidth = 20
	const defaultHeight = 14

	l := list.New(listItems, paginatedItemDelegate{}, defaultWidth, defaultHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true) // Enable filtering
	l.SetShowHelp(false) // Disable built-in help text
	l.Styles.Title = styles.TitleStyle

	return PaginatedList{
		list: l,
	}
}

func (pl *PaginatedList) SetSize(width, height int) {
	pl.width = width
	pl.height = height
	
	// Safety check to prevent nil pointer dereference
	if width > 0 && height > 0 {
		pl.list.SetWidth(width)
		pl.list.SetHeight(height)
	}
}

func (pl PaginatedList) Init() tea.Cmd {
	return nil
}

func (pl PaginatedList) Update(msg tea.Msg) (PaginatedList, tea.Cmd) {
	newListModel, cmd := pl.list.Update(msg)
	pl.list = newListModel
	return pl, cmd
}

func (pl PaginatedList) View() string {
	return pl.list.View()
}

func (pl PaginatedList) SelectedItem() ListItem {
	if selectedItem := pl.list.SelectedItem(); selectedItem != nil {
		if listItem, ok := selectedItem.(ListItem); ok {
			return listItem
		}
	}
	return nil
}

func (pl PaginatedList) SelectedIndex() int {
	return pl.list.Index()
}

func (pl *PaginatedList) SetIndex(i int) {
	pl.list.Select(i)
}

func (pl PaginatedList) IsFiltering() bool {
	return pl.list.FilterState() == list.Filtering
}

type paginatedItemDelegate struct{}

func (d paginatedItemDelegate) Height() int                             { return 1 }
func (d paginatedItemDelegate) Spacing() int                            { return 0 }
func (d paginatedItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d paginatedItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if i, ok := item.(ListItem); ok {
		str := i.GetTitle()

		fn := styles.ListItemStyle.Render
		if index == m.Index() {
			fn = func(s ...string) string {
				return styles.SelectedListItemStyle.Render("> " + s[0])
			}
		}

		fmt.Fprint(w, fn(str))
	}
}