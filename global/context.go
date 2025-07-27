package global

import (
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/backend/history"
	"github.com/maniac-en/req/internal/backend/http"
)

type AppContext struct {
	Collections *collections.CollectionsManager
	Endpoints   *endpoints.EndpointsManager
	HTTP        *http.HTTPManager
	History     *history.HistoryManager
}

var globalAppContext *AppContext

func SetAppContext(ctx *AppContext) {
	globalAppContext = ctx
}

func GetAppContext() *AppContext {
	return globalAppContext
}
