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
	SelectedCollectionViewMode
)

type Model struct {
	ctx                    *Context
	mode                   ViewMode
	collectionsView        views.CollectionsView
	addCollectionView      views.AddCollectionView
	editCollectionView     views.EditCollectionView
	selectedCollectionView views.SelectedCollectionView
	width                  int
	height                 int
	selectedIndex          int
}

func NewModel(ctx *Context) Model {
	collectionsView := views.NewCollectionsView(ctx.Collections)
	if ctx.DummyDataCreated {
		collectionsView.SetDummyDataNotification(true)
	}
	
	m := Model{
		ctx:               ctx,
		mode:              CollectionsViewMode,
		collectionsView:   collectionsView,
		addCollectionView: views.NewAddCollectionView(ctx.Collections),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return m.collectionsView.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		isFiltering := m.mode == CollectionsViewMode && m.collectionsView.IsFiltering()

		if !isFiltering {
			switch msg.String() {
			case "ctrl+c", "q":
				if m.mode == CollectionsViewMode {
					return m, tea.Quit
				}
				m.mode = CollectionsViewMode
				return m, nil
			case "a":
				if m.mode == CollectionsViewMode {
					m.selectedIndex = m.collectionsView.GetSelectedIndex()
					m.mode = AddCollectionViewMode
					if m.width > 0 && m.height > 0 {
						sizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
						m.addCollectionView, _ = m.addCollectionView.Update(sizeMsg)
					}
					return m, nil
				}
			case "enter":
				if m.mode == CollectionsViewMode {
					if selectedItem := m.collectionsView.GetSelectedItem(); selectedItem != nil {
						m.selectedIndex = m.collectionsView.GetSelectedIndex()
						m.mode = SelectedCollectionViewMode
						if m.width > 0 && m.height > 0 {
							m.selectedCollectionView = views.NewSelectedCollectionViewWithSize(m.ctx.Endpoints, m.ctx.HTTP, *selectedItem, m.width, m.height)
						} else {
							m.selectedCollectionView = views.NewSelectedCollectionView(m.ctx.Endpoints, m.ctx.HTTP, *selectedItem)
						}
						return m, m.selectedCollectionView.Init()
					} else {
						log.Error("issue getting currently selected collection")
					}
				}
			case "e":
				if m.mode == CollectionsViewMode {
					if selectedItem := m.collectionsView.GetSelectedItem(); selectedItem != nil {
						m.selectedIndex = m.collectionsView.GetSelectedIndex()
						m.mode = EditCollectionViewMode
						m.editCollectionView = views.NewEditCollectionView(m.ctx.Collections, *selectedItem)
						if m.width > 0 && m.height > 0 {
							sizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
							m.editCollectionView, _ = m.editCollectionView.Update(sizeMsg)
						}
						return m, nil
					} else {
						log.Error("issue getting currently selected collection")
					}
				}
			case "x":
				if m.mode == CollectionsViewMode {
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
		if m.mode == CollectionsViewMode && !m.collectionsView.IsInitialized() {
			m.collectionsView = views.NewCollectionsViewWithSize(m.ctx.Collections, m.width, m.height)
			if m.ctx.DummyDataCreated {
				m.collectionsView.SetDummyDataNotification(true)
			}
			return m, m.collectionsView.Init()
		}
		if m.mode == CollectionsViewMode {
			m.collectionsView, _ = m.collectionsView.Update(msg)
		}
	case views.BackToCollectionsMsg:
		m.mode = CollectionsViewMode
		if m.width > 0 && m.height > 0 {
			m.collectionsView = views.NewCollectionsViewWithSize(m.ctx.Collections, m.width, m.height)
			if m.ctx.DummyDataCreated {
				m.collectionsView.SetDummyDataNotification(true)
			}
		}
		m.collectionsView.SetSelectedIndex(m.selectedIndex)
		return m, m.collectionsView.Init()
	case views.EditCollectionMsg:
		m.mode = EditCollectionViewMode
		m.editCollectionView = views.NewEditCollectionView(m.ctx.Collections, msg.Collection)
		return m, nil
	case views.CollectionDeletedMsg:
		return m, m.collectionsView.Init()
	case views.CollectionDeleteErrorMsg:
		return m, nil
	case views.CollectionCreatedMsg:
		m.addCollectionView.ClearForm()
		m.mode = CollectionsViewMode
		m.selectedIndex = 0
		m.collectionsView.SetSelectedIndex(m.selectedIndex)
		return m, m.collectionsView.Init()
	}

	switch m.mode {
	case CollectionsViewMode:
		m.collectionsView, cmd = m.collectionsView.Update(msg)
	case AddCollectionViewMode:
		m.addCollectionView, cmd = m.addCollectionView.Update(msg)
	case EditCollectionViewMode:
		m.editCollectionView, cmd = m.editCollectionView.Update(msg)
	case SelectedCollectionViewMode:
		m.selectedCollectionView, cmd = m.selectedCollectionView.Update(msg)
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
	case SelectedCollectionViewMode:
		return m.selectedCollectionView.View()
	default:
		return m.collectionsView.View()
	}
}
