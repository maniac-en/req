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

type MainTab int

const (
	RequestBuilderMainTab MainTab = iota
	ResponseViewerMainTab
)

type SelectedCollectionView struct {
	layout           components.Layout
	endpointsManager *endpoints.EndpointsManager
	httpManager      *http.HTTPManager
	collection       collections.CollectionEntity
	sidebar          EndpointSidebarView
	selectedEndpoint *endpoints.EndpointEntity
	requestBuilder   RequestBuilder
	activeMainTab    MainTab
	width            int
	height           int
	notification     string
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
		requestBuilder:   NewRequestBuilder(),
		activeMainTab:    RequestBuilderMainTab,
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

	requestBuilder := NewRequestBuilder()
	requestBuilder.SetSize(innerWidth-sidebarWidth-1, innerHeight)

	return SelectedCollectionView{
		layout:           layout,
		endpointsManager: endpointsManager,
		httpManager:      httpManager,
		collection:       collection,
		sidebar:          sidebar,
		selectedEndpoint: nil,
		requestBuilder:   requestBuilder,
		activeMainTab:    RequestBuilderMainTab,
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
		v.requestBuilder.SetSize(innerWidth-sidebarWidth-1, innerHeight)

	case tea.KeyMsg:
		// Clear notification on any keypress
		v.notification = ""
		
		// If request builder is in component editing mode, only handle esc - forward everything else
		if v.activeMainTab == RequestBuilderMainTab && v.requestBuilder.IsEditingComponent() {
			if msg.String() == "esc" {
				// Forward the Esc to request builder to exit editing mode
				var builderCmd tea.Cmd
				v.requestBuilder, builderCmd = v.requestBuilder.Update(msg)
				return v, builderCmd
			}
			// Forward all other keys to request builder when editing
			var builderCmd tea.Cmd
			v.requestBuilder, builderCmd = v.requestBuilder.Update(msg)
			return v, builderCmd
		}
		
		// Normal key handling when not editing
		switch msg.String() {
		case "esc", "q":
			return v, func() tea.Msg { return BackToCollectionsMsg{} }
		case "1":
			v.activeMainTab = RequestBuilderMainTab
			v.requestBuilder.Focus()
		case "2":
			v.activeMainTab = ResponseViewerMainTab
			v.requestBuilder.Blur()
		case "a":
			v.notification = "Adding endpoints is not yet implemented"
			return v, nil
		case "r":
			v.notification = "Sending requests is not yet implemented"
			return v, nil
		}

	case EndpointSelectedMsg:
		// Store the selected endpoint for display
		v.selectedEndpoint = &msg.Endpoint
		v.requestBuilder.LoadFromEndpoint(msg.Endpoint)
		v.requestBuilder.Focus()

	case RequestSendMsg:
		return v, nil
	}

	// Forward messages to appropriate components (only if not editing)
	if !(v.activeMainTab == RequestBuilderMainTab && v.requestBuilder.IsEditingComponent()) {
		v.sidebar, cmd = v.sidebar.Update(msg)

		// Forward to request builder if it's the active tab
		if v.activeMainTab == RequestBuilderMainTab {
			var builderCmd tea.Cmd
			v.requestBuilder, builderCmd = v.requestBuilder.Update(msg)
			if builderCmd != nil {
				cmd = builderCmd
			}
		}
	}

	return v, cmd
}

func (v SelectedCollectionView) View() string {
	title := "Collection: " + v.collection.Name
	if v.selectedEndpoint != nil {
		title += " > " + v.selectedEndpoint.Name
	}

	sidebarContent := v.sidebar.View()

	// Main tab content
	var mainContent string
	if v.selectedEndpoint != nil {
		// Show main tabs
		tabsContent := v.renderMainTabs()
		tabContent := v.renderMainTabContent()
		mainContent = lipgloss.JoinVertical(lipgloss.Left, tabsContent, "", tabContent)
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

	instructions := "↑↓: navigate endpoints • 1: request • 2: response • enter: edit • esc: stop editing • r: send • esc/q: back"
	if v.notification != "" {
		instructions = v.notification
	}

	return v.layout.FullView(
		title,
		content,
		instructions,
	)
}

func (v SelectedCollectionView) renderMainTabs() string {
	tabs := []string{"Request Builder", "Response Viewer"}
	var renderedTabs []string

	for i, tab := range tabs {
		tabStyle := styles.ListItemStyle.Copy().
			Padding(0, 3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Secondary)

		if MainTab(i) == v.activeMainTab {
			tabStyle = tabStyle.
				Background(styles.Primary).
				Foreground(styles.TextPrimary).
				Bold(true)
		}

		renderedTabs = append(renderedTabs, tabStyle.Render(tab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (v SelectedCollectionView) renderMainTabContent() string {
	switch v.activeMainTab {
	case RequestBuilderMainTab:
		return v.requestBuilder.View()
	case ResponseViewerMainTab:
		return styles.ListItemStyle.Copy().
			Width(v.width/2).
			Height(v.height/2).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Response viewer coming soon...")
	default:
		return ""
	}
}

// Message types for selected collection view
type EndpointSelectedMsg struct {
	Endpoint endpoints.EndpointEntity
}
