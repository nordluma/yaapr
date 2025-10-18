package allanime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeHex(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "797a7b7c7d7e7f",
			expected: "ABCDEFG",
		},
		{
			input:    "595a5b5c5d5e5f",
			expected: "abcdefg",
		},
	}

	for _, c := range testCases {
		actual, err := decodeHex(c.input)

		require.NoError(t, err)
		require.NotEmpty(t, actual)

		assert.Equal(
			t,
			c.expected,
			actual,
			"expected: %s - got: %s",
			c.expected,
			actual,
		)
	}
}
