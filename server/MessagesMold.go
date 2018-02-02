package server

import "net"

type Request struct {
	Message []byte
	Conn    net.Conn
	Action  string
}
type Response struct {
	Message      []byte
	Conn         net.Conn
	sendToSender bool
}

