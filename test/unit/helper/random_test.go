package random_test

import (
	"testing"

	"github.com/misikdmytro/url-shortener/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestNewKey(t *testing.T) {
	f := helper.NewRandomGeneratorFactory()
	g := f.NewRandomGenerator()

	key := g.NewKey(16)
	assert.Len(t, key, 16)
}
