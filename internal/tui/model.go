package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

type Tab struct {
	Title   string
	Content func(contentHeight, contentWidth int) string
	State   interface{}
}

type Tabs struct {
	Home         Tab
	Collections  Tab
	Endpoints    Tab
	Environments Tab
}

type Model struct {
	Tabs      Tabs
	ActiveTab *Tab
	Width     int
	Height    int
	Keybinds  Input
}

func InitModel() (Model, error) {
	model := Model{
		Keybinds: initKeybinds(),
	}

	model.Tabs = model.InitTabs()
	model.ActiveTab = &model.Tabs.Home

	fmt.Printf("State type in model.go: %T\n", model.Tabs.Collections.State)

	return model, nil
}

func (m Model) InitTabs() Tabs {
	return Tabs{
		Home: Tab{
			Title:   "Home",
			Content: m.renderHome,
		},
		Collections: Tab{
			Title:   "Collections",
			Content: m.renderCollections,
			State:   m.createCollectionsState(),
		},
		Endpoints: Tab{
			Title:   "Endpoints",
			Content: m.renderHome,
		},
		Environments: Tab{
			Title:   "Environments",
			Content: m.renderEnvironments,
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

// state stuff
type CollectionState struct {
	Form     *huh.Select[string]
	Selected string
}
