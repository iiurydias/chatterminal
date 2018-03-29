package nodes

import (
	"chat-terminal/logg"
	"strconv"
	"sync"
)

type LogNode struct {
	log *logg.Log
}

func NewLogNode(log *logg.Log) *LogNode {
	return &LogNode{log}
}

func (l *LogNode) Write(response *Response, error chan Error, wg sync.WaitGroup) {
	l.log.Info(response.Msg)
}

func (l *LogNode) Read(request chan<- *Request, error chan Error, wg sync.WaitGroup) {
	return
}

func (l *LogNode) SetName(name string) {
	return
}

func (l *LogNode) SetStatus(status int) {
	return
}

func (l *LogNode) GetName() string {
	return "Log Node"
}

func (l *LogNode) GetStatus() int {
	return SYSTEM
}

func (l *LogNode) GetInfo() string {
	return l.GetName() + " (" + l.log.Name + ") Status: " + strconv.Itoa(l.GetStatus())
}