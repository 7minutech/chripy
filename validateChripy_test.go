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

func TestCleanProfanity(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "I had something interesting for breakfast",
			expected: "I had something interesting for breakfast",
		},
		{
			input:    "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			input:    "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
		{
			input:    "I went to the Sharbert and got some kerfuffle. But it was not at all like fornax!",
			expected: "I went to the **** and got some kerfuffle. But it was not at all like fornax!",
		},
	}

	for _, c := range cases {
		actual := cleanProfanity(c.input)
		if actual != c.expected {
			t.Errorf("cleanProfanity(%s) == %s, expected: %s", c.input, actual, c.expected)
		}
	}
}
