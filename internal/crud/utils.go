package crud

import (
	"fmt"
	"strings"
)

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("name cannot exceed 100 characters")
	}
	return nil
}

func ValidateID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID must be positive")
	}
	return nil
}