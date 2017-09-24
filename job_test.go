package queue

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestJob(t *testing.T) {
	var (
		result interface{}
		wg     sync.WaitGroup
	)

	wg.Add(1)
	job := NewJob("foo", func(v interface{}) {
		result = fmt.Sprintf("%s_bar", v)
		wg.Done()
	})

	go job.Job()
	wg.Wait()

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}
}

func TestSyncJob(t *testing.T) {
	sjob := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return fmt.Sprintf("%s_bar", v), nil
	})

	go sjob.Job()

	result := <-sjob.Wait()
	if err := sjob.Error(); err != nil {
		t.Error(err.Error())
		return
	}

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}
}

func TestSyncJobError(t *testing.T) {
	sjob := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return nil, errors.New("mock error")
	})

	go sjob.Job()

	result := <-sjob.Wait()
	if err := sjob.Error(); err == nil {
		t.Error("mock error")
		return
	}

	if result != nil {
		t.Error("result is nil")
	}
}
