package pool

import (
	"runtime"
	"sync"
)

type WorkerPool struct {
	wg   sync.WaitGroup
	once sync.Once
	jobs chan func()
}

func (p *WorkerPool) Do(f func()) {
	p.once.Do(func() {
		p.jobs = make(chan func(), runtime.NumCPU()*2)
		for i := 0; i < runtime.NumCPU(); i++ {
			p.wg.Add(1)
			go func() {
				defer p.wg.Done()
				for job := range p.jobs {
					job()
				}
			}()
		}
	})
	p.jobs <- f
}

func (p *WorkerPool) Done() {
	if p.jobs != nil {
		close(p.jobs)
		p.wg.Wait()
	}
}
