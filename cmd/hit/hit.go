package main

import (
	"fmt"
	"os"
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
	if err := ParseArgs(&config, os.Args[1:], e.stderr); err != nil {
		return err
	}
	fmt.Fprintf(e.stderr, "%s\n\nSending %d requests to %q (concurrency: %d)\n", logo, config.n, config.url, config.c)

	if e.dry {
		return nil
	}
	return runHit(e, &config)
}

func runHit(e *env, config *Config) error {
	return nil
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
