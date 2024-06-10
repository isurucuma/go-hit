package hit

import (
	"net/http"
	"sync"
	"time"
)

// produces calls fn n times and sends the results to the out
func Produce(out chan<- *http.Request, n int, fn func() *http.Request) {
	for range n {
		out <- fn()
	}
}

func produce(n int, fn func() *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Produce(out, n, fn)
	}()
	return out
}

// Throttle slows down receiving from in by delay and sends what it receives from in to out
func Throttle(in <-chan *http.Request, out chan<- *http.Request, delay time.Duration) {
	t := time.NewTicker(delay)

	defer t.Stop() // to relaese the resources

	// until in channel stopes this will send requests to the out channel but with a minimim delay of delay
	for r := range in {
		<-t.C // wait until timer tick comes
		out <- r
	}
}

// throttle runs Throttle in a go routine
func throttle(in <-chan *http.Request, delay time.Duration) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Throttle(in, out, delay)
	}()
	return out
}

// this is the type which will send the http request and then returns a Result type
type sendFunc func(*http.Request) Result

func Split(in <-chan *http.Request, out chan<- Result, c int, fn sendFunc) {
	// c number of goroutines will concurrently try to get requests from the in channel and send the requests
	send := func() {
		for r := range in {
			out <- fn(r)
		}
	}

	var wg sync.WaitGroup
	wg.Add(c)
	for range c {
		go func() {
			defer wg.Done()
			send()
		}()
	}
	wg.Wait()
}

func split(in <-chan *http.Request, c int, fn sendFunc) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)
		Split(in, out, c, fn)
	}()
	return out
}
