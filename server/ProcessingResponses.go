package server

import (
	"sync"
	"net"
	"project/protocols"
	"encoding/json"
)

func WriteClientResponse(responses <-chan Response, wg3 sync.WaitGroup, clientlist map[net.Conn]string) {
	defer wg3.Done()
	var m protocols.MessageProtocol
	for response := range responses {
		m = protocols.MessageProtocol{"sendMessage", string(response.Message)}
		b, _ := json.Marshal(m)
		if response.sendToSender {
			_, err := response.Conn.Write([]byte(b))
			if err != nil {
				PlotError(err)
				continue
			}
			continue
		}
		for con := range clientlist {
			if response.Conn == con {
				continue
			}
			_, err := con.Write([]byte(b))
			if err != nil {
				PlotError(err)
				continue
			}
		}
	}
}
