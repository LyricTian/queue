package queue

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWorker(t *testing.T) {
	pool := make(chan chan Job, 1)
	w := NewWorker(pool)

	w.Start()

	sjob := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return fmt.Sprintf("%s_bar", v), nil
	})

	go func() {
		job := <-pool
		job <- sjob
	}()

	result := <-sjob.Wait()
	if err := sjob.Error(); err != nil {
		t.Error(err.Error())
		return
	}

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}
}
