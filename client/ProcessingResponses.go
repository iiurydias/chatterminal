package client

import (
	"sync"
	"fmt"
	"net"
	"project/protocols"
	"encoding/json"
)

func PrintResponse(responses chan string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for response := range responses {
		var m protocols.MessageProtocol
		json.Unmarshal([]byte(response), &m)
		fmt.Println(m.Content)
	}
}
func ReadServerResponse(conn net.Conn, response chan<- string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			PlotError(err)
			return
		}
		response <- string(buffer[:n])
	}
}
