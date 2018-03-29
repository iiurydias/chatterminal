package nodes_client

import (
	"net"
	"chat-terminal/settings"
	"log"
	"fmt"
	"os"
)

type Network struct {
	conn net.Conn
	name string
}

func NewNetwork(setting *settings.Settings) *Network {
	addr, err := net.ResolveTCPAddr(setting.Network, setting.Address)
	if err != nil {
		log.Fatalln(err.Error())
	}
	conn, err := net.DialTCP(setting.Network, nil, addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &Network{conn, setting.Cname}
}

func (n *Network) Read(conn net.Conn, response chan string) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("SERVER DISCONNECTED")
				os.Exit(0)
				return
			}
			fmt.Println("Error founded: " + err.Error())
			return
		}
		response <- string(buffer[:n])
	}
}
func (n *Network) Write(conn net.Conn, requests chan string) {
	for request := range requests {
		_, err := conn.Write([]byte(request))
		//if string(request) == `{"Action":"logout","Content":""}` {
		//	os.Exit(0)
		//}
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("SERVER DISCONNECTED")
				os.Exit(0)
				return
			}
			fmt.Println("Error founded: " + err.Error())
			return
		}
	}
}
func (n *Network) GetConn() net.Conn{
	return n.conn
}

func (n *Network) GetName() string{
	return n.name
}
