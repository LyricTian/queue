package queue

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue(1, 1)
	q.Run()

	sjob := NewSyncJob("foo", func(v interface{}) (interface{}, error) {
		return fmt.Sprintf("%s_bar", v), nil
	})
	q.Push(sjob)

	result := <-sjob.Wait()
	if err := sjob.Error(); err != nil {
		t.Error(err.Error())
		return
	}

	if !reflect.DeepEqual(result, "foo_bar") {
		t.Error(result)
	}
}

type tjob struct {
	wg  *sync.WaitGroup
	buf []byte
}

func (t *tjob) Exec() {
	_ = t.buf
	time.Sleep(time.Millisecond)
	t.wg.Done()
}

func BenchmarkQueue(b *testing.B) {
	wg := new(sync.WaitGroup)

	q := NewQueue(10, 10)
	q.Run()

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

	go func() {
		for buf := range ch {
			_ = buf
			time.Sleep(time.Millisecond)
			wg.Done()
		}
	}()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			ch <- []byte("fooooooooooooo")
		}
	})

	wg.Wait()
}
