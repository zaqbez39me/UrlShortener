package generator

import (
	"time"

	"math/rand"
)

type Generator interface {
	Generate(int) string
}

type RandomGenerator struct {
	alphabet []rune
}

func NewRandomGenerator(alphabet string) *RandomGenerator {
	return &RandomGenerator{alphabet: []rune(alphabet)}
}

func (g *RandomGenerator) Generate(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	genRunes := make([]rune, 0, size)
	var rndChar rune
	for len(genRunes) < size {
		rndChar = g.alphabet[rnd.Intn(len(g.alphabet))]
		genRunes = append(genRunes, rndChar)
	}

	return string(genRunes)
}
