package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/log"
	"github.com/maniac-en/req/internal/tui/views"
)

type ViewMode int

const (
	CollectionsViewMode ViewMode = iota
	AddCollectionViewMode
	EditCollectionViewMode
)

type Model struct {
	ctx                *Context
	mode               ViewMode
	collectionsView    views.CollectionsView
	addCollectionView  views.AddCollectionView
	editCollectionView views.EditCollectionView
	width              int
	height             int
}

func NewModel(ctx *Context) Model {
	return Model{
		ctx:               ctx,
		mode:              CollectionsViewMode,
		collectionsView:   views.NewCollectionsView(ctx.Collections),
		addCollectionView: views.NewAddCollectionView(ctx.Collections),
		// editCollectionView will be created on demand
	}
}

func (m Model) Init() tea.Cmd {
	return m.collectionsView.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keybinds only when not in filtering mode
		isFiltering := m.mode == CollectionsViewMode && m.collectionsView.IsFiltering()

		if !isFiltering {
			switch msg.String() {
			case "ctrl+c", "q":
				if m.mode == CollectionsViewMode {
					return m, tea.Quit
				}
				// For other views, 'q' goes back to collections
				m.mode = CollectionsViewMode
				return m, nil
			case "a":
				if m.mode == CollectionsViewMode {
					m.mode = AddCollectionViewMode
					return m, nil
				}
			case "e":
				if m.mode == CollectionsViewMode {
					// Get selected collection and switch to edit mode
					if selectedItem := m.collectionsView.GetSelectedItem(); selectedItem != nil {
						m.mode = EditCollectionViewMode
						m.editCollectionView = views.NewEditCollectionView(m.ctx.Collections, *selectedItem)
						return m, nil
					} else {
						log.Error("issue getting currently selected collection")
					}
				}
			case "x":
				if m.mode == CollectionsViewMode {
					// Delete selected collection
					if selectedItem := m.collectionsView.GetSelectedItem(); selectedItem != nil {
						return m, func() tea.Msg {
							err := m.ctx.Collections.Delete(context.Background(), selectedItem.ID)
							if err != nil {
								return views.CollectionDeleteErrorMsg{Err: err}
							}
							return views.CollectionDeletedMsg{ID: selectedItem.ID}
						}
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case views.BackToCollectionsMsg:
		m.mode = CollectionsViewMode
		// Reload collections to show any changes
		return m, m.collectionsView.Init()
	case views.EditCollectionMsg:
		m.mode = EditCollectionViewMode
		m.editCollectionView = views.NewEditCollectionView(m.ctx.Collections, msg.Collection)
		return m, nil
	case views.CollectionDeletedMsg:
		// Collection deleted, reload collections view
		return m, m.collectionsView.Init()
	case views.CollectionDeleteErrorMsg:
		// Delete failed, just continue
		return m, nil
	}

	// Forward messages to the appropriate view
	switch m.mode {
	case CollectionsViewMode:
		m.collectionsView, cmd = m.collectionsView.Update(msg)
	case AddCollectionViewMode:
		m.addCollectionView, cmd = m.addCollectionView.Update(msg)
	case EditCollectionViewMode:
		m.editCollectionView, cmd = m.editCollectionView.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.mode {
	case CollectionsViewMode:
		return m.collectionsView.View()
	case AddCollectionViewMode:
		return m.addCollectionView.View()
	case EditCollectionViewMode:
		return m.editCollectionView.View()
	default:
		return m.collectionsView.View()
	}
}
