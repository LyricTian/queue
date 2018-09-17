package queue

import (
	"fmt"
	"sync/atomic"
	"testing"
)

type testJob struct {
	payload  int
	callback func(int)
}

func (t *testJob) Job() {
	t.payload++
	t.callback(t.payload)
}

func TestListQueue(t *testing.T) {
	q := NewListQueue(2)
	q.Run()

	var data int64
	q.Push(&testJob{
		payload: 0,
		callback: func(result int) {
			atomic.AddInt64(&data, int64(result))
		},
	})
	q.Push(&testJob{
		payload: 0,
		callback: func(result int) {
			atomic.AddInt64(&data, int64(result))
		},
	})
	q.Push(&testJob{
		payload: 0,
		callback: func(result int) {
			atomic.AddInt64(&data, int64(result))
		},
	})
	q.Terminate()
	if data != 3 {
		t.Error(data)
	}
}

func ExampleListQueue() {
	q := NewListQueue(10)
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

func BenchmarkListQueue(b *testing.B) {
	q := NewListQueue(100)
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
