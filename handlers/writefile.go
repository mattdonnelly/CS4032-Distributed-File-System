package handlers

import (
    "io/ioutil"
    "log"
    "net"
    "os"
)

type WriteFile struct { }

func NewWriteFile() *WriteFile {
    return &WriteFile {}
}

func (h *WriteFile) RequestToken() string {
    return "WRITE_FILE"
}

func (h *WriteFile) Handle(request string, words []string, client *net.TCPConn) StatusCode {
    if _, err := os.Stat("./dfs-files/"); err != nil {
        if dirErr := os.Mkdir("./dfs-files/", 0755); dirErr != nil {
            log.Println("Failed create dfs directory: ", err)
            return STATUS_ERROR
        }
    }

    filename := "./dfs-files/" + words[1]
    startOfData := len(words[0]) + 1 + len(words[1]) + 1
    data := request[startOfData:]

    if err := ioutil.WriteFile(filename, []byte(data), 0755); err != nil {
        log.Println("Failed to open file: " + filename + " Error: " + err.Error())
        return STATUS_ERROR
    } else {
        return STATUS_OK
    }
}
