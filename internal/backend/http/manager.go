package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/maniac-en/req/internal/log"
)

func NewHTTPManager() *HTTPManager {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &HTTPManager{
		Client: client,
	}
}

func validateMethod(method string) error {
	method = strings.ToUpper(strings.TrimSpace(method))
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	if slices.Contains(validMethods, method) {
		return nil
	}
	return fmt.Errorf("invalid HTTP method: %s", method)
}

func validateURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	return nil
}

func (h *HTTPManager) ValidateRequest(req *Request) error {
	if err := validateMethod(req.Method); err != nil {
		log.Error("invalid method", "method", req.Method, "error", err)
		return err
	}
	if err := validateURL(req.URL); err != nil {
		log.Error("invalid URL", "url", req.URL, "error", err)
		return err
	}
	return nil
}

func (h *HTTPManager) ExecuteRequest(req *Request) (*Response, error) {
	if err := h.ValidateRequest(req); err != nil {
		return nil, err
	}

	log.Debug("executing HTTP request", "method", req.Method, "url", req.URL)

	requestURL, err := h.buildURL(req.URL, req.QueryParams)
	if err != nil {
		log.Error("failed to build URL", "error", err)
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	start := time.Now()

	var body io.Reader
	if req.Body != "" && (strings.ToUpper(req.Method) == "POST" || strings.ToUpper(req.Method) == "PUT" || strings.ToUpper(req.Method) == "PATCH") {
		body = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(strings.ToUpper(req.Method), requestURL, body)
	if err != nil {
		log.Error("failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		h.setContentType(httpReq, req.Body)
	}

	if err := h.setHeaders(httpReq, req.Headers); err != nil {
		log.Error("failed to set headers", "error", err)
		return nil, fmt.Errorf("failed to set headers: %w", err)
	}

	resp, err := h.Client.Do(httpReq)
	if err != nil {
		log.Error("HTTP request failed", "error", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error("failed to close response body", "error", closeErr)
		}
	}()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	duration := time.Since(start)

	// Log warnings for concerning HTTP responses
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		log.Warn("HTTP client error", "status", resp.StatusCode, "url", req.URL)
	} else if resp.StatusCode >= 500 {
		log.Warn("HTTP server error", "status", resp.StatusCode, "url", req.URL)
	}

	// Warn about slow requests (>5 seconds)
	if duration > 5*time.Second {
		log.Warn("slow HTTP request", "duration", duration, "url", req.URL)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       string(responseBody),
		Duration:   duration,
	}

	log.Info("HTTP request completed", "status", resp.StatusCode, "duration", duration)
	return response, nil
}

func (h *HTTPManager) buildURL(baseURL string, queryParams map[string]string) (string, error) {
	if len(queryParams) == 0 {
		return baseURL, nil
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	values := parsedURL.Query()
	for key, value := range queryParams {
		values.Set(key, value)
	}
	parsedURL.RawQuery = values.Encode()

	return parsedURL.String(), nil
}

func (h *HTTPManager) setHeaders(req *http.Request, headers map[string]string) error {
	for key, value := range headers {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("header key cannot be empty")
		}
		req.Header.Set(key, value)
	}
	return nil
}

func (h *HTTPManager) setContentType(req *http.Request, body string) {
	if req.Header.Get("Content-Type") != "" {
		return
	}

	body = strings.TrimSpace(body)
	if strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[") {
		if json.Valid([]byte(body)) {
			req.Header.Set("Content-Type", "application/json")
			return
		}
	}

	req.Header.Set("Content-Type", "text/plain")
}
