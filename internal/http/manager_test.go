package http

import (
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPManager(t *testing.T) {
	manager := NewHTTPManager()
	if manager == nil {
		t.Fatal("NewHTTPManager returned nil")
	}
	if manager.Client == nil {
		t.Fatal("HTTPManager client is nil")
	}
	if manager.Client.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", manager.Client.Timeout)
	}
}

func TestValidateMethod(t *testing.T) {
	tests := []struct {
		method string
		valid  bool
	}{
		{"GET", true},
		{"POST", true},
		{"put", true},
		{"delete", true},
		{"INVALID", false},
		{"", false},
	}

	for _, test := range tests {
		err := validateMethod(test.method)
		if test.valid && err != nil {
			t.Errorf("expected %s to be valid, got error: %v", test.method, err)
		}
		if !test.valid && err == nil {
			t.Errorf("expected %s to be invalid, got no error", test.method)
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"https://example.com", true},
		{"http://localhost:8080", true},
		{"ftp://invalid.com", false},
		{"", false},
		{"not-a-url", false},
	}

	for _, test := range tests {
		err := validateURL(test.url)
		if test.valid && err != nil {
			t.Errorf("expected %s to be valid, got error: %v", test.url, err)
		}
		if !test.valid && err == nil {
			t.Errorf("expected %s to be invalid, got no error", test.url)
		}
	}
}

func TestValidateRequest(t *testing.T) {
	manager := NewHTTPManager()

	validReq := &Request{
		Method: "GET",
		URL:    "https://example.com",
	}

	if err := manager.ValidateRequest(validReq); err != nil {
		t.Errorf("expected valid request to pass validation, got: %v", err)
	}

	invalidReq := &Request{
		Method: "INVALID",
		URL:    "not-a-url",
	}

	if err := manager.ValidateRequest(invalidReq); err == nil {
		t.Error("expected invalid request to fail validation")
	}
}

func TestBuildURL(t *testing.T) {
	manager := NewHTTPManager()

	tests := []struct {
		baseURL     string
		queryParams map[string]string
		expected    string
	}{
		{"https://example.com", nil, "https://example.com"},
		{"https://example.com", map[string]string{}, "https://example.com"},
		{"https://example.com", map[string]string{"foo": "bar"}, "https://example.com?foo=bar"},
	}

	for _, test := range tests {
		result, err := manager.buildURL(test.baseURL, test.queryParams)
		if err != nil {
			t.Errorf("buildURL failed: %v", err)
		}
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}

func TestSetHeaders(t *testing.T) {
	manager := NewHTTPManager()
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "req-cli",
	}

	err := manager.setHeaders(req, headers)
	if err != nil {
		t.Errorf("setHeaders failed: %v", err)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type header not set correctly")
	}
}

func TestSetContentType(t *testing.T) {
	manager := NewHTTPManager()

	tests := []struct {
		body     string
		expected string
	}{
		{`{"key": "value"}`, "application/json"},
		{`[1, 2, 3]`, "application/json"},
		{"plain text", "text/plain"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("POST", "https://example.com", nil)
		manager.setContentType(req, test.body)

		if req.Header.Get("Content-Type") != test.expected {
			t.Errorf("for body %q, expected %q, got %q",
				test.body, test.expected, req.Header.Get("Content-Type"))
		}
	}
}
