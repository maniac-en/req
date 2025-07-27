package views

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/crud"
	"github.com/maniac-en/req/internal/tui/components"
)

type CollectionsView struct {
	layout         components.Layout
	list           components.PaginatedList
	collectionsManager *collections.CollectionsManager
	width          int
	height         int
	initialized    bool
	
	// Backend pagination state
	currentPage    int
	pageSize       int
	pagination     crud.PaginationMetadata
}

func NewCollectionsView(collectionsManager *collections.CollectionsManager) CollectionsView {
	return CollectionsView{
		layout:            components.NewLayout(),
		collectionsManager: collectionsManager,
	}
}

func (v CollectionsView) Init() tea.Cmd {
	return v.loadCollections
}

func (v *CollectionsView) loadCollections() tea.Msg {
	return v.loadCollectionsPage(1, 20) // Load first page with 20 items
}

func (v *CollectionsView) loadCollectionsPage(page, pageSize int) tea.Msg {
	offset := (page - 1) * pageSize
	result, err := v.collectionsManager.ListPaginated(context.Background(), pageSize, offset)
	if err != nil {
		return collectionsLoadError{err: err}
	}
	return collectionsLoaded{
		collections: result.Collections,
		pagination:  result.PaginationMetadata,
		currentPage: page,
		pageSize:    pageSize,
	}
}

type collectionsLoaded struct {
	collections []collections.CollectionEntity
	pagination  crud.PaginationMetadata
	currentPage int
	pageSize    int
}

type collectionsLoadError struct {
	err error
}

func (v CollectionsView) Update(msg tea.Msg) (CollectionsView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		v.layout.SetSize(v.width, v.height)
		
		// Only set list size if it's been initialized
		if v.initialized {
			contentHeight := v.height - 4 // Account for header/footer
			v.list.SetSize(v.width-4, contentHeight)
		}
		
	case collectionsLoaded:
		items := make([]components.ListItem, len(msg.collections))
		for i, collection := range msg.collections {
			items[i] = components.NewCollectionItem(collection)
		}
		
		// Update pagination state
		v.currentPage = msg.currentPage
		v.pageSize = msg.pageSize
		v.pagination = msg.pagination
		
		// Create list with pagination info in title
		title := fmt.Sprintf("Collections (Page %d/%d)", v.currentPage, v.pagination.TotalPages)
		v.list = components.NewPaginatedList(items, title)
		
		if v.width > 0 && v.height > 0 {
			contentHeight := v.height - 4
			v.list.SetSize(v.width-4, contentHeight)
		}
		v.initialized = true
		
	case collectionsLoadError:
		// Handle error - for now just mark as initialized
		v.initialized = true
		
	case tea.KeyMsg:
		if !v.initialized {
			break
		}
		
		// Handle pagination keys first, before the list can consume them
		switch msg.String() {
		case "n", "right":
			// Next page
			if v.currentPage < v.pagination.TotalPages {
				return v, func() tea.Msg {
					return v.loadCollectionsPage(v.currentPage+1, v.pageSize)
				}
			}
			return v, nil
		case "p", "left":
			// Previous page
			if v.currentPage > 1 {
				return v, func() tea.Msg {
					return v.loadCollectionsPage(v.currentPage-1, v.pageSize)
				}
			}
			return v, nil
		default:
			// Forward other keys to the list
			v.list, cmd = v.list.Update(msg)
		}
		
	default:
		if v.initialized {
			v.list, cmd = v.list.Update(msg)
		}
	}
	
	return v, cmd
}

func (v CollectionsView) View() string {
	if !v.initialized {
		return v.layout.FullView(
			"Collections",
			"Loading collections...",
			"Please wait",
		)
	}

	content := v.list.View()
	
	// Build instructions with pagination info
	instructions := "↑↓: navigate • a: add • enter: edit • d: delete • q: quit"
	if v.pagination.TotalPages > 1 {
		instructions += fmt.Sprintf(" • p/n: prev/next page (%d/%d)", v.currentPage, v.pagination.TotalPages)
	}
	
	return v.layout.FullView(
		"Collections",
		content,
		instructions,
	)
}

