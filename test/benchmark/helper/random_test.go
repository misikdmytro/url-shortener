package random_test

import (
	"testing"

	"github.com/misikdmytro/url-shortener/internal/helper"
)

func BenchmarkNewKey(b *testing.B) {
	f := helper.NewRandomGeneratorFactory()
	g := f.NewRandomGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.NewKey(8)
	}
}
