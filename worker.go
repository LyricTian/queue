package queue

import (
	"sync/atomic"
)

// Worker a worker who performs a task
type Worker interface {
	Start()
	Terminate()
}

type worker struct {
	pool       chan<- chan Jober
	jobChannel chan Jober
	quit       chan struct{}
	done       func()
	running    uint32
}

// NewWorker Create a worker who performs a task,
// specify a work pool and a callback after the task completes
func NewWorker(pool chan<- chan Jober, done func()) Worker {
	return &worker{
		pool:       pool,
		jobChannel: make(chan Jober),
		quit:       make(chan struct{}),
		done:       done,
	}
}

func (w *worker) Start() {
	if atomic.LoadUint32(&w.running) == 1 {
		return
	}
	atomic.StoreUint32(&w.running, 1)

	go func() {
	LBQUIT:
		for {
			w.pool <- w.jobChannel

			select {
			case jober := <-w.jobChannel:
				jober.Job()
				if fn := w.done; fn != nil {
					fn()
				}
			case <-w.quit:
				break LBQUIT
			}
		}

		close(w.jobChannel)
	}()
}

func (w *worker) Terminate() {
	if atomic.LoadUint32(&w.running) != 1 {
		return
	}

	atomic.StoreUint32(&w.running, 0)
	close(w.quit)
}
