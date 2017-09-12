package queue

// Worker ...
type Worker interface {
	// start worker
	Start()
	// stop worker
	Stop()
}

type worker struct {
	pool       chan<- chan Job
	jobChannel chan Job
	quit       chan bool
}

// NewWorker create a Worker
func NewWorker(pool chan<- chan Job) Worker {
	return &worker{
		pool:       pool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w *worker) Start() {
	go func() {
		for {
			w.pool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				job.Exec()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
