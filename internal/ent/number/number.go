package number

import (
	"log"
	"strconv"
	"time"
	"unicode"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnumfinder/internal/ent/token"
)

type state int

const (
	enter    state = iota
	numStart state = iota
	restoreStart
	exit
)

var char2digit = map[rune]rune{
	'o': '0',
	'O': '0',
	'l': '1',
	'I': '1',
	'з': '3',
	'З': '3',
	'э': '3',
	'Э': '3',
	'g': '9',
}

var maxYear = time.Now().Year() + 2

type Number struct {
	token.TokenN `json:"-"`
	Raw          string `json:"raw"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
	Number       int    `json:"number"`
	NumRestored  int    `json:"numRestored"`
	MaybePage    bool   `json:"maybePage"`
	MaybeYear    bool   `json:"maybeYear"`
}

func NewNumber(t token.TokenN) Number {
	f := t.Features()
	if !f.HasDigits {
		log.Fatalf("Token %s has no number in it", string(t.Raw()))
	}
	num, numRestored := cleanNumber(t)
	res := Number{
		TokenN:      t,
		Raw:         string(t.Raw()),
		Start:       t.Start(),
		End:         t.End(),
		Number:      num,
		NumRestored: numRestored,
	}
	res.MaybeYear = decideIfYear(res)
	return res
}

func (n Number) ToJSON(pretty bool) ([]byte, error) {
	enc := gnfmt.GNjson{Pretty: pretty}
	return enc.Encode(n)
}

func cleanNumber(t token.TokenN) (int, int) {
	f := t.Features()
	if f.IsNumber {
		num, _ := strconv.Atoi(t.Cleaned())
		return num, num
	}
	cleaned := []rune(t.Cleaned())
	num := make([]rune, 0, len(cleaned))
	numRest := make([]rune, 0, len(cleaned))
	state := enter
	for _, v := range cleaned {
		switch state {
		case enter:
			if !unicode.IsDigit(v) {
				continue
			}
			num = append(num, v)
			numRest = append(numRest, v)
			state = numStart
		case numStart:
			if unicode.IsDigit(v) {
				num = append(num, v)
				numRest = append(numRest, v)
				continue
			}
			state = restoreDigit(&numRest, v)
		case restoreStart:
			state = restoreDigit(&numRest, v)
		case exit:
			break
		}
	}
	number, _ := strconv.Atoi(string(num))
	numberRestored, _ := strconv.Atoi(string(numRest))
	return number, numberRestored
}

func restoreDigit(numRest *[]rune, r rune) state {
	if unicode.IsDigit(r) {
		*numRest = append(*numRest, r)
		return restoreStart
	}

	if digit, ok := char2digit[r]; ok {
		*numRest = append(*numRest, digit)
		return restoreStart
	}
	return exit
}

func decideIfYear(n Number) bool {
	f := n.Features()
	if !(f.IsNumber || len(strconv.Itoa(n.NumRestored)) == len(n.Cleaned())) {
		return false
	}
	if n.NumRestored <= maxYear && n.NumRestored > 1749 {
		return true
	}
	return false
}
