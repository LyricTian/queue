package queue

var (
	globalQueue *Queue
)

// Run start run
func Run(maxQueues, maxWorkers int) {
	if globalQueue == nil {
		globalQueue = NewQueue(maxQueues, maxWorkers)
	}
	globalQueue.Run()
}

// Stop stop run
func Stop() {
	if globalQueue == nil {
		return
	}
	globalQueue.Stop()
}

// Push add job to queue
func Push(job Job) {
	if globalQueue == nil {
		return
	}
	globalQueue.Push(job)
}

// Queue ...
type Queue struct {
	maxWorkers int
	jobQueue   chan Job
	workerPool chan chan Job
	workers    []Worker
	quit       chan bool
	isRunning  bool
}

// NewQueue create a Queue
func NewQueue(maxQueues, maxWorkers int) *Queue {
	return &Queue{
		jobQueue:   make(chan Job, maxQueues),
		maxWorkers: maxWorkers,
		workerPool: make(chan chan Job, maxWorkers),
		workers:    make([]Worker, maxWorkers),
		quit:       make(chan bool),
	}
}

// Run start run
func (q *Queue) Run() {
	if q.isRunning {
		return
	}

	q.isRunning = true
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i] = NewWorker(q.workerPool)
		q.workers[i].Start()
	}
	q.dispatcher()
}

func (q *Queue) dispatcher() {
	go func() {
		for {
			select {
			case job := <-q.jobQueue:
				worker := <-q.workerPool
				worker <- job
			case <-q.quit:
				q.isRunning = false
				return
			}
		}
	}()
}

// Stop stop run
func (q *Queue) Stop() {
	if !q.isRunning {
		return
	}

	go func() {
		q.quit <- true
	}()

	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i].Stop()
	}
}

// Push add job to queue
func (q *Queue) Push(job Job) {
	if !q.isRunning {
		return
	}
	q.jobQueue <- job
}
