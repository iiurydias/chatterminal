package server

import (
	"sync"
	"net"
	"strings"
	"fmt"
	"encoding/json"
	"project/protocols"
)

func ProcessClientRequest(requests chan Request, responses chan Response, wg3 sync.WaitGroup, clientlist map[net.Conn]string) {
	defer wg3.Done()
	for request := range requests {
		if string(request.Action) == "getlist" {
			var list string
			for _, client := range clientlist {
				list = list + client + "\n"
			}
			responses <- Response{[]byte(list), request.Conn, true}
			fmt.Println("ACTION MADE")
			continue
		}
		if string(request.Action) == "changename" {
			if string(request.Message) == "empty" || string(request.Message) ==  "" {
				responses <- Response{[]byte("THIS ACTION NEEDS A COMPLEMENT"), request.Conn, true}
				continue
			}
			if val, ok := clientlist[request.Conn]; ok {
				oldname := val
				clientlist[request.Conn] = string(request.Message)
				fmt.Println(strings.ToUpper(oldname) + " CHANGED NAME TO " + strings.ToUpper(clientlist[request.Conn]))
				responses <- Response{[]byte(strings.ToUpper(oldname) + " CHANGED NAME TO " + strings.ToUpper(clientlist[request.Conn])), request.Conn, false}
				continue
			}
			clientlist[request.Conn] = string(request.Message)
			fmt.Println(strings.ToUpper(clientlist[request.Conn]) + " CONNECTED")
			continue
		}
		if string(request.Action) == "logout" {
			responses <- Response{[]byte(strings.ToUpper(clientlist[request.Conn]) + " IS OFFLINE"), request.Conn, false}
			continue
		}
		if string(request.Action) == "sendMessage" {
			responses <- Response{[]byte(clientlist[request.Conn] + ": " + string(request.Message)), request.Conn, false}
			fmt.Println("MESSAGE SENT")
			continue
		}
		responses <- Response{[]byte("ACTION DOES NOT EXIST"), request.Conn, true}
	}
}
func DealWithRequests(conn net.Conn, wg3 sync.WaitGroup, request chan Request, remChan chan net.Conn) {
	defer wg3.Done()
	var wg4 sync.WaitGroup
	wg4.Add(1)
	go ReadClientRequest(conn, request, remChan, wg4)
	wg4.Wait()
}

func ReadClientRequest(conn net.Conn, request chan<- Request, remChan chan net.Conn, wg4 sync.WaitGroup) {
	defer wg4.Done()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			PlotError(err)
			remChan <- conn
			break
		}
		var m protocols.MessageProtocol
		json.Unmarshal(buffer[:n], &m)
		request <- Request{Message: []byte(m.Content), Conn: conn, Action: m.Action}
	}
}
func RemoveClientFromList(remChan chan net.Conn, wg3 sync.WaitGroup, clientlist map[net.Conn]string) {
	defer wg3.Done()
	for connChan := range remChan {
		for connList := range clientlist {
			if connList == connChan {
				delete(clientlist, connList)
				connList.Close()
				break
			}
		}
		continue
	}
}
