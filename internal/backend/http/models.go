// Package http provides HTTP client functionality for making HTTP requests.
package http

import (
	"net/http"
	"time"
)

type HTTPManager struct {
	Client *http.Client
}

type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        string
}

type Response struct {
	StatusCode int
	Status     string
	Headers    map[string][]string
	Body       string
	Duration   time.Duration
}
