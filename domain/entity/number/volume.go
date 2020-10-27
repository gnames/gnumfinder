package number

import (
	"github.com/gnames/gner/domain/entity/txt"
	"github.com/gnames/gnlib/encode"
)

type OutputVolume struct {
	ID       string
	OutPages []txt.OutputPageNER
}

func NewOutputVolumeNER(id string) OutputVolume {
	return OutputVolume{ID: id}
}

func (ov OutputVolume) VolumeID() string {
	return ov.ID
}

func (ov OutputVolume) OutputPages() []txt.OutputPageNER {
	res := make([]txt.OutputPageNER, len(ov.OutPages))
	for i, v := range ov.OutPages {
		res[i] = v
	}
	return res
}

func (ov OutputVolume) ToJSON(pretty bool) ([]byte, error) {
	enc := encode.GNjson{Pretty: pretty}
	return enc.Encode(ov)
}
