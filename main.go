package main

import (
    "log"
    "net"
    "os"
    "requesthandler"
    "tcpserver"
)

func main() {
    args := os.Args[1:]

    if len(args) < 1 {
        log.Fatal("Must specify port number to listen on")
    }

    port := args[0]

    server := tcpserver.New("127.0.0.1", port, 10)

    addrs, err := net.LookupHost("127.0.0.1")
    if err == nil {
        server.AddHandler(requesthandler.NewHelo(addrs[0], port))
    }

    server.Start()
}
