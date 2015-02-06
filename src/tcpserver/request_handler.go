package tcpserver

type StatusCode int

type RequestHandler interface {
    Token() string
    Handle(requestChan <-chan byte, responseChan chan<- byte) <-chan StatusCode
}
