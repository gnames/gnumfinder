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
	nr = numfinder.NewNERecognizer()
	textRunes := []rune("123 one two three")
	text := txt.NewTextNER(textRunes)
	nr.Find(text)
	ents := text.GetEntities()
	nums := make([]number.Number, len(ents))
	for i, v := range ents {
		num, ok := v.(number.Number)
		assert.True(t, ok)
		nums[i] = num
	}
	assert.Equal(t, len(nums), 1)
	assert.Equal(t, nums[0].Raw, "123")
}

func TestCleanNum(t *testing.T) {
	var nr ner.NERecognizer
	nr = numfinder.NewNERecognizer()
	textRunes := []rune("123 ab12зo2I1 062Oabc7")
	text := txt.NewTextNER(textRunes)
	nr.Find(text)
	ents := text.GetEntities()
	nums := make([]number.Number, len(ents))
	for i, v := range ents {
		num, ok := v.(number.Number)
		assert.True(t, ok)
		nums[i] = num
	}
	assert.Equal(t, nums[0].Number, 123)
	assert.Equal(t, nums[0].NumRestored, 123)
	assert.Equal(t, nums[1].Number, 12)
	assert.Equal(t, nums[1].NumRestored, 1230211)
	assert.Equal(t, nums[2].Number, 62)
	assert.Equal(t, nums[2].NumRestored, 620)
}

func TestYear(t *testing.T) {
	var nr ner.NERecognizer
	nr = numfinder.NewNERecognizer()
	textRunes := []rune("123 ab1830FD 18og, 2008, 2005b")
	text := txt.NewTextNER(textRunes)
	nr.Find(text)
	ents := text.GetEntities()
	nums := make([]number.Number, len(ents))
	for i, v := range ents {
		num, ok := v.(number.Number)
		assert.True(t, ok)
		nums[i] = num
	}
	assert.False(t, nums[0].MaybeYear)
	assert.False(t, nums[1].MaybeYear)
	assert.True(t, nums[2].MaybeYear)
	assert.True(t, nums[3].MaybeYear)
	assert.False(t, nums[4].MaybeYear)
}

func TestVolume(t *testing.T) {
	var nr ner.NERecognizer
	nr = numfinder.NewNERecognizer()
	vol := volumeTest1(t)
	nr.FindInVolume(vol)
	outPages := vol.GetPages()
	assert.Equal(t, len(outPages), 99)
}

func volumeTest1(t *testing.T) txt.VolumeNER {
	var pages []txt.PageNER
	res := txt.NewVolumeNER("test1")
	path := filepath.Join("..", "..", "..", "testdata", "test1")
	files, err := ioutil.ReadDir(path)
	pages = make([]txt.PageNER, len(files))
	assert.Nil(t, err)
	for i, v := range files {
		id := v.Name()
		filePath := filepath.Join(path, id)
		text, err := ioutil.ReadFile(filePath)
		assert.Nil(t, err)
		page := txt.NewPageNER(id, []rune(string(text)))
		pages[i] = page
	}
	res.SetPages(pages)
	return res
}
