package hit

import (
	"io"
	"net/http"
	"time"
)

func Send(c *http.Client, r *http.Request) (Result, error) {
	t := time.Now()

	var (
		code  int
		bytes int64
	)

	response, err := c.Do(r)
	if err == nil {
		code = response.StatusCode
		bytes, err = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}

	result := Result{
		Duration: time.Since(t),
		Bytes:    bytes,
		Status:   code,
		Error:    err,
	}
	return result, err
}
