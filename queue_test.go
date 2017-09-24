package queue

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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

type tjob struct {
	wg  *sync.WaitGroup
	buf []byte
}

func (t *tjob) Job() {
	_ = t.buf
	time.Sleep(time.Millisecond)
	t.wg.Done()
}

func BenchmarkQueue(b *testing.B) {
	wg := new(sync.WaitGroup)

	q := NewQueue(10, 100)
	q.Run()
	defer q.Terminate()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			q.Push(&tjob{buf: []byte("fooooooooooooo"), wg: wg})
		}
	})
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	var wg sync.WaitGroup

	ch := make(chan []byte, 10)

	for i := 0; i < 100; i++ {
		go func() {
			for buf := range ch {
				_ = buf
				time.Sleep(time.Millisecond)
				wg.Done()
			}
		}()
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			ch <- []byte("fooooooooooooo")
		}
	})

	wg.Wait()
}
