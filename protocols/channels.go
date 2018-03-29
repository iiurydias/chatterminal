package protocols

import "net"

type Request struct {
	Message []byte
	Conn    net.Conn
	Action  string
}
type Response struct {
	Message      []byte
	Conn         net.Conn
	SendToSender bool
}

