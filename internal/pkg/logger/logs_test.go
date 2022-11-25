package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExtractMapFromKeyvals(t *testing.T) {
	testCases := []struct {
		name     string
		keyvals  []any
		expected map[string]any
	}{
		{
			name:    "one",
			keyvals: []any{"hello"},
			expected: map[string]any{
				"msg": "hello",
			},
		},
		{
			name:    "two",
			keyvals: []any{"msg", "hello"},
			expected: map[string]any{
				"msg": "hello",
			},
		},
		{
			name:    "three",
			keyvals: []any{"msg", "hello", "and"},
			expected: map[string]any{
				"msg": "hello",
				"and": nil,
			},
		},
		{
			name:    "error_on_key",
			keyvals: []any{"msg", "hello", fmt.Errorf("err: %v", "aaa"), "ok"},
			expected: map[string]any{
				"msg":      "hello",
				"err: aaa": "ok",
			},
		},
		{
			name:    "slice_on_key",
			keyvals: []any{"msg", "hello", make([]string, 0, 5), "ok"},
			expected: map[string]any{
				"msg": "hello",
				"[]":  "ok",
			},
		},
		{
			name: "error_with_tags",
			keyvals: []any{
				"msg", "failed to update user",
				"user_id", 5,
				"updated_at", time.Date(2020, 10, 10, 10, 10, 10, 0, time.UTC),
			},
			expected: map[string]any{
				"msg":        "failed to update user",
				"user_id":    5,
				"updated_at": time.Date(2020, 10, 10, 10, 10, 10, 0, time.UTC),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := ExtractMapFromKeyvals(testCase.keyvals...)
			require.Equal(t, testCase.expected, actual)
		})
	}
}
