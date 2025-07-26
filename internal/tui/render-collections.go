package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func renderCollections(contentHeight, contentWidth int, m Model) string {
	collections := getCollections()
	options := []huh.Option[string]{}

	for _, coll := range collections {
		options = append(options, huh.NewOption(coll.Name, coll.ID))
	}

	m.Tabs.Collections.Form = createCollectionForm(options)

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
			ID:   "Some_ID1",
		},
		{
			Name: "Coll2",
			ID:   "Some_ID2",
		},
		{
			Name: "Coll3",
			ID:   "Some_ID3",
		},
	}
}
