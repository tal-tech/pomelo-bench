package metrics

import (
	"time"
)

// A Task is a task that is reported to Metrics.
type Task struct {
	Drop     bool
	Duration time.Duration
	//Description string
}

type Metrics struct {
	name string
	//pid      int
	tasks    []Task
	duration time.Duration // 总耗时
	drops    int
}

type StatReport struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"tm"`
	//Pid       int    `json:"pid"`
	//ReqsPerSecond float32 `json:"qps"`
	Drops     int     `json:"drops"`
	Average   float32 `json:"avg"`
	Median    float32 `json:"med"`
	Top90th   float32 `json:"t90"`
	Top99th   float32 `json:"t99"`
	Top99p9th float32 `json:"t99p9"`
}

func NewMetrics(name string) *Metrics {
	return &Metrics{
		name:     name,
		tasks:    make([]Task, 0, 256),
		duration: 0,
		drops:    0,
	}
}

func (m *Metrics) Execute() StatReport {
	tasks := m.tasks
	duration := m.duration
	drops := m.drops
	size := len(tasks)
	report := StatReport{
		Name:      m.name,
		Timestamp: time.Now().Unix(),
		//Pid:       c.pid,
		//ReqsPerSecond: float32(size) / float32(logInterval/time.Second),
		Drops: drops,
	}

	if size > 0 {
		report.Average = float32(duration/time.Millisecond) / float32(size)

		fiftyPercent := size >> 1
		if fiftyPercent > 0 {
			top50pTasks := topK(tasks, fiftyPercent)
			medianTask := top50pTasks[0]
			report.Median = float32(medianTask.Duration) / float32(time.Millisecond)
			tenPercent := fiftyPercent / 5
			if tenPercent > 0 {
				top10pTasks := topK(top50pTasks, tenPercent)
				task90th := top10pTasks[0]
				report.Top90th = float32(task90th.Duration) / float32(time.Millisecond)
				onePercent := tenPercent / 10
				if onePercent > 0 {
					top1pTasks := topK(top10pTasks, onePercent)
					task99th := top1pTasks[0]
					report.Top99th = float32(task99th.Duration) / float32(time.Millisecond)
					pointOnePercent := onePercent / 10
					if pointOnePercent > 0 {
						topPointOneTasks := topK(top1pTasks, pointOnePercent)
						task99Point9th := topPointOneTasks[0]
						report.Top99p9th = float32(task99Point9th.Duration) / float32(time.Millisecond)
					} else {
						report.Top99p9th = getTopDuration(top1pTasks)
					}
				} else {
					mostDuration := getTopDuration(top10pTasks)
					report.Top99th = mostDuration
					report.Top99p9th = mostDuration
				}
			} else {
				mostDuration := getTopDuration(top50pTasks)
				report.Top90th = mostDuration
				report.Top99th = mostDuration
				report.Top99p9th = mostDuration
			}
		} else {
			mostDuration := getTopDuration(tasks)
			report.Median = mostDuration
			report.Top90th = mostDuration
			report.Top99th = mostDuration
			report.Top99p9th = mostDuration
		}
	}

	return report
}

func (m *Metrics) Drop() {
	m.drops++
}

func (m *Metrics) Clear() {

	m.tasks = make([]Task, 0, 256)
	m.duration = 0
	m.drops = 0
}

func (m *Metrics) Add(duration time.Duration) {
	m.tasks = append(m.tasks, Task{
		Drop:     false,
		Duration: duration,
	})

	m.duration += duration
}

func getTopDuration(tasks []Task) float32 {
	top := topK(tasks, 1)
	if len(top) < 1 {
		return 0
	}

	return float32(top[0].Duration) / float32(time.Millisecond)
}
