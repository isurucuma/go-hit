package hit

import (
	"net/http"
	"time"
)

type Result struct {
	RPS      float64
	Requests int
	Errors   int
	Bytes    int64
	Duration time.Duration
	Fastest  time.Duration
	Slowest  time.Duration
	Status   int
	Error    error
}

func Send(r *http.Request) (Result, error) {
	t := time.Now()

	time.Sleep(100 * time.Millisecond) // temperory for simulation process only

	result := Result{
		Duration: time.Since(t),
		Bytes:    10,
		Status:   http.StatusOK,
	}

	return result, nil
}
