package views

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/crud"
	"github.com/maniac-en/req/internal/tui/components"
)

type AddCollectionView struct {
	layout             components.Layout
	form               components.Form
	collectionsManager *collections.CollectionsManager
	width              int
	height             int
	submitting         bool
}

func NewAddCollectionView(collectionsManager *collections.CollectionsManager) AddCollectionView {
	inputs := []components.TextInput{
		components.NewTextInput("Name", "Enter collection name"),
	}
	
	form := components.NewForm("Add Collection", inputs)
	form.SetSubmitText("Create")
	
	return AddCollectionView{
		layout:             components.NewLayout(),
		form:               form,
		collectionsManager: collectionsManager,
	}
}

func (v AddCollectionView) Init() tea.Cmd {
	return nil
}

func (v AddCollectionView) Update(msg tea.Msg) (AddCollectionView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		v.layout.SetSize(v.width, v.height)
		v.form.SetSize(v.width-4, v.height-8) // Account for layout padding
		
	case tea.KeyMsg:
		if v.submitting {
			// Don't handle keys while submitting
			return v, nil
		}
		
		switch msg.String() {
		case "enter":
			return v, func() tea.Msg { return v.submitForm() }
		case "esc":
			return v, func() tea.Msg { return BackToCollectionsMsg{} }
		}
		
	case CollectionCreatedMsg:
		// Collection was created successfully
		return v, func() tea.Msg { return BackToCollectionsMsg{} }
		
	case CollectionCreateErrorMsg:
		// Handle error - for now just stop submitting
		v.submitting = false
	}
	
	// Update form
	v.form, cmd = v.form.Update(msg)
	return v, cmd
}

func (v *AddCollectionView) submitForm() tea.Msg {
	v.submitting = true
	values := v.form.GetValues()
	
	if len(values) == 0 || values[0] == "" {
		return CollectionCreateErrorMsg{err: crud.ErrInvalidInput}
	}
	
	return v.createCollection(values[0])
}

func (v *AddCollectionView) createCollection(name string) tea.Msg {
	collection, err := v.collectionsManager.Create(context.Background(), name)
	if err != nil {
		return CollectionCreateErrorMsg{err: err}
	}
	return CollectionCreatedMsg{collection: collection}
}

func (v AddCollectionView) View() string {
	if v.submitting {
		return v.layout.FullView(
			"Add Collection",
			"Creating collection...",
			"Please wait",
		)
	}

	content := v.form.View()
	instructions := "tab/↑↓: navigate • enter: create • esc: cancel"
	
	return v.layout.FullView(
		"Add Collection",
		content,
		instructions,
	)
}

// Messages for collection operations
type CollectionCreatedMsg struct {
	collection collections.CollectionEntity
}

type CollectionCreateErrorMsg struct {
	err error
}

type CollectionUpdatedMsg struct {
	collection collections.CollectionEntity
}

type CollectionUpdateErrorMsg struct {
	err error
}

type CollectionDeletedMsg struct {
	id int64
}

type CollectionDeleteErrorMsg struct {
	err error
}

type BackToCollectionsMsg struct{}

type EditCollectionMsg struct {
	collection collections.CollectionEntity
}