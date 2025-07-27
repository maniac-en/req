package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/tui/views"
)

type ViewMode int

const (
	CollectionsViewMode ViewMode = iota
	AddCollectionViewMode
	EditCollectionViewMode
)

type Model struct {
	ctx             *Context
	mode            ViewMode
	collectionsView views.CollectionsView
	addCollectionView views.AddCollectionView
	width           int
	height          int
}

func NewModel(ctx *Context) Model {
	return Model{
		ctx:               ctx,
		mode:              CollectionsViewMode,
		collectionsView:   views.NewCollectionsView(ctx.Collections),
		addCollectionView: views.NewAddCollectionView(ctx.Collections),
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
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case views.BackToCollectionsMsg:
		m.mode = CollectionsViewMode
		// Reload collections to show any changes
		return m, m.collectionsView.Init()
	}

	// Forward messages to the appropriate view
	switch m.mode {
	case CollectionsViewMode:
		m.collectionsView, cmd = m.collectionsView.Update(msg)
	case AddCollectionViewMode:
		m.addCollectionView, cmd = m.addCollectionView.Update(msg)
	}
	
	return m, cmd
}

func (m Model) View() string {
	switch m.mode {
	case CollectionsViewMode:
		return m.collectionsView.View()
	case AddCollectionViewMode:
		return m.addCollectionView.View()
	default:
		return m.collectionsView.View()
	}
}

