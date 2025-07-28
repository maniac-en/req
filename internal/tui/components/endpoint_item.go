package components

import (
	"fmt"

	"github.com/maniac-en/req/internal/backend/endpoints"
)

type EndpointItem struct {
	endpoint endpoints.EndpointEntity
}

func NewEndpointItem(endpoint endpoints.EndpointEntity) EndpointItem {
	return EndpointItem{
		endpoint: endpoint,
	}
}

func (i EndpointItem) FilterValue() string {
	return i.endpoint.Name
}

func (i EndpointItem) GetID() string {
	return fmt.Sprintf("%d", i.endpoint.ID)
}

func (i EndpointItem) GetTitle() string {
	return fmt.Sprintf("%s %s", i.endpoint.Method, i.endpoint.Name)
}

func (i EndpointItem) GetDescription() string {
	return i.endpoint.Url
}

func (i EndpointItem) Title() string {
	return i.GetTitle()
}

func (i EndpointItem) Description() string {
	return i.GetDescription()
}

