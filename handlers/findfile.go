package handlers

import (
    "net"
)

type FindFile struct {
    fileLocations *map[string]string
}

func NewFindFile(locations *map[string]string) *FindFile {
    return &FindFile { fileLocations: locations }
}

func (h *FindFile) RequestToken() string {
    return "FIND"
}

func (h *FindFile) Handle(request string, words []string, client *net.TCPConn) StatusCode {
    var response string

    if fileserver, ok := (*h.fileLocations)[words[1]]; ok {
        response = "FOUND_FILE: " + fileserver
    } else {
        response = "ERROR: Could not locate file"
    }

    client.Write([]byte(response))

    return STATUS_OK
}
