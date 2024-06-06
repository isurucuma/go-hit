package main

import (
	"fmt"
	"strconv"
	"strings"
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
	flagset := map[string]ParseFunc{
		"url": stringVar(&c.url),
		"n":   intVar(&c.n),
		"c":   intVar(&c.c),
		"rps": intVar(&c.rps),
	}

	for _, arg := range args {
		flgName, flgVal, ok := strings.Cut(arg, "=")
		if !ok {
			continue // when it is wrong flag format
		}
		flgName = strings.TrimPrefix(flgName, "-")
		flagParser, ok := flagset[flgName]
		if !ok {
			continue // when it is not a valid flag
		}
		err := flagParser(flgVal)
		if err != nil {
			return fmt.Errorf("error parsing the flag %s=%s : %w", flgName, flgVal, err)
		}
	}
	return nil
}
