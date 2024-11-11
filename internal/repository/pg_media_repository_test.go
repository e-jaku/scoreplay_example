package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name:     "Empty braces",
			input:    []byte("{}"),
			expected: []string{},
		},
		{
			name:     "Single element",
			input:    []byte("{single}"),
			expected: []string{"single"},
		},
		{
			name:     "Multiple elements",
			input:    []byte("{one,two,three}"),
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "Elements with spaces",
			input:    []byte("{ one , two , three }"),
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "Elements with quotes",
			input:    []byte(`{"one","two","three"}`),
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "Elements with commas within quotes",
			input:    []byte(`{"one","two, and half","three"}`),
			expected: []string{"one", "two, and half", "three"},
		},
		{
			name:     "Elements with quotes and extra spaces",
			input:    []byte(`{" one "," two " , " three "}`),
			expected: []string{" one ", " two ", " three "},
		},
		{
			name:     "Mixed quotes and unquoted elements",
			input:    []byte(`{one,"two, and half",three}`),
			expected: []string{"one", "two, and half", "three"},
		},
		{
			name:     "Empty string inside braces",
			input:    []byte("{\"\"}"),
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractArray(tt.input)
			require.Equal(t, tt.expected, result, "extractArray(%s) returned incorrect result", tt.input)
		})
	}
}
