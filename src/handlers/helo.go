package handlers

import "net"

type Helo struct {
    ip string
    port string
}

func NewHelo(ip string, port string) *Helo {
    return &Helo {
        ip: ip,
        port: port,
    }
}

func (h *Helo) RequestToken() string {
    return "HELO"
}

func (h *Helo) Handle(request string, words []string, client *net.TCPConn) <-chan StatusCode {
    ch := make(chan StatusCode, 1)

    response := request + "\n" +
                "IP:" + h.ip + "\n" +
                "Port:" + h.port + "\n" +
                "StudentID:11350561"

    client.Write([]byte(response))

    ch <- STATUS_FINISHED

    return ch
}
