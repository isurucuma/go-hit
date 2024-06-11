package hit

import (
	"context"
	"net/http"
	"runtime"
	"time"
)

// client sends HTTP requests and returns an aggregated result
type Client struct {
	C   int // concurrency level
	RPS int // no of requests per second
}

// Do sends n HTTP requests and returns an aggregated result
func (c *Client) Do(r *http.Request, n int) Result {
	t := time.Now()
	var sum Result
	for result := range c.do(r, n) {
		sum = sum.Merge(result)
	}
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(r *http.Request, n int) <-chan Result {
	pipe := produce(n, func() *http.Request { return r.Clone(context.TODO()) })

	if c.RPS > 0 {
		t := time.Second / time.Duration(c.RPS*c.concurrency())
		pipe = throttle(pipe, t)
	}
	return split(pipe, c.concurrency(), func(r *http.Request) Result {
		// skipping the error handling as it is already handled in the performance result summary
		result, _ := Send(r)
		return result
	})
}

func (c *Client) concurrency() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}
