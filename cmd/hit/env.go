package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	url     string            // Server URL to send the requests
	n       int               // Number of requests to send
	c       int               // Number of concurrent requests to send
	rps     int               // Number of requests per second to send
	headers map[string]string // HTTP headers to include in the requests
}

type env struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
	dry    bool
}

type positiveIntValue int

type headersFlag struct {
	headers map[string]string
}

func ParseArgs(c *Config, args []string, stderr io.Writer) error {
	fs := flag.NewFlagSet("hit", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage: %s [options] url\n", fs.Name())
		fs.PrintDefaults()
	}

	fs.Var(newPositiveIntValue(&c.n), "n", "Number of requests to send")
	fs.Var(newPositiveIntValue(&c.c), "c", "Number of concurrent requests to send")
	fs.Var(newPositiveIntValue(&c.rps), "rps", "Number of requests per second to send")

	headers := &headersFlag{headers: make(map[string]string)}
	fs.Var(headers, "H", "HTTP headers to include in the requests")

	err := fs.Parse(args)
	if err != nil {
		return err
	}
	c.url = fs.Arg(0)
	c.headers = headers.headers
	if err := validateArgs(c); err != nil {
		fmt.Fprintln(fs.Output(), err)
		fs.Usage()
		return err
	}
	return nil
}

func newPositiveIntValue(p *int) *positiveIntValue {
	return (*positiveIntValue)(p)
}

func (n *positiveIntValue) String() string {
	return strconv.Itoa(int(*n))
}

func (n *positiveIntValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}

	if v <= 0 {
		return errors.New("should be a positive integer")
	}

	*n = positiveIntValue(v)
	return nil
}

func (h *headersFlag) String() string {
	var result string
	for k, v := range h.headers {
		result += fmt.Sprintf("%s: %s\n", k, v)
	}
	return result
}

func (h *headersFlag) Set(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid header format, should be 'Key: Value'")
	}
	key := strings.TrimSpace(parts[0])
	key = strings.TrimPrefix(key, "\"")
	value := strings.TrimSpace(parts[1])
	value = strings.TrimSuffix(value, "\"")
	h.headers[key] = value
	return nil
}

func validateArgs(c *Config) error {
	const urlArg = "url argument"

	u, err := url.Parse(c.url)
	if err != nil {
		return argError(c.url, urlArg, err)
	}
	if c.url == "" || u.Host == "" || u.Scheme == "" {
		return argError(c.url, urlArg, errors.New("requires a valid url"))
	}
	if c.n < c.c {
		err := fmt.Errorf(`should be greater than -c: "%d"`, c.c)
		return argError(c.n, "flag -n", err)
	}

	return nil
}

func argError(value any, arg string, err error) error {
	return fmt.Errorf(`invalid value "%v" for %s: %w`, value, arg, err)
}
