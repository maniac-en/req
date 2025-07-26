package tui

import (
	"github.com/charmbracelet/huh"
)

type Tab struct {
	Title   string
	Content func(contentHeight, contentWidth int, model Model) string
	Form    *huh.Select[string]
}

type Tabs struct {
	Home         Tab
	Collections  Tab
	Endpoints    Tab
	Environments Tab
}

type Model struct {
	Tabs        Tabs
	ActiveTab   *Tab
	Width       int
	Height      int
	Keybinds    Input
	GlobalState States
}

type States struct {
	Collection *CollectionState
}

func InitModel() (Model, error) {
	model := Model{
		Keybinds: initKeybinds(),
	}

	model.Tabs = InitTabs()
	model.ActiveTab = &model.Tabs.Home

	return model, nil
}

func InitTabs() Tabs {
	return Tabs{
		Home: Tab{
			Title:   "Home",
			Content: renderHome,
		},
		Collections: Tab{
			Title:   "Collections",
			Content: renderCollections,
			Form:    createCollectionForm([]huh.Option[string]{}),
		},
		Endpoints: Tab{
			Title:   "Endpoints",
			Content: renderHome,
		},
		Environments: Tab{
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

// state stuff
type CollectionState struct {
	Selected string
	Error    error
}
