package main

import "testing"

func TestInProfainWords(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "clean",
			expected: false,
		},
		{
			input:    "kerfuffle",
			expected: true,
		},
		{
			input:    "Sharbert",
			expected: true,
		},
		{
			input:    "fornax!",
			expected: false,
		},
		{
			input:    "FORNAX",
			expected: true,
		},
	}

	for _, c := range cases {
		actual := inProfainWords(c.input)
		if actual != c.expected {
			t.Errorf("inProfainWords(%s) == %t, expected: %t", c.input, actual, c.expected)
		}
	}
}
