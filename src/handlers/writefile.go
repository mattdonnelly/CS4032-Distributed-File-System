package handlers

import "net"

type WriteFile struct { }

func NewWriteFile(ip string, port string) *WriteFile {
    return &WriteFile {}
}

func (h *WriteFile) RequestToken() string {
    return "WRITE_FILE"
}

func (h *WriteFile) Handle(request string, words []string, client *net.TCPConn) <-chan StatusCode {
    ch := make(chan StatusCode, 1)

    ch <- STATUS_FINISHED

    return ch
}
