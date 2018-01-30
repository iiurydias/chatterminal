package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"sync"
	"strings"
	"log"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Ip      string
	Cname   string
	Port    string
	Address string
	Network string
	File    string
}

func main() {
	var wg sync.WaitGroup
	s := loadSettings()
	wg.Add(1)
	go handleConnection(s, wg)
	wg.Wait()
}
func handleConnection(s Config, wg sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.Network, s.Address)
	checkFatalError(err)
	conn, err := net.DialTCP(s.Network, nil, addr)
	checkFatalError(err)
	_, err = conn.Write([]byte("@changename "+ s.Cname))
	checkFatalError(err)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go processResquests(conn, wg2)
	wg2.Wait()
}
func processResquests(conn net.Conn, wg2 sync.WaitGroup) {
	defer wg2.Done()
	var response = make(chan string)
	var request = make(chan string)
	var wg3 sync.WaitGroup
	wg3.Add(4)
	go printResponse(response, wg3)
	go readServerResponse(conn, response, wg3)
	go writeServerRequest(conn, request, wg3)
	go writeRequest(request, wg3)
	wg3.Wait()
}
func loadSettings() Config {
	var s Config
	s.Network = "tcp"
	s.File = os.Args[1]
	data, err := ioutil.ReadFile(s.File)
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatalln(err)
	}
	s.Address = s.Ip + s.Port
	return s
}
func checkFatalError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
func plotError(err error) {
	if err.Error() == "EOF" {
		fmt.Println("SERVER DISCONNECTED")
		return
	}
	fmt.Println("Error founded: " + err.Error())
}
func readServerResponse(conn net.Conn, response chan<- string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			plotError(err)
			return
		}
		response <- string(buffer[:n])
	}
}
func writeServerRequest(conn net.Conn, requests <-chan string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for request := range requests {
		_, err := conn.Write([]byte(request))
		if string(request) == "@logout" {
			os.Exit(0)
		}
		if err != nil {
			plotError(err)
			return
		}
	}
}
func printResponse(responses chan string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for response := range responses {
		fmt.Println(response)
	}
}
func writeRequest(request chan<- string, wg3 sync.WaitGroup) {
	defer wg3.Done()
	for {
		reader := bufio.NewReader(os.Stdin)
		bText, _ := reader.ReadString('\n')
		bText = strings.Replace(bText, "\n", "", -1)
		request <- bText
	}
}
