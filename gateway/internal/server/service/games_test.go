package service

import (
	"testing"
)

func TestMapToGame(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]any
		expected *Game
	}{
		{
			name: "valid",
			input: map[string]any{
				"id":    1.0,
				"title": "platformer",
			},
			expected: &Game{ID: 1, Title: "platformer"},
		},
		{
			name: "invalid",
			input: map[string]any{
				"id":    "a",
				"title": "test game",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		if mapToGame(tc.input) == nil && tc.expected == nil {
			return
		}

		if *tc.expected != *mapToGame(tc.input) {
			t.Error("Error in", tc.name, mapToGame(tc.input), tc.expected)
		}
	}
}
