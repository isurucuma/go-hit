package hit

import (
	"net/http"
	"time"
)

// client sends HTTP requests and returns an aggregated result
type Client struct {
	C   int // concurrency level
	RPS int // no of requests per second
}

// Do sends n HTTP requests and returns an aggregated result
func (c *Client) DO(r *http.Request, n int) Result {
	t := time.Now()

	var sum Result
	for range n {
		// final report will include all the errors included, therefore skipping the errors in the result here
		result, _ := Send(r)
		sum = sum.Merge(result)
	}
	return sum.Finalize(time.Since(t))
}
