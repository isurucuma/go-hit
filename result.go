package hit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (r Result) Merge(other Result) Result {
	r.Requests++
	r.Bytes += other.Bytes

	if r.Fastest == 0 || other.Duration < r.Fastest {
		r.Fastest = other.Duration
	}
	if other.Duration > r.Slowest {
		r.Slowest = other.Duration
	}
	if other.Error != nil || other.Status != http.StatusOK {
		r.Errors++
	}
	return r
}

func (r Result) Finalize(total time.Duration) Result {
	r.Duration = total
	r.RPS = float64(r.Requests / int(total.Seconds()))
	return r
}

func (r Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}

	p("\nSummary:\n")
	p("\tSuccess \t: %.0f%%\n", r.successRatio())
	p("\tRPS \t\t: %.1f\n", r.RPS)
	p("\tRequests \t: %d\n", r.Requests)
	p("\tErrors \t\t: %d\n", r.Errors)
	p("\tBytes \t\t: %d\n", r.Bytes)
	p("\tDuration \t: %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastest \t: %s\n", round(r.Fastest))
		p("\tSlowest \t: %s\n", round(r.Slowest))
	}

}

func (r Result) String() string {
	var s strings.Builder
	r.Fprint(&s)
	return s.String()
}

func (r Result) successRatio() float64 {
	rr, e := float64(r.Requests), float64(r.Errors)
	return (rr - e) / rr * 100
}
func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)
}
