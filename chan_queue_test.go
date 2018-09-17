package queue

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue(2, 10)
	q.Run()

	var count int64

	for i := 0; i < 10; i++ {
		job := NewJob("foo", func(v interface{}) {
			atomic.AddInt64(&count, 1)
		})
		q.Push(job)
	}

	q.Terminate()

	if count != 10 {
		t.Error(count)
	}
}

func TestSyncQueue(t *testing.T) {
	q := NewQueue(1, 2)
	q.Run()
	defer q.Terminate()

	sjob := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return fmt.Sprintf("%s_bar", v), nil
	})
	q.Push(sjob)

	result := <-sjob.Wait()
	if err := sjob.Error(); err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}
}

func ExampleQueue() {
	q := NewQueue(1, 10)
	q.Run()

	var count int64

	for i := 0; i < 10; i++ {
		job := NewJob("foo", func(v interface{}) {
			atomic.AddInt64(&count, 1)
		})
		q.Push(job)
	}

	q.Terminate()
	fmt.Println(count)
	// output: 10
}

func BenchmarkQueue(b *testing.B) {
	q := NewQueue(10, 100)
	q.Run()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			job := NewJob("", func(v interface{}) {
				_ = v
			})
			q.Push(job)
		}
	})
	q.Terminate()
}
