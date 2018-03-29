package nodes

import (
	"log"
	"strconv"
	//"errors"
	"sync"
)

type Gatekeeper struct {
	nodeList []Node
}

func NewGatekeeper() *Gatekeeper {
	return &Gatekeeper{make([]Node, 0)}
}

func (gatekeeper *Gatekeeper) AddNode(node Node) {
	log.Println(node.GetInfo() + " was addded on the list")
	gatekeeper.nodeList = append(gatekeeper.nodeList, node)
}

func (gatekeeper *Gatekeeper) Read(wg sync.WaitGroup) (chan *Request, chan Error) {
	request := make(chan *Request)
	err := make(chan Error)
	var wg2 sync.WaitGroup
	for _, node := range gatekeeper.nodeList {
		wg2.Add(1)
		node.Read(request, err, wg)
	}
	return request, err
}

func (gatekeeper *Gatekeeper) Write(Response chan *Response, wg3 sync.WaitGroup) (chan Error) {
	defer wg3.Done()
	err := make(chan Error)
	var wg4 sync.WaitGroup
	for res := range Response {
		log.Println(res.Node.GetInfo() + " wants to responde with: \"" + res.Msg + "\" SendToSender: " + strconv.FormatBool(res.SendToSender))
		//if res.Node.GetStatus() == LOGOUT {
		//	err <- errors.New(res.Node.GetInfo() + " is off-line and will not respond anything")
		//	continue
		//}
		if res.Node.GetStatus() == LOGOUT && res.SendToSender == true {
			wg4.Add(1)
			res.Node.Write(res, err, wg4)
			continue
		}
		for _, node := range gatekeeper.nodeList {
			if node.GetStatus() != LOGOUT {
				wg4.Add(1)
				node.Write(res, err, wg4)
			}
		}
	}
	return err
}

func (gatekeeper *Gatekeeper) GetLoggedNodes() []Node {
	var list []Node
	for _, node := range gatekeeper.nodeList {
		if node.GetStatus() != LOGOUT && node.GetStatus() != SYSTEM {
			list = append(list, node)
		}
	}
	return list
}