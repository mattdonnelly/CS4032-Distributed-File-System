package requesthandler

import "net"

type StatusCode int

const (
    STATUS_ERROR StatusCode = iota
    STATUS_UNDEFINED
    STATUS_FINISHED
)

type RequestHandler interface {
    RequestToken() string
    Handle(request string, words []string, client *net.TCPConn) <-chan StatusCode
}
