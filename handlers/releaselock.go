package handlers

import (
    "net"
)

type ReleaseLock struct {
    locks *map[string]bool
}

func NewReleaseLock(locks *map[string]bool) *ReleaseLock {
    return &ReleaseLock { locks: locks }
}

func (h *ReleaseLock) RequestToken() string {
    return "RELEASE_LOCK"
}

func (h *ReleaseLock) Handle(request string, words []string, client *net.TCPConn) StatusCode {
    filename := words[1]

    var response string

    if fileLock, ok := (*h.locks)[filename]; ok {
        if fileLock {
            response = "RELEASED_LOCK " + filename
        } else {
            response = "RELEASED_LOCK " + filename
        }
    } else {
        response = "RELEASED_LOCK " + filename
    }

    client.Write([]byte(response))

    return STATUS_OK
}
