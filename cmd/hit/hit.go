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

func main() {
	var config Config = Config{
		n: 100,
		c: 1,
	}
	if err := ParseArgs(&config, os.Args[1:]); err != nil {
		os.Exit(1)
	}
	fmt.Printf(
		"%s\n\nSending %d requests to %q (concurrency: %d)\n",
		logo, config.n, config.url, config.c,
	)
}
