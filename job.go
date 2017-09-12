package queue

// Job ...
type Job interface {
	Exec()
}

// SyncJob sync job
type SyncJob struct {
	err      error
	result   chan interface{}
	v        interface{}
	callback func(interface{}) (interface{}, error)
}

// NewSyncJob create a sync job
func NewSyncJob(v interface{}, callback func(interface{}) (interface{}, error)) *SyncJob {
	return &SyncJob{
		result:   make(chan interface{}, 1),
		v:        v,
		callback: callback,
	}
}

// Exec ...
func (j *SyncJob) Exec() {
	result, err := j.callback(j.v)
	if err != nil {
		j.err = err
		close(j.result)
		return
	}

	j.result <- result

	close(j.result)
}

// Wait ...
func (j *SyncJob) Wait() <-chan interface{} {
	return j.result
}

func (j *SyncJob) Error() error {
	return j.err
}
