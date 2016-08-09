package pool

import "testing"

func TestBasic(t *testing.T) {
	count := 1000
	p := new(WorkerPool)
	for i := 0; i < count; i++ {
		p.Do(func(n uint64) {
			// NOOP
		})
	}
	if n := p.Done(); n != uint64(count) {
		t.Fatalf("Extected %v, got %v", count, n)
	}
}
