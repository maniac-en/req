package demo

import (
	"context"

	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/log"
)

type DemoGenerator struct {
	collectionsManager *collections.CollectionsManager
	endpointsManager   *endpoints.EndpointsManager
}

func NewDemoGenerator(collectionsManager *collections.CollectionsManager, endpointsManager *endpoints.EndpointsManager) *DemoGenerator {
	return &DemoGenerator{
		collectionsManager: collectionsManager,
		endpointsManager:   endpointsManager,
	}
}

func (d *DemoGenerator) PopulateDummyData(ctx context.Context) (bool, error) {
	log.Info("populating dummy data for demo")

	// Check if we already have collections
	result, err := d.collectionsManager.ListPaginated(ctx, 1, 0)
	if err != nil {
		log.Error("failed to check existing collections", "error", err)
		return false, err
	}

	if len(result.Collections) > 0 {
		log.Debug("dummy data already exists, skipping population", "collections_count", len(result.Collections))
		return false, nil
	}

	// Create demo collections and endpoints
	if err := d.createJSONPlaceholderCollection(ctx); err != nil {
		return false, err
	}

	if err := d.createReqresCollection(ctx); err != nil {
		return false, err
	}

	if err := d.createHTTPBinCollection(ctx); err != nil {
		return false, err
	}

	log.Info("dummy data populated successfully")
	return true, nil
}

func (d *DemoGenerator) createJSONPlaceholderCollection(ctx context.Context) error {
	collection, err := d.collectionsManager.Create(ctx, "JSONPlaceholder API")
	if err != nil {
		log.Error("failed to create JSONPlaceholder collection", "error", err)
		return err
	}

	endpoints := []endpoints.EndpointData{
		{
			CollectionID: collection.ID,
			Name:         "Get All Posts",
			Method:       "GET",
			URL:          "https://jsonplaceholder.typicode.com/posts",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  "",
		},
		{
			CollectionID: collection.ID,
			Name:         "Get Single Post",
			Method:       "GET",
			URL:          "https://jsonplaceholder.typicode.com/posts/1",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  "",
		},
		{
			CollectionID: collection.ID,
			Name:         "Create Post",
			Method:       "POST",
			URL:          "https://jsonplaceholder.typicode.com/posts",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  `{"title": "My New Post", "body": "This is the content of my new post", "userId": 1}`,
		},
		{
			CollectionID: collection.ID,
			Name:         "Update Post",
			Method:       "PUT",
			URL:          "https://jsonplaceholder.typicode.com/posts/1",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  `{"id": 1, "title": "Updated Post", "body": "This post has been updated", "userId": 1}`,
		},
		{
			CollectionID: collection.ID,
			Name:         "Delete Post",
			Method:       "DELETE",
			URL:          "https://jsonplaceholder.typicode.com/posts/1",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  "",
		},
	}

	return d.createEndpoints(ctx, endpoints)
}

func (d *DemoGenerator) createReqresCollection(ctx context.Context) error {
	collection, err := d.collectionsManager.Create(ctx, "ReqRes API")
	if err != nil {
		log.Error("failed to create ReqRes collection", "error", err)
		return err
	}

	endpoints := []endpoints.EndpointData{
		{
			CollectionID: collection.ID,
			Name:         "List Users",
			Method:       "GET",
			URL:          "https://reqres.in/api/users",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{"page": "2"},
			RequestBody:  "",
		},
		{
			CollectionID: collection.ID,
			Name:         "Single User",
			Method:       "GET",
			URL:          "https://reqres.in/api/users/2",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  "",
		},
		{
			CollectionID: collection.ID,
			Name:         "Create User",
			Method:       "POST",
			URL:          "https://reqres.in/api/users",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  `{"name": "morpheus", "job": "leader"}`,
		},
		{
			CollectionID: collection.ID,
			Name:         "Login",
			Method:       "POST",
			URL:          "https://reqres.in/api/login",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  `{"email": "eve.holt@reqres.in", "password": "cityslicka"}`,
		},
	}

	return d.createEndpoints(ctx, endpoints)
}

func (d *DemoGenerator) createHTTPBinCollection(ctx context.Context) error {
	collection, err := d.collectionsManager.Create(ctx, "HTTPBin Testing")
	if err != nil {
		log.Error("failed to create HTTPBin collection", "error", err)
		return err
	}

	endpoints := []endpoints.EndpointData{
		{
			CollectionID: collection.ID,
			Name:         "Test GET",
			Method:       "GET",
			URL:          "https://httpbin.org/get",
			Headers:      `{"User-Agent": "Req-Terminal-Client/1.0"}`,
			QueryParams:  map[string]string{"test": "value", "demo": "true"},
			RequestBody:  "",
		},
		{
			CollectionID: collection.ID,
			Name:         "Test POST JSON",
			Method:       "POST",
			URL:          "https://httpbin.org/post",
			Headers:      `{"Content-Type": "application/json", "User-Agent": "Req-Terminal-Client/1.0"}`,
			QueryParams:  map[string]string{},
			RequestBody:  `{"message": "Hello from Req!", "timestamp": "2024-01-15T10:30:00Z", "data": {"key": "value"}}`,
		},
		{
			CollectionID: collection.ID,
			Name:         "Test Headers",
			Method:       "GET",
			URL:          "https://httpbin.org/headers",
			Headers:      `{"Authorization": "Bearer demo-token", "X-Custom-Header": "req-demo"}`,
			QueryParams:  map[string]string{},
			RequestBody:  "",
		},
		{
			CollectionID: collection.ID,
			Name:         "Test Status Codes",
			Method:       "GET",
			URL:          "https://httpbin.org/status/200",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{},
			RequestBody:  "",
		},
	}

	return d.createEndpoints(ctx, endpoints)
}

func (d *DemoGenerator) createEndpoints(ctx context.Context, endpointData []endpoints.EndpointData) error {
	for _, data := range endpointData {
		_, err := d.endpointsManager.CreateEndpoint(ctx, data)
		if err != nil {
			log.Error("failed to create endpoint", "name", data.Name, "error", err)
			return err
		}
		log.Debug("created demo endpoint", "name", data.Name, "method", data.Method, "url", data.URL)
	}
	return nil
}