package nodes

import "sync"

const SYSTEM = 0
const LOGIN = 1
const LOGOUT = 2

type Node interface {
	Write(response *Response, error chan Error, wg sync.WaitGroup)
	Read(request chan <- *Request, error chan Error, wg sync.WaitGroup)
	SetStatus(status int)
	SetName(name string)
	GetStatus() int
	GetName() string
	GetInfo() string
}

type Request struct{
	Msg string
	Node Node
}

type Response struct {
	Msg string
	Node Node
	SendToSender bool
}
type Error struct{
	Node Node
	Error error
}