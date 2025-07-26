package http

import (
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
