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

func TestEncodeDecode(t *testing.T) {
	samples := []string{
		"HelloWorld",
		"ABCDEFG",
		"abcdefg",
		"0123456789",
		"!@#_-~",
	}

	for _, s := range samples {
		encoded := encode(s)
		decoded, err := decodeHex(encoded)
		require.NoError(t, err)

		assert.Equal(
			t, s, decoded,
			"Encode/decode mismatch: got %q, want %q",
			decoded, s,
		)
	}
}

func TestInvalidHex(t *testing.T) {
	_, err := decodeHex("zzzz")
	require.Error(t, err, "expected error for invalid hex input, got nil")
}

func BenchmarkDecodeXORFunc(b *testing.B) {
	v := "797a7b7c7d7e7f"
	for b.Loop() {
		_, _ = decodeHex(v)
	}
}
