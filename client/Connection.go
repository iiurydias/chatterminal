package client

import (
	"sync"
	"net"
	"project/settings"
)

func HandleConnection(s settings.Settings, wg sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.Network, s.Address)
	CheckFatalError(err)
	conn, err := net.DialTCP(s.Network, nil, addr)
	CheckFatalError(err)
	_, err = conn.Write([]byte(`{"Action":"changename","Content":"` + s.Cname + `"}`))
	CheckFatalError(err)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go ProcessRequests(conn, wg2)
	wg2.Wait()
}
