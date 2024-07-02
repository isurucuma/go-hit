package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

type Config struct {
	url string //Server URL to send the requests
	n   int    // Number of requests to send
	c   int    // Numebr of concurrent requests to send
	rps int    // Number of requests per second to send
}

type env struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
	dry    bool
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

	err := fs.Parse(args)
	if err != nil {
		return err
	}
	c.url = fs.Arg(0)
	if err := validateArgs(c); err != nil {
		fmt.Fprintln(fs.Output(), err)
		fs.Usage()
		return err
	}
	return nil
}

type positiveIntValue int

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
