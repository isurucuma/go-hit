package hit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestSendN(t *testing.T) {
	t.Parallel()
	var goHits atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			goHits.Add(1)
		},
	))
	defer server.Close()

	const wantHits, wantErrors = 10, 0
	headers := make(map[string]string)
	sum, err := SendN(context.Background(), server.URL, headers, wantHits, Concurrency(1), RequestsPerSecond(10))
	if err != nil {
		t.Fatalf("SendN() err=%q; want nil", err)
	}
	if got := goHits.Load(); got != wantHits {
		t.Errorf("hits=%d; want %d", got, wantHits)
	}
	if got := sum.Errors; got != wantErrors {
		t.Errorf("Errors=%d; want %d", got, wantErrors)
	}
}

func TestSendNWithValidHeaders(t *testing.T) {
	t.Parallel()
	var goHits atomic.Int64
	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer YOUR_ACCESS_TOKEN"}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			goHits.Add(1)
			for k, v := range headers {
				if r.Header.Get(k) != v {
					t.Fatalf("header %q=%q; want %q", k, r.Header.Get(k), v)
				}
			}
		},
	))
	defer server.Close()

	const wantHits, wantErrors = 1, 0
	sum, err := SendN(context.Background(), server.URL, headers, wantHits, Concurrency(1), RequestsPerSecond(1))
	if err != nil {
		t.Fatalf("SendN() err=%q; want nil", err)
	}
	if got := goHits.Load(); got != wantHits {
		t.Errorf("hits=%d; want %d", got, wantHits)
	}
	if got := sum.Errors; got != wantErrors {
		t.Errorf("Errors=%d; want %d", got, wantErrors)
	}
}
