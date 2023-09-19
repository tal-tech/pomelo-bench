package metrics

import (
	"sync"
	"time"
)

type SimpleMetrics struct {
	duration int64 // 总耗时
	count    int64
	mu       sync.Mutex
}

func (s *SimpleMetrics) Add(duration time.Duration) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	s.duration += int64(duration)

}
func (s *SimpleMetrics) Duration() (duration time.Duration) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count == 0 {
		return 0
	}

	duration = time.Duration(s.duration / s.count)

	return duration
}

func (s *SimpleMetrics) Clear() {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count == 0 {
		return
	}

	s.count = 0
	s.duration = 0
}
