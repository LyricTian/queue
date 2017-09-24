package queue

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
	isRuning   bool
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
	if w.isRuning {
		return
	}

	w.isRuning = true
	go func() {
		for w.isRuning {
			w.pool <- w.jobChannel

			select {
			case jober := <-w.jobChannel:
				jober.Job()
				if fn := w.done; fn != nil {
					fn()
				}
			case <-w.quit:
			}

		}

		close(w.jobChannel)
	}()
}

func (w *worker) Terminate() {
	if !w.isRuning {
		return
	}

	w.isRuning = false
	close(w.quit)
}
