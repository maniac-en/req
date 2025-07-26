package messages

type SwitchTabMsg struct {
	TabIndex int
}

type EditCollectionMsg struct {
	Label string
	Value string
}
