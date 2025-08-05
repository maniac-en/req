package views

import (
	optionsProvider "github.com/maniac-en/req/internal/tui/components/OptionsProvider"
)

func defaultListConfig[T any]() *optionsProvider.ListConfig[T] {
	config := optionsProvider.ListConfig[T]{
		ShowPagination:   false,
		ShowHelp:         false,
		ShowTitle:        false,
		FilteringEnabled: false,
	}
	return &config
}
