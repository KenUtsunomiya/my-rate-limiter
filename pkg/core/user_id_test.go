package core

import (
	"testing"
)

func TestUserID_IsValid_Valid(t *testing.T) {
	testCases := []string{
		"0", "1", "1234567890123456789",
	}

	for _, tc := range testCases {
		id := UserID(tc)
		if !id.IsValid() {
			t.Errorf("expected valid for input %s, got invalid", tc)
		}
	}
}

func TestUserID_IsValid_Invalid(t *testing.T) {
	testCases := []string{
		"", "a", "123abc", "-123", "123 ", " 123",
		"12345678901234567890",
	}

	for _, tc := range testCases {
		id := UserID(tc)
		if id.IsValid() {
			t.Errorf("expected invalid for input %s, got valid", tc)
		}
	}
}
