package pool

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	wg   sync.WaitGroup
	once sync.Once
	jobs chan func(uint64)
	n    uint64
}

func (p *WorkerPool) Do(f func(n uint64)) {
	p.once.Do(func() {
		p.jobs = make(chan func(uint64), runtime.NumCPU()*2)
		for i := 0; i < runtime.NumCPU(); i++ {
			p.wg.Add(1)
			go func() {
				defer p.wg.Done()
				var n uint64
				defer func() {
					atomic.AddUint64(&p.n, n)
				}()
				for job := range p.jobs {
					job(n)
					n++
				}
			}()
		}
	})
	p.jobs <- f
}

func (p *WorkerPool) Done() uint64 {
	if p.jobs != nil {
		close(p.jobs)
		p.wg.Wait()
	}
	return p.n
}
