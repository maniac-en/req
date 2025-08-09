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
type DeactivateView struct{}
