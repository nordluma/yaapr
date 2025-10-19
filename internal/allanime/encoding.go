package allanime

import "encoding/hex"

// Replacement table used for decoding and encoding
//
// Uppercase letters
// 79 => A, 7a => B, 7b => C, 7c => D, 7d => E, 7e => F, 7f => G, 70 => H,
// 71 => I, 72 => J, 73 => K, 74 => L, 75 => M, 76 => N, 77 => O, 68 => P,
// 69 => Q, 6a => R, 6b => S, 6c => T, 6d => U, 6e => V, 6f => W, 60 => X,
// 61 => Y, 62 => Z,
//
// Lowercase letters
// 59 => a, 5a => b, 5b => c, 5c => d, 5d => e, 5e => f, 5f => g, 50 => h,
// 51 => i, 52 => j, 53 => k, 54 => l, 55 => m, 56 => n, 57 => o, 48 => p,
// 49 => q, 4a => r, 4b => s, 4c => t, 4d => u, 4e => v, 4f => w, 40 => x,
// 41 => y, 42 => z,
//
// Numbers
// 08 => 0, 09 => 1, 0a => 2, 0b => 3, 0c => 4, 0d => 5, 0e => 6, 0f => 7,
// 00 => 8, 01 => 9,
//
// Special characters
// 15 => -, 16 => ., 67 => _, 46 => ~, 02 => :, 17 => /, 07 => ?, 1b => #,
// 63 => [, 65 => ], 78 => @, 19 => !, 1c => $, 1e => &, 10 => (, 11 => ),
// 12 => *, 13 => +, 14 => ,, 03 => ;, 05 => =, 1d => %,

const key = 0x38

func decodeHex(encoded string) (string, error) {
	data, err := hex.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	for i := range data {
		data[i] ^= key
	}

	return string(data), nil
}

var _ = encode // this is to silence gopls complaining about a unused function

func encode(s string) string {
	out := make([]byte, len(s))
	for i, b := range []byte(s) {
		out[i] = b ^ key
	}

	return hex.EncodeToString(out)
}
