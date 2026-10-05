package metrics

import (
	"fmt"
	"os"
	"time"
)

// quick-and-dirty, very hacky, recorder for per-tick performance metrics

type metricRecorder struct {
	data []*tickMetricLine
}

type tickMetricLine struct {
	tickNumber       int
	duration         time.Duration
	changesProcessed int
}

var globalRecorder = new()

func new() *metricRecorder {
	return &metricRecorder{
		data: make([]*tickMetricLine, 0, 200), // allocate some initial capacity
	}
}

func RecordTick(tickNumber int, duration time.Duration, changesProcessed int) {
	globalRecorder.data = append(globalRecorder.data, &tickMetricLine{
		tickNumber,
		duration,
		changesProcessed,
	})
}

func WriteToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Error creating file: %v", err))
	}
	defer file.Close()

	for _, line := range globalRecorder.data {
		lineString := fmt.Sprintf("%v,%v,%v\n", line.tickNumber, line.duration.Microseconds(), line.changesProcessed)
		file.WriteString(lineString)
	}
	file.Sync()
}
