package gnumfinder

import "github.com/gnames/gner/ent/txt"

type GNumfinder interface {
	Find(txt.TextNER)
	FindInVolume(txt.VolumeNER)
}
