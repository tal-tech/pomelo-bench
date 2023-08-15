package metrics

import (
	"encoding/json"
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
