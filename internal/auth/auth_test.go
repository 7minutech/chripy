package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {

	password1 := "password123"
	hash1, err := HashPassword(password1)
	if err != nil {
		t.Errorf("HashPassword(%s) err was not nil", password1)
	}

	password2 := "seasame456"
	hash2False, err := HashPassword(password1)
	if err != nil {
		t.Errorf("HashPassword(%s) err was not nil", password1)
	}

	hash2True, err := HashPassword(password2)
	if err != nil {
		t.Errorf("HashPassword(%s) err was not nil", password2)
	}

	cases := []struct {
		input struct {
			password string
			hash     string
		}
		expected bool
	}{
		{
			input: struct {
				password string
				hash     string
			}{password: password1, hash: hash1},
			expected: true,
		},
		{
			input: struct {
				password string
				hash     string
			}{password: password2, hash: hash2False},
			expected: false,
		},
		{
			input: struct {
				password string
				hash     string
			}{password: password2, hash: hash2True},
			expected: true,
		},
	}

	for _, c := range cases {

		ok, err := CheckPasswordHash(c.input.password, c.input.hash)
		if err != nil {
			t.Errorf("CheckPasswordHash(%s, %s) err was not nil", c.input.password, c.input.hash)
		}

		if ok != c.expected {
			t.Errorf("CheckPasswordHash(%s, %s) == %t", c.input.password, c.input.hash, ok)
		}
	}
}
