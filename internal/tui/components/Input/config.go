package input

import "github.com/charmbracelet/bubbles/key"

type InputConfig struct {
	Prompt      string
	Placeholder string
	CharLimit   int
	Width       int
	KeyMap      InputKeyMaps
}

type InputKeyMaps struct {
	Accept key.Binding
	Back   key.Binding
}
