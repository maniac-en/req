package crud

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"valid-name", true},
		{"", false},
		{"   ", false},
		{string(make([]byte, 101)), false},
		{"normal name", true},
	}

	for _, test := range tests {
		err := ValidateName(test.name)
		if test.valid && err != nil {
			t.Errorf("expected %q to be valid, got error: %v", test.name, err)
		}
		if !test.valid && err == nil {
			t.Errorf("expected %q to be invalid, got no error", test.name)
		}
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		id    int64
		valid bool
	}{
		{1, true},
		{100, true},
		{0, false},
		{-1, false},
	}

	for _, test := range tests {
		err := ValidateID(test.id)
		if test.valid && err != nil {
			t.Errorf("expected ID %d to be valid, got error: %v", test.id, err)
		}
		if !test.valid && err == nil {
			t.Errorf("expected ID %d to be invalid, got no error", test.id)
		}
	}
}
