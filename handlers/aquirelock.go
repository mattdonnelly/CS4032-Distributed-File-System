package handlers

import (
    "net"
)

type AquireLock struct {
    locks *map[string]bool
}

func NewAquireLock(locks *map[string]bool) *AquireLock {
    return &AquireLock { locks: locks }
}

func (h *AquireLock) RequestToken() string {
    return "AQUIRE_LOCK"
}

func (h *AquireLock) Handle(request string, words []string, client *net.TCPConn) StatusCode {
    filename := words[1]

    var response string

    if fileLock, ok := (*h.locks)[filename]; ok {
        if !fileLock {
            (*h.locks)[filename] = true
            response = "AQUIRED_LOCK " + filename
        } else {
            response = "ERROR: File is currently in use"
        }
    } else {
        (*h.locks)[filename] = true
        response = "AQUIRED_LOCK " + filename
    }

    client.Write([]byte(response))

    return STATUS_OK
}
