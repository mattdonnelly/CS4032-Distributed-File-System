package main

import (
    "flag"
    "log"
    "net"
    "handlers"
    "tcpserver"
)

func main() {
    port := flag.String("p", "", "Port for server to listen on")
    serverType := flag.String("t", "", "Type of server to start (either 'DS' or 'FS'")

    flag.Parse()

    if *port == "" {
        log.Fatal("Must specify port for server to listen on.")
    }

    if *serverType != "FS" && *serverType != "DS" {
        log.Fatal("Must specifY type of server to start")
    }

    server := tcpserver.New("127.0.0.1", *port, 10)

    addrs, err := net.LookupHost("127.0.0.1")
    if err == nil {
        server.AddHandler(handlers.NewHelo(addrs[0], *port))
    }

    if *serverType == "FS" {
        server.AddHandler(handlers.NewWriteFile())
    }

    server.Start()
}
