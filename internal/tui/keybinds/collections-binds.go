package keybinds

import "github.com/charmbracelet/bubbles/key"

type ListKeyMap struct {
	CursorUp             key.Binding
	CursorDown           key.Binding
	NextPage             key.Binding
	PrevPage             key.Binding
	Filter               key.Binding
	ClearFilter          key.Binding
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding
	AddItem              key.Binding
	EditItem             key.Binding
	DeleteItem           key.Binding
	Accept               key.Binding
	Back                 key.Binding
}

func (c ListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{c.CursorUp, c.CursorDown, c.NextPage, c.PrevPage, c.AddItem, c.EditItem, c.DeleteItem}
}

func (c ListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{c.CursorUp, c.CursorDown, c.NextPage, c.PrevPage},
	}
}

func NewListKeyMap() *ListKeyMap {
	return &ListKeyMap{
		CursorUp:             Keys.Up,
		CursorDown:           Keys.Down,
		NextPage:             Keys.NextPage,
		PrevPage:             Keys.PrevPage,
		Filter:               Keys.Filter,
		ClearFilter:          Keys.ClearFilter,
		CancelWhileFiltering: Keys.CancelWhileFiltering,
		AcceptWhileFiltering: Keys.AcceptWhileFiltering,
		AddItem:              Keys.InsertItem,
		DeleteItem:           Keys.Remove,
		EditItem:             Keys.EditItem,
		Accept:               Keys.Choose,
		Back:                 Keys.Back,
	}
}
