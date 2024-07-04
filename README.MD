## Go-Hit tool

### Description

Go-Hit is a lightweight load testing tool developed using Golang. It allows you to send HTTP requests to web servers, currently supporting only GET requests. The tool provides the ability to send multiple requests per second in parallel, and you can define the concurrency level and the threshold for the requests per second (RPS) value.

### Motivation

This tool was created while following the book <B>"GO by Example - by Inanc Gumus"</B>

### Features

- Supports sending GET HTTP requests to web servers
- Ability to send multiple requests per second in parallel
- Configurable concurrency level
- Threshold setting for the requests per second (RPS) value

### Installation

To install Go-Hit, you can use the following command:

```
go get github.com/isurucuma/go-hit
```

### Usage

#### Using the Go API

To use the Go-Hit tool, follow these steps:

1. Import the package in your Go code:

```go
import "github.com/isurucuma/go-hit"
```

2. Create a new instance of the `Client` struct and configure it with the desired settings:

```go
client := go_hit.NewClient()
client.C = 10 // concurrency level
client.RPS = 50 // requests per second
client.Timeout = 10 * time.Second // timeout per request
```

3. Send the HTTP requests using the `SendN` function:

```go
result, err := client.SendN(context.Background(), "http://www.some-web-page.com", 1000)
if err != nil {
    // handle error
}
```

The `result` variable will contain the response times and other relevant information for the sent requests.

#### Using the CLI

If you prefer to use the command-line interface, you can run the `gohit` binary with the following arguments:

```
./gohit -c=10 -rps=15 -H="Content-Type: application/json" -H="Authorization: Bearer YOUR_ACCESS_TOKEN" http://www.example.com
```

Here, the arguments are:

- `-c=10`: The concurrency level (represents the simulated number of parallel users).
- `-n=1000`: The total number of requests to send.
- `-rps=50`: The number of requests per second each simulated user is eligible to send.
- `-H="Content-Type: application/json"`: The headers to include in the HTTP requests (possible to add multiple headers).
- `"http://www.some-web-page.com"`: The URL to send the GET requests to.

The end result will contain the following information:

```go
type Result struct {
    RPS      float64 // supported RPS value from the server
    Requests int    // total requests sent
    Errors   int    // total number of errors
    Bytes    int64  // total bytes received
    Duration time.Duration // duration for the whole process
    Fastest  time.Duration
    Slowest  time.Duration
    Status   int    // overall status
    Error    error  // if there are errors wrapped
}
```

## Pipeline

![Pipeline architecture](https://github.com/isurucuma/go-hit/blob/main/docs/pipeline.png)

### Contributing

If you would like to contribute to the Go-Hit project, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Implement your changes
4. Test your changes
5. Submit a pull request

### License

Go-Hit is licensed under the [MIT License](LICENSE).
