package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/tui/components"
)

type SelectedCollectionView struct {
	layout           components.Layout
	endpointsManager *endpoints.EndpointsManager
	collection       collections.CollectionEntity
	width            int
	height           int
	initialized      bool
}

func NewSelectedCollectionView(endpointsManager *endpoints.EndpointsManager, collection collections.CollectionEntity) SelectedCollectionView {
	return SelectedCollectionView{
		layout:           components.NewLayout(),
		endpointsManager: endpointsManager,
		collection:       collection,
	}
}

func (v SelectedCollectionView) Init() tea.Cmd {
	return nil
}

func (v SelectedCollectionView) Update(msg tea.Msg) (SelectedCollectionView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		v.layout.SetSize(v.width, v.height)
		v.initialized = true

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return v, func() tea.Msg { return BackToCollectionsMsg{} }
		}
	}

	return v, cmd
}

func (v SelectedCollectionView) View() string {
	if !v.initialized {
		return v.layout.FullView(
			"Loading...",
			"Initializing collection view...",
			"Please wait",
		)
	}

	title := "Collection: " + v.collection.Name
	content := "Selected collection view - endpoints will be displayed here"
	instructions := "esc/q: back to collections"

	return v.layout.FullView(
		title,
		content,
		instructions,
	)
}
