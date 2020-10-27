package numfinder_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gnames/gner/domain/entity/txt"
	"github.com/gnames/gner/domain/usecase/ner"
	"github.com/gnames/gnumfind/domain/entity/number"
	"github.com/gnames/gnumfind/domain/usecase/numfinder"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	var nr ner.NERecognizer
	text := "123 one two three"
	nr = numfinder.NewNERecognizer()
	recog := nr.Find([]rune(text))
	nums := make([]number.Number, len(recog))
	for i, v := range recog {
		num, ok := v.(number.Number)
		assert.True(t, ok)
		nums[i] = num
	}
	assert.Equal(t, len(nums), 1)
	assert.Equal(t, nums[0].Raw, "123")
}

func TestVolume(t *testing.T) {
	var nr ner.NERecognizer
	vol := volumeTest1(t)
	nr = numfinder.NewNERecognizer()
	out := nr.FindInVolume(vol)
	outPages := out.OutputPages()
	assert.Equal(t, len(outPages), 99)
}

func volumeTest1(t *testing.T) txt.Volume {
	var pages []txt.Page
	res := txt.Volume{ID: "test1"}
	path := filepath.Join("..", "..", "..", "testdata", "test1")
	files, err := ioutil.ReadDir(path)
	assert.Nil(t, err)
	for _, v := range files {
		id := v.Name()
		filePath := filepath.Join(path, id)
		text, err := ioutil.ReadFile(filePath)
		assert.Nil(t, err)
		page := txt.Page{ID: id, Text: []rune(string(text))}
		pages = append(pages, page)
	}
	res.Pages = pages
	return res
}
