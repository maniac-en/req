package app

import (
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/backend/history"
	"github.com/maniac-en/req/internal/backend/http"
)

type Context struct {
	Collections      *collections.CollectionsManager
	Endpoints        *endpoints.EndpointsManager
	HTTP             *http.HTTPManager
	History          *history.HistoryManager
	DummyDataCreated bool
	Version          string
}

func NewContext(
	collections *collections.CollectionsManager,
	endpoints *endpoints.EndpointsManager,
	httpManager *http.HTTPManager,
	history *history.HistoryManager,
	version string,
) *Context {
	return &Context{
		Collections:      collections,
		Endpoints:        endpoints,
		HTTP:             httpManager,
		History:          history,
		DummyDataCreated: false,
		Version:          version,
	}
}

func (c *Context) SetDummyDataCreated(created bool) {
	c.DummyDataCreated = created
}
