package tui

import "github.com/charmbracelet/huh"

func (m Model) createCollectionsState() *CollectionState {
	var selected string
	selectInput := huh.NewSelect[string]().
		Value(&selected)

	return &CollectionState{
		Form:     selectInput,
		Selected: selected,
	}
}
