package queue

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestSyncJob(t *testing.T) {
	job := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return fmt.Sprintf("%s_bar", v), nil
	})

	go func() {
		time.Sleep(time.Millisecond * 10)
		job.Exec()
	}()

	result := <-job.Wait()
	if err := job.Error(); err != nil {
		t.Error(err.Error())
		return
	}

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}
}

func TestSyncJobError(t *testing.T) {
	job := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return nil, errors.New("wow")
	})

	go func() {
		time.Sleep(time.Millisecond * 10)
		job.Exec()
	}()

	result := <-job.Wait()
	if err := job.Error(); err == nil {
		t.Error("error is wow")
		return
	}

	if result != nil {
		t.Error("result is nil")
	}
}
