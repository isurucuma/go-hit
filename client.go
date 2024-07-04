package hit

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// client sends HTTP requests and returns an aggregated result
type Client struct {
	C         int               // concurrency level
	RPS       int               // no of requests per second
	Timeout   time.Duration     // timeout per request
	Transport http.RoundTripper // use this to send the request and get the response
}

func (c *Client) client() *http.Client {
	transport := c.Transport
	if transport == nil {
		transport = &http.Transport{
			MaxIdleConnsPerHost: c.concurrency(), // otherwise by default http client keeps less number of idle connections in the pool
		}
	}
	return &http.Client{
		Timeout: c.Timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}
}

// Do sends n HTTP requests and returns an aggregated result
func (c *Client) Do(ctx context.Context, r *http.Request, n int) Result {
	t := time.Now()
	var sum Result
	for result := range c.do(ctx, r, n) {
		sum = sum.Merge(result)
	}
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(ctx context.Context, r *http.Request, n int) <-chan Result {
	pipe := produce(ctx, n, func() *http.Request { return r.Clone(ctx) })

	if c.RPS > 0 {
		t := time.Second / time.Duration(c.RPS*c.concurrency())
		pipe = throttle(pipe, t)
	}
	return split(pipe, c.concurrency(), func(r *http.Request) Result {
		// skipping the error handling as it is already handled in the performance result summary
		result, _ := Send(c.client(), r)
		return result
	})
}

func (c *Client) concurrency() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}

type Option func(*Client)

func Concurrency(n int) Option {
	return func(c *Client) {
		c.C = n
	}
}

func RequestsPerSecond(n int) Option {
	return func(c *Client) {
		c.RPS = n
	}
}

func Timeout(d time.Duration) Option {
	return func(c *Client) {
		c.Timeout = d
	}
}

func SendN(ctx context.Context, url string, headers map[string]string, n int, opts ...Option) (Result, error) {
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return Result{}, fmt.Errorf("new http request: %w", err)
	}
	for key, val := range headers {
		r.Header.Add(key, val)
	}
	var c Client
	for _, o := range opts {
		o(&c)
	}
	return c.Do(ctx, r, n), nil
}
