package http

import (
	"fmt"
	"net/http"
	"net/url"
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
	for _, valid := range validMethods {
		if method == valid {
			return nil
		}
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
	httpReq, err := http.NewRequest(strings.ToUpper(req.Method), requestURL, nil)
	if err != nil {
		log.Error("failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
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
	defer resp.Body.Close()

	duration := time.Since(start)

	response := &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
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
