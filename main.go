package main

import (
    "flag"
    "log"
    "net"
    "github.com/mattdonnelly/CS4032-Distributed-File-System/handlers"
    "github.com/mattdonnelly/CS4032-Distributed-File-System/tcpserver"
)

func main() {
    port := flag.String("p", "", "Port for server to listen on")
    serverType := flag.String("t", "", "Type of server to start (either 'DS', 'FS' or 'LS'")

    flag.Parse()

    if *port == "" {
        log.Fatal("Must specify port for server to listen on.")
    }

    if *serverType != "FS" && *serverType != "DS" && *serverType != "LS" {
        log.Fatal("Must specify type of server to start")
    }

    server := tcpserver.New("127.0.0.1", *port, 10)

    addrs, err := net.LookupHost("127.0.0.1")
    if err == nil {
        server.AddHandler(handlers.NewHelo(addrs[0], *port))
    }

    if *serverType == "FS" {
        server.AddHandler(handlers.NewWriteFile())
        server.AddHandler(handlers.NewReadFile())
    } else if *serverType == "DS" {
        fileLocations := make(map[string]string)

        server.AddHandler(handlers.NewPutFile(&fileLocations))
        server.AddHandler(handlers.NewFindFile(&fileLocations))
    } else if *serverType == "LS" {
        locks := make(map[string]bool)

        server.AddHandler(handlers.NewAquireLock(&locks))
        server.AddHandler(handlers.NewReleaseLock(&locks))
    }

    server.Start()
}
