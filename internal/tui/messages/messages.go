package messages

type ItemAdded struct {
	Item string
}
type ItemEdited struct {
	Item   string
	ItemID int64
}

type DeleteItem struct {
	ItemID int64
}

type ChooseItem[T any] struct {
	Item   T
	Source string
}

type DeactivateView struct{}

type NavigateToView struct {
	ViewName string
	Data     interface{}
}
