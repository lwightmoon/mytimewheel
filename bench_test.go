package mytimewheel

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkAfter(b *testing.B) {
	w := NewWheel(1, 1000)
	for i := 0; i < b.N; i++ {
		d := rand.Intn(10000) + 1
		w.AfterFunc(time.Duration(d)*time.Millisecond, func() {
		})
	}
}
