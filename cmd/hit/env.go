package main

import (
	"errors"
	"flag"
	"strconv"
	"time"
)

type Config struct {
	url string //Server URL to send the requests
	n   int    // Number of requests to send
	c   int    // Numebr of concurrent requests to send
	rps int    // Number of requests per second to send
}

type ParseFunc func(string) error

func intVar(p *int) ParseFunc {
	return func(s string) error {
		var err error
		*p, err = strconv.Atoi(s)
		return err
	}
}

func stringVar(v *string) ParseFunc {
	return func(s string) error {
		*v = s
		return nil
	}
}

func boolVar(b *bool) ParseFunc {
	return func(s string) error {
		var err error
		*b, err = strconv.ParseBool(s)
		return err
	}
}

func durationVar(d *time.Duration) ParseFunc {
	return func(s string) error {
		var err error
		*d, err = time.ParseDuration(s)
		return err
	}
}

func ParseArgs(c *Config, args []string) error {
	// flagset := map[string]ParseFunc{
	// 	"url": stringVar(&c.url),
	// 	"n":   intVar(&c.n),
	// 	"c":   intVar(&c.c),
	// 	"rps": intVar(&c.rps),
	// }

	// for _, arg := range args {
	// 	flgName, flgVal, ok := strings.Cut(arg, "=")
	// 	if !ok {
	// 		continue // when it is wrong flag format
	// 	}
	// 	flgName = strings.TrimPrefix(flgName, "-")
	// 	flagParser, ok := flagset[flgName]
	// 	if !ok {
	// 		continue // when it is not a valid flag
	// 	}
	// 	err := flagParser(flgVal)
	// 	if err != nil {
	// 		return fmt.Errorf("error parsing the flag %s=%s : %w", flgName, flgVal, err)
	// 	}
	// }
	// return nil

	fs := flag.NewFlagSet("hit", flag.ContinueOnError)

	fs.StringVar(&c.url, "url", "", "Server URL to send the requests")
	fs.Var(newPositiveIntValue(&c.n), "n", "Number of requests to send")
	fs.Var(newPositiveIntValue(&c.c), "c", "Number of concurrent requests to send")
	fs.Var(newPositiveIntValue(&c.rps), "rps", "Number of requests per second to send")

	return fs.Parse(args)
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
