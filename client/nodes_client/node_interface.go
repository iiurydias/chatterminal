package nodes_client

import "net"

type Node interface {
	Write(conn net.Conn, response chan string)
	Read(conn net.Conn, requests chan string)
	GetConn() net.Conn
	GetName() string
}
