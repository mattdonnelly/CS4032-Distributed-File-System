package main

import (
    "log"
    "os"
    "tcpserver"
)

func main() {
    args := os.Args[1:]

    if len(args) < 1 {
        log.Fatal("Must specify port number to listen on")
    }

    port := args[0]

    server := tcpserver.New("127.0.0.1", port, 10)
    server.Start()
}
