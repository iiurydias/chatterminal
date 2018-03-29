package servers

import (
	"chat-terminal/servers/nodes"
	"chat-terminal/protocols"
	"strings"
	"encoding/json"
	"strconv"
	"sync"
	"fmt"
)

type Server struct {
	Gatekeeper *nodes.Gatekeeper
}

func NewServer(gatekeeper *nodes.Gatekeeper) *Server {
	return &Server{gatekeeper}
}

func (s *Server) Process(requests chan *nodes.Request, wg sync.WaitGroup) (chan *nodes.Response, chan nodes.Error) {
	var m protocols.MessageProtocol
	response := make(chan *nodes.Response)
	errChan := make(chan nodes.Error)
	go func() {
		defer wg.Done()
		for request := range requests {
			err := json.Unmarshal([]byte(request.Msg), &m)
			fmt.Println(request.Msg)
			if err != nil{
				errChan <- nodes.Error{request.Node, err}
				continue
			}
			if request.Node.GetStatus() == nodes.LOGOUT && request.Node.GetName() != ""{
				response <- &nodes.Response{"YOU NEED TO RECONNECT", request.Node, true}
				continue
			}
			if m.Action == "getlist" {
				var list string
				count := 1
				for _, node := range s.Gatekeeper.GetLoggedNodes() {
					list = list + strconv.Itoa(count) + ". " + node.GetName() + "       "
					count++
				}
				response <- &nodes.Response{list, request.Node, true}
				continue
			}
			if m.Action == "changename" {
				if m.Content == ""  || m.Content == "empty"{
					response <- &nodes.Response{"THIS ACTION NEEDS A COMPLEMENT", request.Node, true}
					continue
				}
				if request.Node.GetStatus() != nodes.LOGOUT {
					oldname := request.Node.GetName()
					request.Node.SetName(m.Content)
					response <- &nodes.Response{strings.ToUpper(oldname) + " CHANGED NAME TO " + request.Node.GetName(), request.Node, false}
					continue
				}
				request.Node.SetName(m.Content)
				request.Node.SetStatus(nodes.LOGIN)
				response <- &nodes.Response{request.Node.GetName() + " IS ONLINE", request.Node, false}
				continue
			}
			if m.Action == "logout" {
				response <- &nodes.Response{request.Node.GetName() + " IS OFFLINE", request.Node, false}
				request.Node.SetStatus(nodes.LOGOUT)
				continue
			}
			if m.Action == "sendMessage" {
				response <- &nodes.Response{request.Node.GetName() + ": " + m.Content, request.Node, false}
				continue
			}
			response <- &nodes.Response{"ACTION DOES NOT EXIST", request.Node, true}
		}
	}()
	return response, errChan
}

//
//func (servers *Server) HandleNewClients() {
//	connChan := make(chan net.Conn)
//	var wg2 sync.WaitGroup
//	wg2.Add(2)
//	go func() {
//		defer wg2.Done()
//		servers.gatekeeper.AcceptClient(connChan, servers.gatekeeper.listener)
//	}()
//	go func() {
//		defer wg2.Done()
//		servers.RouteClient(connChan)
//	}()
//	wg2.Wait()
//}
//func (servers *Server) RouteClient(connChan chan net.Conn) {
//	var wg3 sync.WaitGroup
//	var clientlist = make(map[net.Conn]string)
//	var response = make(chan protocols.Response)
//	var request = make(chan protocols.Request)
//	var remChan = make(chan net.Conn)
//	wg3.Add(3)
//	go func() {
//		var wg7 sync.WaitGroup
//		defer wg3.Done()
//		for responses := range response {
//			for _, node := range servers.gatekeeper.nodeList {
//				wg7.Add(1)
//				go func() {
//					defer wg7.Done()
//					node.Read(responses, &clientlist)
//				}()
//			}
//		}
//		wg7.Wait()
//	}()
//	go func() {
//		defer wg3.Done()
//		Process(request, response, clientlist)
//	}()
//	go func() {
//		defer wg3.Done()
//		RemoveClientFromList(remChan, clientlist)
//	}()
//	for conn := range connChan {
//		wg3.Add(1)
//		go func() {
//			defer wg3.Done()
//			servers.DealWithRequests(conn, request, remChan)
//		}()
//	}
//	wg3.Wait()
//}
//func (servers *Server) DealWithRequests(conn net.Conn, request chan protocols.Request, remChan chan net.Conn) {
//	var wg4 sync.WaitGroup
//	for _, node := range servers.gatekeeper.nodeList {
//		wg4.Add(1)
//		go func() {
//			defer wg4.Done()
//			node.Read(conn, request, remChan)
//		}()
//	}
//	wg4.Wait()
//}
//func RemoveClientFromList(remChan chan net.Conn, clientlist map[net.Conn]string) {
//	for connChan := range remChan {
//		for connList := range clientlist {
//			if connList == connChan {
//				delete(clientlist, connList)
//				connList.Close()
//				break
//			}
//		}
//		continue
//	}
//}
