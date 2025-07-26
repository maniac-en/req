package tui

import (
	"errors"
)

type Tab struct {
	Title   string
	Content func(contentHeight, contentWidth int) string
}

type Model struct {
	Tabs      []Tab
	ActiveTab int
	Width     int
	Height    int
	Keybinds  Input
}

func InitModel(tabs []Tab) (Model, error) {
	if len(tabs) == 0 {
		return Model{}, errors.New("Tabs array cannot be empty")
	}
	return Model{
		Tabs:      tabs,
		ActiveTab: 0,
		Keybinds:  initKeybinds(),
	}, nil
}

func InitTabs() []Tab {
	return []Tab{
		{
			Title:   "Home",
			Content: renderHome,
		},
		{
			Title:   "Collections",
			Content: renderCollections,
		},
		{
			Title:   "Endpoints",
			Content: renderHome,
		},
		{
			Title:   "Environments",
			Content: renderEnvironments,
		},
	}
}

// keybinds
type Input struct {
	Quit              string
	KeyboardInterrupt string
	Collections       string
	Endpoints         string
	Environments      string
}

func initKeybinds() Input {
	return Input{
		Quit:              "q",
		KeyboardInterrupt: "ctrl+c",
		Collections:       "c",
		Endpoints:         "e",
		Environments:      "n",
	}
}
