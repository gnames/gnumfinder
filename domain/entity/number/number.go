package number

import (
	"log"

	"github.com/gnames/gner/domain/entity/token"
	"github.com/gnames/gnlib/encode"
)

type Number struct {
	Raw         string `json:"raw"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
	Number      int    `json:"number"`
	NumRestored int    `json:"numRestored"`
	MaybePage   bool   `json:"maybePage"`
	MaybeYear   bool   `json:"maybeYear"`
}

func NewFindResult(t token.Token) Number {
	if !t.HasDigits {
		log.Fatalf("Token %s has no number in it", string(t.Raw))
	}
	return Number{
		Raw:   string(t.Raw),
		Start: t.Start,
		End:   t.End,
	}
}

func (num Number) ToJSON(pretty bool) ([]byte, error) {
	enc := encode.GNjson{Pretty: pretty}
	return enc.Encode(num)
}
