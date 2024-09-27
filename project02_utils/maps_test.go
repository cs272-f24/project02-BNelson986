package project02utils

import (
	"testing"
)

func TestStop(t *testing.T) {
	m := &Maps{}

	m.stopWords = map[string]struct{}{
		"the": {},
		"and": {},
		"to":  {},
	}

	tests := []struct {
		name     string
		words    []string
		expected []string
	}{
		{
			name:     "No stop words",
			words:    []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "All stop words",
			words:    []string{"the", "and", "to"},
			expected: []string{},
		},
		{
			name:     "Mixed words",
			words:    []string{"hello", "the", "world", "and"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Empty input",
			words:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.removeStopWords(tt.words)
			// Same as 'resultCompare' function in sort_test.go
			if !func(a, b []string) bool {
				if len(a) != len(b) {
					return false
				}
				for i := range a {
					if a[i] != b[i] {
						return false
					}
				}
				return true
			}(result, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
