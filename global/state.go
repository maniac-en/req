package global

type State struct {
	currentCollection string
}

func NewGlobalState() *State {
	return &State{
		currentCollection: "",
	}
}

// Gets the current collection from the app state
func (s *State) GetCurrentCollection() string {
	return s.currentCollection
}

// Sets the current collection to the app state
func (s *State) SetCurrentCollection(collection string) {
	s.currentCollection = collection
}
