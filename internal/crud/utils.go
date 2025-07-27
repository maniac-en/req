package crud

import (
	"fmt"
	"strings"
	"time"

	"github.com/maniac-en/req/internal/log"
)

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		log.Debug("validation failed: empty name")
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 100 {
		log.Debug("validation failed: name too long", "length", len(name))
		return fmt.Errorf("name cannot exceed 100 characters")
	}
	return nil
}

func ValidateID(id int64) error {
	if id <= 0 {
		log.Debug("validation failed: invalid ID", "id", id)
		return fmt.Errorf("ID must be positive")
	}
	return nil
}

// ParseTimestamp safely parses RFC3339 timestamp strings from database
func ParseTimestamp(timestamp string) time.Time {
	parsed, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Debug("failed to parse timestamp", "timestamp", timestamp, "error", err)
		return time.Time{}
	}
	return parsed
}
