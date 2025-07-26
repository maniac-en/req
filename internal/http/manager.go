package http

import (
	"fmt"
	"net/http"
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
