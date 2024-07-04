package main

import (
	"io"
	"testing"
)

type parseArgsTest struct {
	name string
	args []string
	want Config
}

func TestParseArgsValidInput(t *testing.T) {
	t.Parallel()
	for _, tt := range []parseArgsTest{
		{
			name: "all_flags",
			args: []string{"-n=10", "-c=5", "-rps=5", "-H=Content-Type: application/json", "-H=Authorization: Bearer YOUR_ACCESS_TOKEN", "http://test"},
			want: Config{n: 10, c: 5, rps: 5, url: "http://test", headers: map[string]string{"Content-Type": "application/json", "Authorization": "Bearer YOUR_ACCESS_TOKEN"}},
		},
		// exercise: test with a mixture of flags
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got Config
			if err := ParseArgs(&got, tt.args, io.Discard); err != nil {
				t.Fatalf("parseArgs() error = %v, want no error", err)
			}
			if !got.isEqual(tt.want) {
				t.Errorf("flags = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func (c Config) isEqual(other Config) bool {
	headerEqual := true
	for k, v := range c.headers {
		if v != other.headers[k] {
			headerEqual = false
			break
		}
	}
	return c.url == other.url && c.n == other.n && c.c == other.c && c.rps == other.rps && headerEqual
}

func TestParseArgsInvalidInput(t *testing.T) {
	t.Parallel()
	for _, tt := range []parseArgsTest{
		{name: "n_syntax", args: []string{"-n=ONE", "http://test"}},
		{name: "n_zero", args: []string{"-n=0", "http://test"}},
		{name: "n_negative", args: []string{"-n=-1", "http://test"}},
		// exercise: test other error conditions
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ParseArgs(&Config{}, tt.args, io.Discard)
			if err == nil {
				t.Fatal("parseArgs() = nil, want error")
			}
		})
	}
}
