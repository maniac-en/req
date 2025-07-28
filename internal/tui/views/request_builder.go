package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/tui/components"
	"github.com/maniac-en/req/internal/tui/styles"
)

type RequestBuilderTab int

const (
	RequestBodyTab RequestBuilderTab = iota
	HeadersTab
	QueryParamsTab
)

type RequestBuilder struct {
	endpoint         *endpoints.EndpointEntity
	method           string
	url              string
	requestBody      string
	activeTab        RequestBuilderTab
	bodyTextarea     components.Textarea
	width            int
	height           int
	focused          bool
	componentFocused bool // Whether we're actually editing a component
}

func NewRequestBuilder() RequestBuilder {
	bodyTextarea := components.NewTextarea("Body", "Enter request body (JSON, text, etc.)")

	return RequestBuilder{
		method:           "GET",
		url:              "",
		requestBody:      "",
		activeTab:        RequestBodyTab,
		bodyTextarea:     bodyTextarea,
		focused:          false,
		componentFocused: false,
	}
}

func (rb *RequestBuilder) SetSize(width, height int) {
	rb.width = width
	rb.height = height

	// Set size for body textarea (use most of available width)
	// Use about 90% of available width for better JSON editing
	textareaWidth := int(float64(width) * 0.9)
	if textareaWidth > 120 {
		textareaWidth = 120 // Cap at reasonable max width
	}
	if textareaWidth < 60 {
		textareaWidth = 60 // Ensure minimum usable width
	}
	
	// Set height for textarea (leave space for method/URL, tabs)
	textareaHeight := height - 8 // Account for method/URL row + tabs + spacing
	if textareaHeight < 5 {
		textareaHeight = 5
	}
	if textareaHeight > 15 {
		textareaHeight = 15 // Cap at reasonable height
	}
	
	rb.bodyTextarea.SetSize(textareaWidth, textareaHeight)
}

func (rb *RequestBuilder) Focus() {
	rb.focused = true
	// Don't auto-focus any component - user needs to explicitly focus in
	rb.componentFocused = false
	rb.bodyTextarea.Blur()
}

func (rb *RequestBuilder) Blur() {
	rb.focused = false
	rb.componentFocused = false
	rb.bodyTextarea.Blur()
}

func (rb RequestBuilder) Focused() bool {
	return rb.focused
}

func (rb RequestBuilder) IsEditingComponent() bool {
	return rb.componentFocused
}

func (rb *RequestBuilder) LoadFromEndpoint(endpoint endpoints.EndpointEntity) {
	rb.endpoint = &endpoint
	rb.method = endpoint.Method
	rb.url = endpoint.Url
	rb.requestBody = endpoint.RequestBody
	rb.bodyTextarea.SetValue(endpoint.RequestBody)
	rb.bodyTextarea.Blur() // Make sure it's not focused by default
	rb.componentFocused = false
}

func (rb RequestBuilder) Update(msg tea.Msg) (RequestBuilder, tea.Cmd) {
	if !rb.focused {
		return rb, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			// Only handle tab switching if not editing a component
			if !rb.componentFocused {
				if msg.String() == "tab" {
					rb.activeTab = (rb.activeTab + 1) % 3
				} else {
					rb.activeTab = (rb.activeTab + 2) % 3 // Go backwards
				}
			}
		case "enter":
			if !rb.componentFocused {
				// Focus into the current tab's component for editing
				rb.componentFocused = true
				switch rb.activeTab {
				case RequestBodyTab:
					rb.bodyTextarea.Focus()
				case HeadersTab:
					// TODO: Focus headers editor
				case QueryParamsTab:
					// TODO: Focus query params editor
				}
			}
		// case "r":
		// 	// Send request - only when not editing a component
		// 	if !rb.componentFocused {
		// 		return rb, func() tea.Msg {
		// 			return RequestSendMsg{Method: rb.method, URL: rb.url, Body: rb.bodyInput.Value()}
		// 		}
		// 	}
		case "esc":
			// Exit component editing mode
			if rb.componentFocused {
				rb.componentFocused = false
				rb.bodyTextarea.Blur()
				// TODO: Blur other components
			}
		}
	}

	// Only update components if we're in component editing mode
	if rb.componentFocused {
		switch rb.activeTab {
		case RequestBodyTab:
			rb.bodyTextarea, cmd = rb.bodyTextarea.Update(msg)
		case HeadersTab:
			// TODO: Update headers editor
		case QueryParamsTab:
			// TODO: Update query params editor
		}
	}

	return rb, cmd
}

func (rb RequestBuilder) View() string {
	if rb.width < 10 || rb.height < 10 {
		return "Request Builder (resize window)"
	}

	var sections []string

	// Method and URL row - aligned properly
	methodStyle := styles.ListItemStyle.Copy().
		Background(styles.Primary).
		Foreground(styles.TextPrimary).
		Padding(0, 2).
		Bold(true).
		Height(1)

	urlStyle := styles.ListItemStyle.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(0, 2).
		Width(rb.width - 20).
		Height(1)

	methodView := methodStyle.Render(rb.method)
	urlView := urlStyle.Render(rb.url)
	methodUrlRow := lipgloss.JoinHorizontal(lipgloss.Center, methodView, " ", urlView)
	sections = append(sections, methodUrlRow, "")

	// Tab headers
	tabHeaders := rb.renderTabHeaders()
	sections = append(sections, tabHeaders, "")

	// Tab content
	tabContent := rb.renderTabContent()
	sections = append(sections, tabContent)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (rb RequestBuilder) renderTabHeaders() string {
	tabs := []string{"Request Body", "Headers", "Query Params"}
	var renderedTabs []string

	for i, tab := range tabs {
		tabStyle := styles.ListItemStyle.Copy().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Secondary)

		if RequestBuilderTab(i) == rb.activeTab {
			tabStyle = tabStyle.
				Background(styles.Primary).
				Foreground(styles.TextPrimary).
				Bold(true)
		}

		renderedTabs = append(renderedTabs, tabStyle.Render(tab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (rb RequestBuilder) renderTabContent() string {
	switch rb.activeTab {
	case RequestBodyTab:
		return rb.bodyTextarea.View()
	case HeadersTab:
		return rb.renderPlaceholderTab("Headers editor coming soon...")
	case QueryParamsTab:
		return rb.renderPlaceholderTab("Query params editor coming soon...")
	default:
		return ""
	}
}

func (rb RequestBuilder) renderPlaceholderTab(message string) string {
	// Calculate the same dimensions as the textarea
	textareaWidth := int(float64(rb.width) * 0.9)
	if textareaWidth > 120 {
		textareaWidth = 120
	}
	if textareaWidth < 60 {
		textareaWidth = 60
	}
	
	textareaHeight := rb.height - 8
	if textareaHeight < 5 {
		textareaHeight = 5
	}
	if textareaHeight > 15 {
		textareaHeight = 15
	}

	// Create a placeholder with the same structure as textarea
	labelStyle := styles.TitleStyle.Copy().
		Width(12).
		Align(lipgloss.Right)

	containerWidth := textareaWidth - 12 - 1 - 2 // Same calculation as textarea
	if containerWidth < 20 {
		containerWidth = 20
	}

	container := styles.ListItemStyle.Copy().
		Width(containerWidth).
		Height(textareaHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Align(lipgloss.Center, lipgloss.Center)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render("Coming Soon:"),
		" ",
		container.Render(message),
	)
}

// Message types
type RequestSendMsg struct {
	Method string
	URL    string
	Body   string
}

