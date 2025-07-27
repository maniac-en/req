package messages

type SwitchTabMsg struct {
	TabIndex int
}

type EditCollectionMsg struct {
	Label string
	Value string
}

type EditEndpointMsg struct {
	Name   string
	Method string
	URL    string
	ID     string
}
