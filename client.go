package main

import (
	"sync"
	"project/client"
)

func main() {
	var wg sync.WaitGroup
	s := client.LoadSettings()
	wg.Add(1)
	go client.HandleConnection(s, wg)
	wg.Wait()
}
