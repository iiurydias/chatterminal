package client

import (
	"net"
	"sync"
	"os"
	"encoding/json"
	"bufio"
	"strings"
	"project/protocols"
)

func ProcessRequests(conn net.Conn, wg2 sync.WaitGroup) {
	defer wg2.Done()
	var response = make(chan string)
	var request = make(chan string)
	var wg3 sync.WaitGroup
	wg3.Add(4)
	go PrintResponse(response, wg3)
	go ReadServerResponse(conn, response, wg3)
	go WriteServerRequest(conn, request, wg3)
	go WriteRequest(request, wg3)
	wg3.Wait()
}

func WriteServerRequest(conn net.Conn, requests <-chan string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for request := range requests {
		_, err := conn.Write([]byte(request))
		if string(request) == `{"Action":"logout","Content":""}` {
			os.Exit(0)
		}
		if err != nil {
			PlotError(err)
			return
		}
	}
}
func WriteRequest(request chan<- string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for {
		reader := bufio.NewReader(os.Stdin)
		bText, _ := reader.ReadString('\n')
		bText = strings.Replace(bText, "\n", "", -1)
		action := strings.Split(string(bText), " ")
		var m protocols.MessageProtocol
		if string(action[0][0]) == "@" {
			switch string(action[0][1:]) {
			case "changename":
				if len(action)<2{
					action = append(action, "empty")
				}
				m = protocols.MessageProtocol{string(action[0][1:]), string(action[1])}
				break
			default:
				m = protocols.MessageProtocol{string(action[0][1:]), ""}
				break
			}
			b, _ := json.Marshal(m)
			request <- string(b)
			continue
		}
		m = protocols.MessageProtocol{"sendMessage", strings.Join(action, " ")}
		b, _ := json.Marshal(m)
		request <- string(b)
	}
}
