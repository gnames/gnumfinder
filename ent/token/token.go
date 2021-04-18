package token

import (
	"strings"
	"unicode"

	gner "github.com/gnames/gner/ent/token"
)

type tokenN struct {
	gner.TokenNER
	features *Features
	runeSet  map[rune]struct{}
}

func New(token gner.TokenNER) gner.TokenNER {
	t := &tokenN{
		TokenNER: token,
	}
	return t
}

// Features is a fixed set of general properties determined during the
// the text traversal.
type Features struct {
	// HasStartParens token starts with '('.
	HasStartParens bool

	// HasEndParens token end with '('.
	HasEndParens bool

	// HasStartSqParens token starts with '['.
	HasStartSqParens bool

	// HasEndSqParens token ends with ']'.
	HasEndSqParens bool

	// HasEndDot token ends with '.'
	HasEndDot bool

	// HasEndComma token ends with ','
	HasEndComma bool

	// HasDigits token includes at least one '0-9'.
	HasDigits bool

	// HasLetters token includes at least one character for which
	// unicode.IsLetter(ch) is true.
	HasLetters bool

	// HasDash token includes '-'
	HasDash bool

	// HasSpecialChars internal part of a token includes non-letters, non-digits.
	HasSpecialChars bool

	// IsCapitalized is true if the furst letter of a token is capitalized.
	// The first letter does not have to be the first character.
	IsCapitalized bool

	// IsNumber internal part of a token has only numbers.
	IsNumber bool

	// IsWord internal part of a token includes only letters.
	IsWord bool
}

func (t *tokenN) Features() *Features {
	return t.features
}

func (t *tokenN) SetFeatures(f *Features) {
	t.features = f
}

// calculateProperties takes raw and cleaned values of a token and computes
// properties of these values, saving them into Properties object.
func calculateProperties(raw, cleaned []rune, f *Features) {
	f.HasStartParens = raw[0] == rune('(')
	f.HasEndParens = raw[len(raw)-1] == rune(')')
	f.HasStartSqParens = raw[0] == rune('[')
	f.HasEndSqParens = raw[len(raw)-1] == rune(']')
	f.HasEndDot = raw[len(raw)-1] == rune('.')
	f.HasEndComma = raw[len(raw)-1] == rune(',')
	for _, v := range cleaned {
		if v == rune('-') {
			f.HasDash = true
		}

		if !f.HasLetters && unicode.IsLetter(v) {
			f.HasLetters = true
			continue
		}

		if !f.HasDigits && unicode.IsDigit(v) {
			f.HasDigits = true
			continue
		}

		if !f.HasSpecialChars && v == rune('�') {
			f.HasSpecialChars = true
			continue
		}
	}

	if f.HasDigits && !f.HasLetters && !f.HasSpecialChars && !f.HasDash {
		f.IsNumber = true
	}

	if f.HasLetters && !f.HasDigits && !f.HasSpecialChars {
		f.IsWord = true
	}
}

func (t *tokenN) ProcessToken() {
	t.normalizeRaw()
	t.features = &Features{}
	calculateProperties(t.Raw(), []rune(t.Cleaned()), t.features)
}

func (t *tokenN) normalizeRaw() {
	var runes []rune
	t.runeSet = make(map[rune]struct{})
	firstLetter := true
	for _, v := range t.Raw() {
		t.runeSet[v] = struct{}{}
		hasDash := v == rune('-')
		if unicode.IsLetter(v) || unicode.IsNumber(v) || hasDash {
			if firstLetter {
				firstLetter = false
			}
			runes = append(runes, v)
		} else {
			t.runeSet['�'] = struct{}{}
			runes = append(runes, rune('�'))
		}
	}
	res := string(runes)
	t.SetCleaned(strings.Trim(res, "�"))
}
