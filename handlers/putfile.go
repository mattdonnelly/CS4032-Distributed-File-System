package handlers

import (
    "log"
    "net"
)

type PutFile struct {
    fileLocations *map[string]string
}

func NewPutFile(locations *map[string]string) *PutFile {
    return &PutFile { fileLocations: locations }
}

func (h *PutFile) RequestToken() string {
    return "PUT_FILE"
}

func (h *PutFile) Handle(request string, words []string, client *net.TCPConn) StatusCode {
    servAddr := words[1]
    filename := words[2]
    startOfData := len(words[0]) + 1 + len(words[1]) + 1 + len(words[2]) + 1
    data := request[startOfData:]

    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        log.Println("ResolveTCPAddr failed:", err.Error())
        return STATUS_ERROR
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        log.Println("Dial failed:", err.Error())
        return STATUS_ERROR
    }

    new_request := "WRITE_FILE " + filename + " " + data

    _, err = conn.Write([]byte(new_request))
    if err != nil {
        log.Println("Write to server failed:", err.Error())
        return STATUS_ERROR
    }

    (*h.fileLocations)[filename] = servAddr

    response := "PUT_FILE_OK " + servAddr + " " + filename

    client.Write([]byte(response))

    return STATUS_OK
}
