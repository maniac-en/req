package views

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/crud"
	"github.com/maniac-en/req/internal/tui/components"
)

type EditCollectionView struct {
	layout             components.Layout
	form               components.Form
	collectionsManager *collections.CollectionsManager
	collection         collections.CollectionEntity
	width              int
	height             int
	submitting         bool
}

func NewEditCollectionView(collectionsManager *collections.CollectionsManager, collection collections.CollectionEntity) EditCollectionView {
	inputs := []components.TextInput{
		components.NewTextInput("Name", "Enter collection name"),
	}
	
	// Pre-populate with existing collection name
	inputs[0].SetValue(collection.Name)
	
	form := components.NewForm("Edit Collection", inputs)
	form.SetSubmitText("Update")
	
	return EditCollectionView{
		layout:             components.NewLayout(),
		form:               form,
		collectionsManager: collectionsManager,
		collection:         collection,
	}
}

func (v EditCollectionView) Init() tea.Cmd {
	return nil
}

func (v EditCollectionView) Update(msg tea.Msg) (EditCollectionView, tea.Cmd) {
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
		
	case CollectionUpdatedMsg:
		// Collection was updated successfully
		return v, func() tea.Msg { return BackToCollectionsMsg{} }
		
	case CollectionUpdateErrorMsg:
		// Handle error - for now just stop submitting
		v.submitting = false
	}
	
	// Update form
	v.form, cmd = v.form.Update(msg)
	return v, cmd
}

func (v *EditCollectionView) submitForm() tea.Msg {
	v.submitting = true
	values := v.form.GetValues()
	
	if len(values) == 0 || values[0] == "" {
		return CollectionUpdateErrorMsg{err: crud.ErrInvalidInput}
	}
	
	return v.updateCollection(values[0])
}

func (v *EditCollectionView) updateCollection(name string) tea.Msg {
	// Update the collection using the manager's Update method
	updatedCollection, err := v.collectionsManager.Update(context.Background(), v.collection.ID, name)
	if err != nil {
		return CollectionUpdateErrorMsg{err: err}
	}
	return CollectionUpdatedMsg{collection: updatedCollection}
}

func (v EditCollectionView) View() string {
	if v.submitting {
		return v.layout.FullView(
			"Edit Collection",
			"Updating collection...",
			"Please wait",
		)
	}

	content := v.form.View()
	instructions := "tab/↑↓: navigate • enter: update • esc: cancel"
	
	return v.layout.FullView(
		"Edit Collection",
		content,
		instructions,
	)
}