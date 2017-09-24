package queue

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWorker(t *testing.T) {
	var done interface{}
	pool := make(chan chan Jober, 1)
	w := NewWorker(pool, func() {
		done = "done"
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

	if !reflect.DeepEqual(done, "done") {
		t.Error(done)
	}
}
