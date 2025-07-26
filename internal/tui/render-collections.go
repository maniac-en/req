package tui

import (
	// "log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func renderCollections(contentHeight, contentWidth int, m Model) string {
	collections := getCollections()
	options := []huh.Option[string]{}

	for _, coll := range collections {
		options = append(options, huh.NewOption(coll.Name, coll.ID))
	}

	// m.Tabs.Collections.Form = createCollectionForm(options)
	// log.Printf("options %v", m.Tabs.Collections.Form.View())

	selectRender := lipgloss.NewStyle().
		Width(contentWidth - 100).
		Height(contentHeight - 5).
		AlignVertical(lipgloss.Bottom).
		BorderStyle(lipgloss.RoundedBorder()).
		Render(m.Tabs.Collections.Form.View())

	innerBox := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		AlignHorizontal(lipgloss.Center).
		Render(selectRender)

	return innerBox
}

type Collection struct {
	Name string
	ID   string
}

// Dummy func that returns a bunch of collections
func getCollections() []Collection {
	return []Collection{
		{
			Name: "Coll1",
			ID:   "Some_ID",
		},
		{
			Name: "Coll2",
			ID:   "Some_ID",
		},
		{
			Name: "Coll3",
			ID:   "Some_ID",
		},
	}
}
