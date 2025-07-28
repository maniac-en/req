package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/tui/components"
	"github.com/maniac-en/req/internal/tui/styles"
)

type SelectedCollectionView struct {
	layout           components.Layout
	endpointsManager *endpoints.EndpointsManager
	collection       collections.CollectionEntity
	sidebar          EndpointSidebarView
	width            int
	height           int
}

func NewSelectedCollectionView(endpointsManager *endpoints.EndpointsManager, collection collections.CollectionEntity) SelectedCollectionView {
	return SelectedCollectionView{
		layout:           components.NewLayout(),
		endpointsManager: endpointsManager,
		collection:       collection,
		sidebar:          NewEndpointSidebarView(endpointsManager, collection),
	}
}

func (v SelectedCollectionView) Init() tea.Cmd {
	return v.sidebar.Init()
}

func (v SelectedCollectionView) Update(msg tea.Msg) (SelectedCollectionView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		v.layout.SetSize(v.width, v.height)

		sidebarWidth := v.width / 3
		v.sidebar.width = sidebarWidth
		v.sidebar.height = v.height - 4


	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return v, func() tea.Msg { return BackToCollectionsMsg{} }
		}
	}


	// Forward messages to sidebar
	v.sidebar, cmd = v.sidebar.Update(msg)

	return v, cmd
}

func (v SelectedCollectionView) View() string {
	title := "Collection: " + v.collection.Name

	sidebarContent := v.sidebar.View()
	mainContent := "Endpoint details will be displayed here"
	sidebarStyle := styles.SidebarStyle.Copy().
		Width(v.width / 3).
		Height(v.height - 4)

	mainStyle := styles.MainContentStyle.Copy().
		Width((v.width * 2 / 3) - 1).
		Height(v.height - 4)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarStyle.Render(sidebarContent),
		mainStyle.Render(mainContent),
	)

	instructions := "↑↓: navigate endpoints • esc/q: back to collections"

	return v.layout.FullView(
		title,
		content,
		instructions,
	)
}
