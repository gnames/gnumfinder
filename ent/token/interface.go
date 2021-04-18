package token

import (
	gner "github.com/gnames/gner/ent/token"
)

type TokenN interface {
	gner.TokenNER
	Features() *Features
	SetFeatures(*Features)
}
