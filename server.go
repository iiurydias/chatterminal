package main

import (
	"chat-terminal/servers"
	"chat-terminal/settings"
	"chat-terminal/servers/nodes"
	"log"
	"chat-terminal/logg"
	"sync"
)

func main() {
	setting := settings.NewSettings()
	setting.LoadSettings()
	P := nodes.NewGatekeeper()
	LisNode := nodes.NewListenerNode(P, setting)
	log1 := logg.NewLogger("IURY'S SERVER")
	screenH := logg.NewScreenHandle(logg.INFO)
	fileH := logg.NewFileHandle(logg.INFO, "teste.txt")
	remoteH := logg.NewRemoteHandle(logg.INFO, "tcp", "192.168.40.245:8080")
	log1.AddHandle(screenH)
	log1.AddHandle(fileH)
	log1.AddHandle(remoteH)
	Snode := nodes.NewLogNode(log1)
	P.AddNode(LisNode)
	P.AddNode(Snode)
	var wg sync.WaitGroup
	wg.Add(1)
	request, errPRead := P.Read(wg)
	LogErrors(errPRead)
	server := servers.NewServer(P)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	res, errServer := server.Process(request, wg2)
	LogErrors(errServer)
	var wg3 sync.WaitGroup
	wg3.Add(1)
	errPWrite := P.Write(res, wg3)
	LogErrors(errPWrite)
	wg.Wait()
	wg2.Wait()
	wg3.Wait()
}

func LogErrors(errChannel chan nodes.Error){
	go func(){
		for err := range errChannel {
			if err.Error.Error() == "EOF"{
				err.Node.SetStatus(nodes.LOGOUT)
			}
			log.Print(err.Error.Error())
		}
	}()
}