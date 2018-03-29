package client

import (
	"net"
	"sync"
	"log"
	"chat-terminal/client/nodes_client"

)

type Client struct {
	nodeA nodes_client.Node
	nodeB nodes_client.Node
}

func NewClient(nodeA, nodeB nodes_client.Node) *Client {
	return &Client{nodeA, nodeB}
}


func (client *Client) HandleConnection() {
	_, err := client.nodeA.GetConn().Write([]byte(`{"Action":"changename","Content":"` + client.nodeA.GetName() + `"}`))
	if err != nil {
		log.Fatalln(err.Error())
	}
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		client.ProcessRequests(client.nodeA.GetConn())
	}()
	wg2.Wait()
}
func (client *Client) ProcessRequests(conn net.Conn) {
	var response = make(chan string)
	var request = make(chan string)
	var wg3 sync.WaitGroup
	wg3.Add(4)
	go func() {
		defer wg3.Done()
		client.nodeB.Write(conn,response)
	}()
	go func() {
		defer wg3.Done()
		client.nodeA.Read(conn, response)
	}()
	go func() {
		defer wg3.Done()
		client.nodeA.Write(conn, request)
	}()
	go func() {
		defer wg3.Done()
		client.nodeB.Read(conn, request)
	}()
	wg3.Wait()
}
