package tui

import "github.com/charmbracelet/huh"

func createCollectionForm(options []huh.Option[string]) *huh.Select[string] {
	var selected string

	return huh.NewSelect[string]().
		Value(&selected).
		Options(options...)
}
