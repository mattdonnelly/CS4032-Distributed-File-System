package handlers

import (
    "io/ioutil"
    "log"
    "net"
)

type ReadFile struct { }

func NewReadFile() *ReadFile {
    return &ReadFile {}
}

func (h *ReadFile) RequestToken() string {
    return "WRITE_FILE"
}

func (h *ReadFile) Handle(request string, words []string, client *net.TCPConn) StatusCode {
    filename := "./dfs-files/" + words[1]
    data, err := ioutil.ReadFile(filename)

    if err != nil {
        log.Println("Failed to open file: " + filename + " Error: " + err.Error())
        return STATUS_ERROR
    }

    response := "FILE: " + filename + "\n" +
                "DATA: " + string(data)

    client.Write([]byte(response))

    return STATUS_OK
}
