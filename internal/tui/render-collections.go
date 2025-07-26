package tui

import (
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderCollections(contentHeight, contentWidth int) string {
	collections := getCollections()
	options := []huh.Option[string]{}

	for _, coll := range collections {
		options = append(options, huh.NewOption(coll.Name, coll.ID))
	}

	collectionState, ok := m.Tabs.Collections.State.(*CollectionState)
	if !ok {
		log.Printf("State type: %T\n", m.Tabs.Collections.State)
		log.Printf("Error: State is not of type CollectionState")
		return ""
	}

	newForm := huh.NewSelect[string]().
		Options(options...)

	collectionState.Form = newForm

	selectRender := lipgloss.NewStyle().
		Width(contentWidth - 100).
		Height(contentHeight - 5).
		AlignVertical(lipgloss.Bottom).
		BorderStyle(lipgloss.RoundedBorder()).
		Render(collectionState.Form.View())

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
