package server

import (
	"sync"
	"net"
	"log"
	"project/settings"
)

func HandleNewClients(wg sync.WaitGroup, s settings.Settings) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.Network, s.Address)
	if err != nil {
		log.Fatalln(err.Error())
	}
	listener, err := net.ListenTCP(s.Network, addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	connChan := make(chan net.Conn)
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go AcceptClient(connChan, listener, wg2)
	go RouteClient(connChan, wg2)
	wg2.Wait()
}
func AcceptClient(conexao chan net.Conn, listener net.Listener, wg2 sync.WaitGroup) {
	defer wg2.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			PlotError(err)
			continue
		}
		conexao <- conn
	}
}
func RouteClient(connChan chan net.Conn, wg2 sync.WaitGroup) {
	defer wg2.Done()
	var wg3 sync.WaitGroup
	var clientlist = make(map[net.Conn]string)
	var response = make(chan Response)
	var request = make(chan Request)
	var remChan = make(chan net.Conn)
	wg3.Add(3)
	go WriteClientResponse(response, wg3, clientlist)
	go ProcessClientRequest(request, response, wg3, clientlist)
	go RemoveClientFromList(remChan, wg3, clientlist)
	for conn := range connChan {
		wg3.Add(1)
		go DealWithRequests(conn, wg3, request, remChan)
	}
	wg3.Wait()
}
