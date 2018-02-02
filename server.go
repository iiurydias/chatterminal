package main

import (
	"project/server"
	"sync"
)

//devolver a mensagem para o client, em json, fazer unmarshal do objeto da msm forma... (devolver igual
// com action e content).
func main() {
	var wg sync.WaitGroup
	s := server.LoadSettings()
	wg.Add(1)
	go server.HandleNewClients(wg, s)
	wg.Wait()
}
