package utils

import "testing"

func TestSemverGreater(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   bool
	}{
		{"3.78.1", "3.78.0", true},
		{"3.78.0", "3.78.1", false},
		{"3.78.1", "3.78.1", false},
		{"4.0.0", "3.99.9", true},
		{"3.100.0", "3.99.0", true},
		{"3.1", "3.0.1", true},
		{"3.0.1", "3.1", false},
		{"10.0.0", "9.9.9", true},
	}

	for _, tc := range tests {
		got := SemverGreater(tc.v1, tc.v2)
		if got != tc.want {
			t.Errorf("SemverGreater(%q, %q) = %v, want %v", tc.v1, tc.v2, got, tc.want)
		}
	}
}
