package queue

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
)

func TestWorker(t *testing.T) {
	var done int64
	pool := make(chan chan Jober, 1)
	w := NewWorker(pool, func() {
		atomic.AddInt64(&done, 1)
	})

	w.Start()

	defer w.Terminate()

	sjob := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return fmt.Sprintf("%s_bar", v), nil
	})

	<-pool <- sjob

	result := <-sjob.Wait()
	if err := sjob.Error(); err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}

	if atomic.LoadInt64(&done) != 1 {
		t.Error(done)
	}
}
