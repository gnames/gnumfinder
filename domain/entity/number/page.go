package number

import (
	"github.com/gnames/gner/domain/entity/txt"
	"github.com/gnames/gnlib/encode"
)

type OutputPage struct {
	txt.Page
	Output []txt.OutputNER
}

func NewOutputPageNER(page txt.Page) OutputPage {
	return OutputPage{Page: page}
}

func (op OutputPage) PageID() string {
	return op.Page.ID
}

func (op OutputPage) ToJSON(pretty bool) ([]byte, error) {
	enc := encode.GNjson{Pretty: pretty}
	return enc.Encode(op)
}
