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

	sum, err := SendN(context.Background(), server.URL, wantHits, Concurrency(1), RequestsPerSecond(10))
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
