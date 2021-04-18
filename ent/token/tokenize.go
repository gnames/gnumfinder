package token

import (
	gner "github.com/gnames/gner/ent/token"
)

// Tokenize creates a slice containing every word in the document tokenized.
func Tokenize(text []rune) []TokenN {
	gts := gner.Tokenize(text, New)
	res := make([]TokenN, len(gts))
	for i := range gts {
		t := gts[i].(TokenN)
		res[i] = t
	}
	return res
}
