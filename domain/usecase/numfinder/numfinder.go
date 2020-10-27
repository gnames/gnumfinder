package numfinder

import (
	"sync"

	"github.com/gnames/gner/domain/entity/token"
	"github.com/gnames/gner/domain/entity/txt"
	"github.com/gnames/gnumfind/domain/entity/number"
)

type numfinder struct {
	jobs int
}

func NewNERecognizer() numfinder {
	return numfinder{jobs: 8}
}

func (n numfinder) Find(text []rune) []txt.OutputNER {
	tokens := token.Tokenize(text)
	res := make([]txt.OutputNER, 0, len(tokens))
	for _, v := range tokens {
		if v.HasDigits {
			res = append(res, number.NewFindResult(v))
		}
	}
	return res
}

func (n numfinder) FindInVolume(vol txt.Volume) txt.OutputVolumeNER {
	res := number.NewOutputVolumeNER(vol.ID)
	res.OutPages = make([]txt.OutputPageNER, 0, len(vol.Pages))
	chIn := make(chan txt.Page)
	chOut := make(chan txt.OutputPageNER)
	var wgWork, wgProc sync.WaitGroup
	wgWork.Add(n.jobs)
	wgProc.Add(1)

	go func() {
		for _, p := range vol.Pages {
			chIn <- p
		}
		close(chIn)
	}()
	for i := 0; i < n.jobs; i++ {
		go n.finderWorker(chIn, chOut, &wgWork)
	}
	go func() {
		defer wgProc.Done()
		for out := range chOut {
			res.OutPages = append(res.OutPages, out)
		}
	}()
	wgWork.Wait()
	close(chOut)
	wgProc.Wait()
	return res
}

func (n numfinder) finderWorker(chIn <-chan txt.Page, chOut chan<- txt.OutputPageNER,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range chIn {
		out := number.NewOutputPageNER(p)
		out.Output = n.Find(p.Text)
		chOut <- out
	}
}
