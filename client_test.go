package hit

import (
	"context"
	"net/http"
	"testing"
)

type fakeTripper func(*http.Request) (*http.Response, error)

func (f fakeTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestClientdDo(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatalf("new http request: %v", err)
	}

	fail := func(req *http.Request) (*http.Response, error) {
		t.Logf("fakeTripper: %s", req.URL)
		resp := &http.Response{
			StatusCode: http.StatusInternalServerError,
		}
		return resp, nil
	}

	client := &Client{
		C:         1,
		RPS:       1,
		Transport: fakeTripper(fail),
	}

	sum := client.Do(context.Background(), req, 5)
	if got := sum.Errors; got != 5 {
		t.Errorf("Errors=%d; want 5", got)
	}
}
