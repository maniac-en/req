package views

import (
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
)

func defaultListConfig[T, U any]() *optionsProvider.ListConfig[T, U] {
	config := optionsProvider.ListConfig[T, U]{
		ShowPagination:   false,
		ShowStatusBar:    false,
		ShowHelp:         false,
		ShowTitle:        false,
		FilteringEnabled: true,
	}
	return &config
}
