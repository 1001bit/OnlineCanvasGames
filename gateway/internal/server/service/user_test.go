package service

import "testing"

func TestMapToUser(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]any
		expected *User
	}{
		{
			name: "valid",
			input: map[string]any{
				"name": "bob",
				"date": "once upon a time",
			},
			expected: &User{Name: "bob", Date: "once upon a time"},
		},
		{
			name: "dateless",
			input: map[string]any{
				"id":   2.0,
				"name": "jack",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		if mapToUser(tc.input) == nil && tc.expected == nil {
			return
		}

		if *tc.expected != *mapToUser(tc.input) {
			t.Error("Error in", tc.name, mapToUser(tc.input), tc.expected)
		}
	}
}
