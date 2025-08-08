package views

import (
	"github.com/charmbracelet/bubbles/list"
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
	"github.com/maniac-en/req/internal/tui/styles"
)

func createDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = styles.SelectedListStyle
	d.Styles.SelectedDesc = styles.SelectedListStyle

	return d
}

func defaultListConfig[T, U any]() *optionsProvider.ListConfig[T, U] {
	config := optionsProvider.ListConfig[T, U]{
		ShowPagination:   false,
		ShowStatusBar:    false,
		ShowHelp:         false,
		ShowTitle:        false,
		FilteringEnabled: true,

		Delegate: createDelegate(),
	}
	return &config
}
