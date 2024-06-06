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

const usage = `
Usage:
	- url: Server URL to hit (required)
	- n: Number of requests to send (default: 1)
	- c: Number of concurrent requests to send (default: 1)
	- rps: Number of requests per second to send (default: 1)	
`

func main() {
	var config Config
	if err := ParseArgs(&config, os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(
		"%s\n\nSending %d requests to %q (concurrency: %d)\n",
		logo, config.n, config.url, config.c,
	)
}
