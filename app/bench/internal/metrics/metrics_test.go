package metrics

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestMetrics_Execute(t *testing.T) {

	m := NewMetrics("test")

	n := time.Now()

	for i := 0; i < 100; i++ {

		if i == 99 {

			time.Sleep(1 * time.Second)
		}
		m.Add(time.Now().Sub(n))

	}

	report := m.Execute()

	d, _ := json.Marshal(report)
	t.Log(string(d))

}

func TestMetrics_AppendExecute(t *testing.T) {

	m := NewMetrics("test")

	n := time.Now()

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {

		wg.Add(1)

		go func(index int) {
			if index == 99 {

				time.Sleep(1 * time.Second)
			}

			m.Add(time.Now().Sub(n))

			wg.Done()

		}(i)

	}

	wg.Wait()

	report := m.Execute()

	d, _ := json.Marshal(report)
	t.Log(string(d))

	t.Log(len(m.tasks))

}

func Benchmark_ADD(b *testing.B) {

	me := SimpleMetrics{}
	n := time.Now()

	for i := 0; i < b.N; i++ {

		me.Add(time.Now().Sub(n))

	}

	b.Log(me.Duration())
}
