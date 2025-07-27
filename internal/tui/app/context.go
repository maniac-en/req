package app

import (
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/backend/history"
	"github.com/maniac-en/req/internal/backend/http"
)

type Context struct {
	Collections *collections.CollectionsManager
	Endpoints   *endpoints.EndpointsManager
	HTTP        *http.HTTPManager
	History     *history.HistoryManager
}

func NewContext(
	collections *collections.CollectionsManager,
	endpoints *endpoints.EndpointsManager,
	httpManager *http.HTTPManager,
	history *history.HistoryManager,
) *Context {
	return &Context{
		Collections: collections,
		Endpoints:   endpoints,
		HTTP:        httpManager,
		History:     history,
	}
}

