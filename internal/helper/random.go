package helper

import (
	"math/rand"
	"time"
)

type RandomGenerator interface {
	NewKey(length int) string
}

type RandomGeneratorFactory interface {
	NewRandomGenerator() RandomGenerator
}

type randomGenerator struct {
	r *rand.Rand
}

type randomGeneratorFactory struct{}

func NewRandomGeneratorFactory() RandomGeneratorFactory {
	return &randomGeneratorFactory{}
}

func (f *randomGeneratorFactory) NewRandomGenerator() RandomGenerator {
	return &randomGenerator{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (r *randomGenerator) NewKey(length int) string {
	key := make([]byte, length)
	for i := range key {
		key[i] = symbols[r.r.Intn(len(symbols))]
	}
	return string(key)
}
