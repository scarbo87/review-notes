package tricks

import (
	"fmt"
)

type worker struct {
	jobFn   func(job int)
	jobCh   chan int
	closeCh chan struct{}
}

func NewWorker(size uint64, jobFn func(job int)) *worker {

	if jobFn == nil {
		jobFn = func(job int) {
			fmt.Println(job)
		}
	}

	w := &worker{
		jobCh:   make(chan int, size),
		closeCh: make(chan struct{}),
		jobFn:   jobFn,
	}
	go w.run()
	return w
}

func (w *worker) AddJob(job int) {
	w.jobCh <- job
}

// Close implement a waiting close strategy.
// We send a signal to special not-buffering channel and then we are waiting when this channel will be closed.
func (w *worker) Close() {
	w.closeCh <- struct{}{}
	<-w.closeCh
}

func (w *worker) run() {
	defer func() {
		for {
			// if we haven't processed  all tasks yet, we'll do it here
			if len(w.jobCh) == 0 {
				break
			}
			w.jobFn(<-w.jobCh)
		}
		close(w.jobCh)
		close(w.closeCh)
	}()

	for {
		select {
		case job := <-w.jobCh:
			w.jobFn(job)
		case <-w.closeCh: // when we received a signal we'll call the defer function
			return
		}
	}
}
