package main

import (
	"sync"
	"chat-terminal/client"
	"chat-terminal/settings"
	"chat-terminal/client/nodes_client"
)

func main() {
	settings := settings.NewSettings()
	settings.LoadSettings()
	network := nodes_client.NewNetwork(settings)
	screen := nodes_client.NewScreen()
	client := client.NewClient(network, screen)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.HandleConnection()
	}()
	wg.Wait()
}
