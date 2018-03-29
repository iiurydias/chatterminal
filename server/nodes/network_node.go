package nodes

import (
	"net"
	"strconv"
	"log"
	"sync"
)

type NetworkNode struct {
	Conn   *net.TCPConn
	Status int
	Name   string
}

func NewNetworkNode(conn *net.TCPConn) *NetworkNode {
	return &NetworkNode{conn, LOGOUT, ""}
}

func (n *NetworkNode) Write(response *Response, error chan Error, wg sync.WaitGroup) {
	defer wg.Done()
	if response.Node == n && response.SendToSender || response.Node != n && !response.SendToSender {
		msg := "{\"Action\":\"sendMessage\",\"Content\":\"" + response.Msg +"\"}"
		_, err := n.Conn.Write([]byte(msg))
		if err!= nil {
			error <- Error{n, err}
			return
		}
	}
}

func (n *NetworkNode) Read(request chan<- *Request,  error chan Error, wg3 sync.WaitGroup) {
	go func() {
		defer wg3.Done()
		for {
			buffer := make([]byte, 1024)
			num, err := n.Conn.Read(buffer)
			if err!= nil {
				error <- Error{n, err}
				return
			}
			message := string(buffer[:num])
			log.Println(n.GetInfo() + " sent a request: " +  message)
			request <- &Request{message, n}
		}
	}()
}

func (n *NetworkNode) SetName(name string) {
	n.Name = name
}

func (n *NetworkNode) SetStatus(status int) {
	n.Status = status
}

func (n *NetworkNode) GetName() string {
	return n.Name
}

func (n *NetworkNode) GetStatus() int {
	return n.Status
}

func (n *NetworkNode) GetInfo() string {
	return "Network Node (" + n.GetName() + "): " + n.Conn.RemoteAddr().String() + " Status: " + strconv.Itoa(n.GetStatus())
}