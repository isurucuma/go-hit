package hit

import (
	"net/http"
	"time"
)

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
