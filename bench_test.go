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

func runmulti(w *Wheel) {
	for i := 0; i < 10; i++ {
		go func() {
			w.AfterFunc(20*time.Millisecond, func() {
			})
		}()
	}
	time.Sleep(25 * time.Millisecond)
}
func BenchmarkMultiAfter(b *testing.B) {
	w := NewDefaultWheel()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runmulti(w)
	}
}
