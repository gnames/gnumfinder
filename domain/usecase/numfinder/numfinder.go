package numfinder

import (
	"github.com/gnames/gner/domain/entity/token"
	"github.com/gnames/gner/domain/entity/txt"
	"github.com/gnames/gnumfind/domain/entity/number"
)

type numfinder struct{}

func NewNERecognizer() numfinder {
	return numfinder{}
}

func (n numfinder) Find(tn txt.TextNER) {
	tokens := token.Tokenize(tn.GetText())
	nums := make([]txt.EntityNER, 0, len(tokens))
	lines := make(map[int]int)
	for _, v := range tokens {
		lines[v.Line] += 1
		if v.HasDigits {
			nums = append(nums, number.NewNumber(v))
		}
	}
	tn.SetLines(lines)
	tn.SetEntities(nums)
}

func (n numfinder) FindInVolume(vol txt.VolumeNER) {
	pages := vol.GetPages()
	for _, p := range pages {
		n.Find(p)
	}
	pageNums(vol)
}

type pgNums struct {
	page         txt.PageNER
	pageNumFirst number.Number
	pageNumLast  number.Number
}

func pageNums(vol txt.VolumeNER) {
	pages := vol.GetPages()
	pns := make([]pgNums, len(pages))
	for i, page := range pages {
		pn := pgNums{page: page}
		markPageNums(&pn)
		pns[i] = pn
	}
}

func markPageNums(pn *pgNums) {
	ents := pn.page.GetEntities()
	if len(ents) == 0 {
		return
	}

	firstNum := ents[0].(number.Number)
	pn.pageNumFirst = firstNum
	lastNum := ents[len(ents)-1].(number.Number)
	pn.pageNumLast = lastNum
}
