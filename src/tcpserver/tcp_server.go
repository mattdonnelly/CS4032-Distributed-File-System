package tcpserver

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "requesthandler"
    "strings"
)

type TCPServer struct {
    addr *net.TCPAddr
    requestHandlers map[string]requesthandler.RequestHandler
    threadCount int
    clientChan chan *net.TCPConn
    killChan chan bool
}

func New(host string, port string, threadCount int) *TCPServer {
    addr, _ := net.ResolveTCPAddr("tcp", host + ":" + port)
    return &TCPServer{
        addr: addr,
        requestHandlers: make(map[string]requesthandler.RequestHandler),
        threadCount: threadCount,
        clientChan: make(chan *net.TCPConn, threadCount),
        killChan: make(chan bool),
    }
}

func (server* TCPServer) Start() {
    listener, err := net.ListenTCP("tcp", server.addr)
    if err != nil {
        log.Fatal("Couldn't start server: " + err.Error())
    }

    fmt.Println("Listening on " + server.addr.String() + "...")

    server.acceptConnections(listener)
}

func (server* TCPServer) acceptConnections(listener *net.TCPListener) {
    acceptChan := make(chan *net.TCPConn)

    go func() {
        for {
            client, err := listener.AcceptTCP()

            if err != nil {
                fmt.Printf("Couldn't accept client: " + err.Error())
                continue
            }

            fmt.Printf("Accepted Connection: %v <-> %v\n", client.LocalAddr(), client.RemoteAddr())
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
                        fmt.Println("Could not accept connection")
                }
            case <- server.killChan:
                fmt.Println("Killing Service...")
                return
        }
    }
}

func (server* TCPServer) handleConnection() {
    status := requesthandler.STATUS_UNDEFINED

    for status != requesthandler.STATUS_FINISHED || status != requesthandler.STATUS_ERROR {
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
            status = <- server.RouteRequest(string(buf), client)
        } else {
            break
        }
    }
}

func (server* TCPServer) AddHandler(handler requesthandler.RequestHandler) {
    if _, full := server.requestHandlers[handler.RequestToken()]; full {
        fmt.Println("Protocol already exists: " + handler.RequestToken())
    }

    server.requestHandlers[handler.RequestToken()] = handler

    _ , success := server.requestHandlers[handler.RequestToken()]
    if !success {
        log.Fatal("Failed to add request handler: " + handler.RequestToken())
    }
}

func (server* TCPServer) RouteRequest(request string, client *net.TCPConn) <-chan requesthandler.StatusCode {
    words := strings.Fields(request)

    handler, success := server.requestHandlers[words[0]]
    if success {
        return handler.Handle(words, client)
    } else {
        statusChan := make(chan requesthandler.StatusCode, 1)
        statusChan <- requesthandler.STATUS_ERROR
        return statusChan
    }
}
