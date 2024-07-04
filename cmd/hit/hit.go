package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/isurucuma/go-hit"
)

const logo = `
 ___  ___  ___  _________       
|\  \|\  \|\  \|\___   ___\     
\ \  \\\  \ \  \|___ \  \_|     
 \ \   __  \ \  \   \ \  \      
  \ \  \ \  \ \  \   \ \  \     
   \ \__\ \__\ \__\   \ \__\    
    \|__|\|__|\|__|    \|__| `

func run(e *env) error {
	var config Config = Config{
		n: 100,
		c: 1,
	}
	if err := ParseArgs(&config, e.args[1:], e.stderr); err != nil {
		return err
	}
	fmt.Fprintf(e.stdout, "%s\n\nSending %d requests to %q (concurrency: %d)\n", logo, config.n, config.url, config.c)

	if e.dry {
		return nil
	}
	return runHit(e, &config)
}

func runHit(e *env, config *Config) error {
	handleErr := func(err error) error {
		if err != nil {
			fmt.Fprintf(e.stderr, "\nerror occurred: %v\n", err)
			return err
		}
		return nil
	}

	const (
		timeout           = 5 * time.Second
		timeoutPerRequest = 30 * time.Second
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()

	sum, err := hit.SendN(ctx,
		config.url,
		config.headers,
		config.n,
		hit.Concurrency(10),
		hit.RequestsPerSecond(20))
	sum.Fprint(e.stdout)

	if err = ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		return handleErr(fmt.Errorf("timed out in %s", timeout))
	}
	return handleErr(err)
}

func main() {
	e := &env{
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   os.Args,
	}

	if err := run(e); err != nil {
		os.Exit(1)
	}
}
