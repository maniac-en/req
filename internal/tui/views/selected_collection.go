package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/backend/http"
	"github.com/maniac-en/req/internal/tui/components"
	"github.com/maniac-en/req/internal/tui/styles"
)

type SelectedCollectionView struct {
	layout           components.Layout
	endpointsManager *endpoints.EndpointsManager
	httpManager      *http.HTTPManager
	collection       collections.CollectionEntity
	sidebar          EndpointSidebarView
	selectedEndpoint *endpoints.EndpointEntity
	width            int
	height           int
}

func NewSelectedCollectionView(endpointsManager *endpoints.EndpointsManager, httpManager *http.HTTPManager, collection collections.CollectionEntity) SelectedCollectionView {
	sidebar := NewEndpointSidebarView(endpointsManager, collection)
	sidebar.Focus() // Make sure sidebar starts focused
	
	return SelectedCollectionView{
		layout:           components.NewLayout(),
		endpointsManager: endpointsManager,
		httpManager:      httpManager,
		collection:       collection,
		sidebar:          sidebar,
		selectedEndpoint: nil,
	}
}

func NewSelectedCollectionViewWithSize(endpointsManager *endpoints.EndpointsManager, httpManager *http.HTTPManager, collection collections.CollectionEntity, width, height int) SelectedCollectionView {
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
	sidebar.Focus() // Make sure sidebar starts focused

	return SelectedCollectionView{
		layout:           layout,
		endpointsManager: endpointsManager,
		httpManager:      httpManager,
		collection:       collection,
		sidebar:          sidebar,
		selectedEndpoint: nil,
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

	case EndpointSelectedMsg:
		// Store the selected endpoint for display
		v.selectedEndpoint = &msg.Endpoint
	}

	// Always forward messages to sidebar for now, but it only handles them when focused
	v.sidebar, cmd = v.sidebar.Update(msg)

	return v, cmd
}


func (v SelectedCollectionView) View() string {
	title := "Collection: " + v.collection.Name
	if v.selectedEndpoint != nil {
		title += " > " + v.selectedEndpoint.Name
	}

	sidebarContent := v.sidebar.View()

	// Simple endpoint information display
	var mainContent string
	if v.selectedEndpoint != nil {
		var lines []string
		lines = append(lines, "Selected Endpoint:")
		lines = append(lines, "")
		lines = append(lines, "Name: "+v.selectedEndpoint.Name)
		lines = append(lines, "Method: "+v.selectedEndpoint.Method)
		lines = append(lines, "URL: "+v.selectedEndpoint.Url)
		
		if v.selectedEndpoint.Headers != "" {
			lines = append(lines, "")
			lines = append(lines, "Headers: "+v.selectedEndpoint.Headers)
		}
		
		if v.selectedEndpoint.QueryParams != "" {
			lines = append(lines, "")
			lines = append(lines, "Query Params: "+v.selectedEndpoint.QueryParams)
		}
		
		if v.selectedEndpoint.RequestBody != "" {
			lines = append(lines, "")
			lines = append(lines, "Request Body:")
			lines = append(lines, v.selectedEndpoint.RequestBody)
		}
		
		mainContent = lipgloss.JoinVertical(lipgloss.Left, lines...)
	} else {
		// Check if there are no endpoints at all
		if len(v.sidebar.endpoints) == 0 {
			mainContent = "Create an endpoint to get started"
		} else {
			mainContent = "Select an endpoint from the sidebar to view details"
		}
	}

	if v.width < 10 || v.height < 10 {
		return v.layout.FullView(title, sidebarContent, "esc/q: back to collections")
	}

	windowWidth := int(float64(v.width) * 0.85)
	windowHeight := int(float64(v.height) * 0.8)
	innerWidth := windowWidth
	innerHeight := windowHeight - 6

	sidebarWidth := innerWidth / 4
	mainWidth := innerWidth - sidebarWidth - 1

	// Sidebar styling
	sidebarStyle := styles.SidebarStyle.Copy().
		Width(sidebarWidth).
		Height(innerHeight).
		BorderForeground(styles.Primary)

	mainStyle := styles.MainContentStyle.Copy().
		Width(mainWidth).
		Height(innerHeight)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarStyle.Render(sidebarContent),
		mainStyle.Render(mainContent),
	)

	instructions := "↑↓: navigate endpoints • esc/q: back"

	return v.layout.FullView(
		title,
		content,
		instructions,
	)
}

// Message types for selected collection view
type EndpointSelectedMsg struct {
	Endpoint endpoints.EndpointEntity
}
