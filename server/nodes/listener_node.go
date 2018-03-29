package nodes

import (
	"net"
	"log"
	"chat-terminal/settings"
	"strconv"
	"sync"
)

type ListenerNode struct{
	gatekeeper *Gatekeeper
	listener   *net.TCPListener
}

func NewListenerNode(gatekeeper *Gatekeeper, settings *settings.Settings) *ListenerNode{
	addr, err := net.ResolveTCPAddr(settings.Network, settings.Address)
	if err != nil {
		log.Fatalln(err.Error())
	}
	listener, err := net.ListenTCP(settings.Network, addr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &ListenerNode{gatekeeper, listener}
}

func (l *ListenerNode) Read(request chan <- *Request, error chan Error, wg2 sync.WaitGroup){
	var wg3 sync.WaitGroup
	go func() {
		defer wg2.Done()
		for {
			conn, err := l.listener.AcceptTCP()
			if err != nil {
				error <- Error{l, err}
				continue
			}
			wg3.Add(1)
			node := NewNetworkNode(conn)
			node.Read(request, error, wg3)
			l.gatekeeper.AddNode(node)
		}
	}()
}
func (l *ListenerNode) Write(response *Response, error chan Error, wg sync.WaitGroup){
	return
}

func (l *ListenerNode) SetName(name string) {
	return
}
func (l *ListenerNode) SetStatus(status int) {
	return
}

func (l *ListenerNode) GetName() string {
	return "Listener Node"
}
func (l *ListenerNode) GetStatus() int {
	return SYSTEM
}

func (l *ListenerNode) GetInfo() string {
	return l.GetName() + ": " + l.listener.Addr().String() + " Status: " + strconv.Itoa(l.GetStatus())
}