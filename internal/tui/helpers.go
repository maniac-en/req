package tui

import "github.com/charmbracelet/huh"

func createCollectionForm(options []huh.Option[string]) *huh.Form {
	var selected string

	formGroup := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&selected).
				Options(
					huh.NewOption("hel", "world"),
				),
		),
	)
	return formGroup
}
