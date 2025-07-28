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

func NewSelectedCollectionViewWithSize(endpointsManager *endpoints.EndpointsManager, collection collections.CollectionEntity, width, height int) SelectedCollectionView {
	layout := components.NewLayout()
	layout.SetSize(width, height)

	windowWidth := int(float64(width) * 0.85)
	windowHeight := int(float64(height) * 0.8)
	innerWidth := windowWidth - 4
	innerHeight := windowHeight - 6
	sidebarWidth := innerWidth / 4

	sidebar := NewEndpointSidebarView(endpointsManager, collection)
	sidebar.width = sidebarWidth
	sidebar.height = innerHeight

	return SelectedCollectionView{
		layout:           layout,
		endpointsManager: endpointsManager,
		collection:       collection,
		sidebar:          sidebar,
		width:            width,
		height:           height,
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

		windowWidth := int(float64(v.width) * 0.85)
		windowHeight := int(float64(v.height) * 0.8)
		innerWidth := windowWidth - 4
		innerHeight := windowHeight - 6
		sidebarWidth := innerWidth / 4

		v.sidebar.width = sidebarWidth
		v.sidebar.height = innerHeight

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return v, func() tea.Msg { return BackToCollectionsMsg{} }
		}
	}

	v.sidebar, cmd = v.sidebar.Update(msg)

	return v, cmd
}

func (v SelectedCollectionView) View() string {
	title := "Collection: " + v.collection.Name

	sidebarContent := v.sidebar.View()
	mainContent := "Endpoint details will be displayed here"

	if v.width < 10 || v.height < 10 {
		return v.layout.FullView(title, sidebarContent, "esc/q: back to collections")
	}

	windowWidth := int(float64(v.width) * 0.85)
	windowHeight := int(float64(v.height) * 0.8)
	innerWidth := windowWidth
	innerHeight := windowHeight - 6

	sidebarWidth := innerWidth / 4
	mainWidth := innerWidth - sidebarWidth - 1

	sidebarStyle := styles.SidebarStyle.Copy().
		Width(sidebarWidth).
		Height(innerHeight)

	mainStyle := styles.MainContentStyle.Copy().
		Width(mainWidth).
		Height(innerHeight).
		Align(lipgloss.Center, lipgloss.Center)

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
