package tcpserver

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

type TCPServer struct {
    addr *net.TCPAddr
    requestHandlers map[string]RequestHandler
    threadCount int
    clientChan chan *net.TCPConn
    killChan chan bool
}

func New(host string, port string, threadCount int) *TCPServer {
    addr, _ := net.ResolveTCPAddr("tcp", host + ":" + port)
    return &TCPServer{
        addr: addr,
        requestHandlers: make(map[string]Protocol),
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
    for {
        client := <- server.clientChan

        reader := bufio.NewReader(client)
        responseChan := make(chan byte)
        requestChan := make(chan byte)

        go func(){
            for responseByte := range responseChan {
                client.Write([]byte{ responseByte })
            }
        }()

        go func(){
            for {
                b, err := reader.ReadByte()
                if err != nil {
                    break
                }

                requestChan <- b
            }
        }()

        status := 0

        for status != 1 {
            buffer := make([]byte,0)

            for {
                b := <- requestChan

                if b != '\n' && b != ' ' && b != ':' && b != '\r' {
                    buffer = append(buffer, b)
                }
                else {
                    break
                }
            }

            token := string(buffer)

            if token != "" {
                if token == "KILL_SERVICE" {
                    server.killChan <- true
                    log.Println("Killing service")
                }
                else {
                    status = 1
                }
            } else {
                log.Println("Encountered empty token")
            }
        }
    }
}
