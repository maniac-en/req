package views

import (
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
)

func defaultListConfig[T, C any]() *optionsProvider.ListConfig[T, C] {
	config := optionsProvider.ListConfig[T, C]{
		ShowPagination:   false,
		ShowStatusBar:    false,
		ShowHelp:         false,
		ShowTitle:        false,
		FilteringEnabled: true,
	}
	return &config
}
