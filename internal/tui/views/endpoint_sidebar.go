package views

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/tui/components"
)

type EndpointSidebarView struct {
	list             components.PaginatedList
	endpointsManager *endpoints.EndpointsManager
	collection       collections.CollectionEntity
	width            int
	height           int
	initialized      bool
	selectedIndex    int
	endpoints        []endpoints.EndpointEntity
	focused          bool
}

func NewEndpointSidebarView(endpointsManager *endpoints.EndpointsManager, collection collections.CollectionEntity) EndpointSidebarView {
	return EndpointSidebarView{
		endpointsManager: endpointsManager,
		collection:       collection,
		selectedIndex:    0,
		focused:          false,
	}
}

func (v *EndpointSidebarView) Focus() {
	v.focused = true
}

func (v *EndpointSidebarView) Blur() {
	v.focused = false
}

func (v EndpointSidebarView) Focused() bool {
	return v.focused
}

func (v EndpointSidebarView) Init() tea.Cmd {
	return v.loadEndpoints
}

func (v *EndpointSidebarView) loadEndpoints() tea.Msg {
	result, err := v.endpointsManager.ListByCollection(context.Background(), v.collection.ID, 100, 0)
	if err != nil {
		return endpointsLoadError{err: err}
	}
	return endpointsLoaded{
		endpoints: result.Endpoints,
	}
}

func (v EndpointSidebarView) Update(msg tea.Msg) (EndpointSidebarView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		if v.initialized {
			v.list.SetSize(v.width, v.height)
		}

	case endpointsLoaded:
		v.endpoints = msg.endpoints
		items := make([]components.ListItem, len(msg.endpoints))
		for i, endpoint := range msg.endpoints {
			items[i] = components.NewEndpointItem(endpoint)
		}

		title := fmt.Sprintf("Endpoints (%d)", len(msg.endpoints))
		v.list = components.NewPaginatedList(items, title)
		v.list.SetIndex(v.selectedIndex)

		if v.width > 0 && v.height > 0 {
			v.list.SetSize(v.width, v.height)
		}
		v.initialized = true

		// Auto-select first endpoint if available
		if len(msg.endpoints) > 0 {
			return v, func() tea.Msg {
				return EndpointSelectedMsg{Endpoint: msg.endpoints[0]}
			}
		}

	case endpointsLoadError:
		v.initialized = true

	case tea.KeyMsg:
		if v.initialized {
			switch msg.String() {
			// case "enter":
			// 	if selectedEndpoint := v.GetSelectedEndpoint(); selectedEndpoint != nil {
			// 		return v, func() tea.Msg {
			// 			return EndpointSelectedMsg{Endpoint: *selectedEndpoint}
			// 		}
			// 	}
			default:
				// Forward navigation keys to the list even if not explicitly focused
				oldIndex := v.list.SelectedIndex()
				v.list, cmd = v.list.Update(msg)
				newIndex := v.list.SelectedIndex()

				// If the selected index changed, auto-select the new endpoint
				if oldIndex != newIndex && newIndex >= 0 && newIndex < len(v.endpoints) {
					return v, func() tea.Msg {
						return EndpointSelectedMsg{Endpoint: v.endpoints[newIndex]}
					}
				}
			}
		}
	}

	return v, cmd
}

func (v EndpointSidebarView) GetSelectedEndpoint() *endpoints.EndpointEntity {
	if !v.initialized || len(v.endpoints) == 0 {
		return nil
	}

	selectedIndex := v.list.SelectedIndex()
	if selectedIndex >= 0 && selectedIndex < len(v.endpoints) {
		return &v.endpoints[selectedIndex]
	}
	return nil
}

func (v EndpointSidebarView) GetSelectedIndex() int {
	if v.initialized {
		return v.list.SelectedIndex()
	}
	return v.selectedIndex
}

func (v *EndpointSidebarView) SetSelectedIndex(index int) {
	v.selectedIndex = index
	if v.initialized {
		v.list.SetIndex(index)
	}
}

func (v EndpointSidebarView) View() string {
	if !v.initialized {
		title := "Endpoints"
		content := "Loading endpoints..."
		return v.formatEmptyState(title, content)
	}
	if len(v.endpoints) == 0 {
		title := "Endpoints (0)"
		content := "No endpoints found"
		return v.formatEmptyState(title, content)
	}
	return v.list.View()
}

func (v EndpointSidebarView) formatEmptyState(title, content string) string {
	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")
	lines = append(lines, content)

	for len(lines) < v.height-2 {
		lines = append(lines, "")
	}

	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return result
}

type endpointsLoaded struct {
	endpoints []endpoints.EndpointEntity
}

type endpointsLoadError struct {
	err error
}
