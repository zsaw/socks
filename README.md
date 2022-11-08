# Socks Library

Socks is a library written in Go (Golang).

Socks is implemented using SOCKS5 algorithm.

## Contents

- [Socks Library](#socks-library)
  - [Installation](#installation)
  - [Quick start](#quick-start)
    - [Server](#server)
    - [Client](#client)

## Installation

To install socks package, you need to install Go and set your Go workspace first.

1. You first need [Go](https://golang.org/) installed (**version 1.7+ is required**), then you can use the below Go command to install socks.

    ```sh
    go get github.com/zsaw/socks
    ```

2. Import it in your code:

    ```go
    import "github.com/zsaw/socks"
    ```

## Quick start

### Server
```go
package main

import (
    "log"

    "github.com/zsaw/socks"
)

func main() {
    err := socks.ListenAndServe(":1080")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
```

### Client
```go
package main

import (
    "bufio"
	"log"
	"net"
	"net/http"
	"os"

    "github.com/zsaw/socks"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:1080")
    if err != nil {
        log.Println(err.Error())
        return
    }
    defer conn.Close()

    conn, err = socks.Client(conn, "93.184.216.34:80")
    if err != nil {
        log.Println(err.Error())
        return
    }

    req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
    err = req.Write(conn)
    if err != nil {
        log.Println(err.Error())
        return
    }

    resp, err := http.ReadResponse(bufio.NewReader(conn), req)
    if err != nil {
        log.Println(err.Error())
        return
    }

    resp.Write(os.Stdout)
}
```
