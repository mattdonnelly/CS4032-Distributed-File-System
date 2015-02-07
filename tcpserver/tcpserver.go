package tcpserver

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    "github.com/mattdonnelly/CS4032-Distributed-File-System/handlers"
)

type TCPServer struct {
    addr *net.TCPAddr
    requestHandlers map[string]handlers.RequestHandler
    threadCount int
    clientChan chan *net.TCPConn
    killChan chan bool
}

func New(host string, port string, threadCount int) *TCPServer {
    addr, _ := net.ResolveTCPAddr("tcp", host + ":" + port)
    return &TCPServer{
        addr: addr,
        requestHandlers: make(map[string]handlers.RequestHandler),
        threadCount: threadCount,
        clientChan: make(chan *net.TCPConn, threadCount),
        killChan: make(chan bool),
    }
}

func (server *TCPServer) Start() {
    listener, err := net.ListenTCP("tcp", server.addr)
    if err != nil {
        log.Fatal("Couldn't start server: " + err.Error())
    }

    fmt.Println("Listening on " + server.addr.String() + "...")

    server.acceptConnections(listener)
}

func (server *TCPServer) AddHandler(handler handlers.RequestHandler) {
    if _, full := server.requestHandlers[handler.RequestToken()]; full {
        log.Println("Protocol already exists: " + handler.RequestToken())
    }

    server.requestHandlers[handler.RequestToken()] = handler

    _ , success := server.requestHandlers[handler.RequestToken()]
    if !success {
        log.Fatal("Failed to add request handler: " + handler.RequestToken())
    }
}

func (server *TCPServer) RouteRequest(request string, client *net.TCPConn) <-chan handlers.StatusCode {
    words := strings.Fields(request)

    handler, success := server.requestHandlers[words[0]]
    if success {
        return handler.Handle(request, words, client)
    } else {
        log.Println("UKNOWN_REQUEST: " + request)
        statusChan := make(chan handlers.StatusCode, 1)
        statusChan <- handlers.STATUS_ERROR
        return statusChan
    }
}

func (server *TCPServer) acceptConnections(listener *net.TCPListener) {
    acceptChan := make(chan *net.TCPConn)

    go func() {
        for {
            client, err := listener.AcceptTCP()

            if err != nil {
                log.Printf("Couldn't accept client: " + err.Error())
                continue
            }

            log.Printf("Accepted Connection: %v <-> %v\n", client.LocalAddr(), client.RemoteAddr())
            acceptChan <- client
        }
    }()

    for i := 0; i < server.threadCount; i++ {
        go server.handleConnection()
    }

    for {
        select {
            case client := <- acceptChan:
                select {
                    case server.clientChan <- client:
                    default:
                        log.Println("Could not accept connection")
                }
            case <- server.killChan:
                log.Println("Killing Service...")
                return
        }
    }
}

func (server *TCPServer) handleConnection() {
    status := handlers.STATUS_UNDEFINED

    for status != handlers.STATUS_DISCONNECT {
        client := <- server.clientChan

        reader := bufio.NewReader(client)

        buf := []byte{}
        readErr := false
        for {
            line, err := reader.ReadBytes('\n')
            if err != nil {
                readErr = true
                break
            }

            buf = append(buf, line...)

            if peek, err := reader.Peek(1); err == nil && string(peek) == "\n" {
                break
            }
        }

        if !readErr {
            request :=  strings.TrimSpace(string(buf))
            log.Println(request)
            status = <- server.RouteRequest(request, client)
        } else {
            break
        }
    }
}
