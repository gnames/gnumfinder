package token_test

import (
	"testing"

	"github.com/gnames/gnumfinder/ent/token"
	"github.com/tj/assert"
)

func TestTokenize(t *testing.T) {
	str := "one\vtwo 1779, poma-  \t\r\ntomus " +
		"dash -\nstandalone (23) " +
		"Tora-\nBora\n\rthree \n"
	tokens := token.Tokenize([]rune(str))
	assert.Equal(t, len(tokens), 10)
	for i := range tokens {
		if i == 2 || i == 7 {
			assert.True(t, tokens[i].Features().IsNumber)
		}
	}
}
