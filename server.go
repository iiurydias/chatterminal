package main

import (
	"net"
	"fmt"
	"sync"
	"log"
	"strings"
	"os"
	"io/ioutil"
	"encoding/json"
)

type Request struct {
	Message      []byte
	Conn         net.Conn
	sendToSender bool
}

type Response struct {
	Request
}
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
	go handleNewClients(wg, s)
	wg.Wait()
}
func handleNewClients(wg sync.WaitGroup, s Config) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.Network, s.Address)
	if err != nil {
		log.Fatalln(err.Error())
	}
	listener, err := net.ListenTCP(s.Network, addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	connChan := make(chan net.Conn)
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go acceptClient(connChan, listener, wg2)
	go routeClient(connChan, wg2)
	wg2.Wait()
}
func routeClient(connChan chan net.Conn, wg2 sync.WaitGroup) {
	defer wg2.Done()
	var wg3 sync.WaitGroup
	var clientlist = make(map[net.Conn]string)
	var response = make(chan Response)
	var request = make(chan Request)
	var remChan = make(chan Request)
	wg3.Add(3)
	go writeClientResponse(response, wg3, clientlist)
	go processClientRequest(request, response, wg3, clientlist)
	go removeClientFromList(remChan, wg3, clientlist)
	for conn := range connChan {
		wg3.Add(1)
		go dealWithRequests(conn, wg3, clientlist, request, remChan)
	}
	wg3.Wait()
}
func dealWithRequests(conn net.Conn, wg3 sync.WaitGroup, clientlist map[net.Conn]string, request chan Request, remChan chan Request) {
	defer wg3.Done()
	var wg4 sync.WaitGroup
	wg4.Add(1)
	go readClientRequest(conn, request, remChan, wg4)
	wg4.Wait()
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
func acceptClient(conexao chan net.Conn, listener net.Listener, wg2 sync.WaitGroup) {
	defer wg2.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			plotError(err)
			continue
		}
		conexao <- conn
	}
}
func readAndSaveClientName(conn net.Conn, clientlist map[net.Conn]string) error {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}
	clientlist[conn] = string(buffer[:n])
	return nil
}
func plotError(err error) {
	if err.Error() == "EOF" {
		fmt.Println("CLIENT DISCONNECTED")
		return
	}
	fmt.Println("Error founded: " + err.Error())
}
func readClientRequest(conn net.Conn, request chan<- Request, remChan chan Request, wg4 sync.WaitGroup) {
	defer wg4.Done()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			plotError(err)
			remChan <- Request{Message: []byte("logout"), Conn: conn}
			break
		}
		request <- Request{Message: buffer[:n], Conn: conn}
	}
}
func writeClientResponse(responses <-chan Response, wg3 sync.WaitGroup, clientlist map[net.Conn]string) {
	defer wg3.Done()
	for response := range responses {
		if response.sendToSender {
			_, err := response.Conn.Write([]byte(response.Request.Message))
			if err != nil {
				plotError(err)
				continue
			}
			continue
		}
		for con := range clientlist {
			if response.Conn == con {
				continue
			}
			_, err := con.Write([]byte(response.Request.Message))
			if err != nil {
				plotError(err)
				continue
			}
		}
	}
}
func processClientRequest(requests chan Request, responses chan Response, wg3 sync.WaitGroup, clientlist map[net.Conn]string) {
	defer wg3.Done()
	for request := range requests {
		action := strings.Split(string(request.Message), " ")
		if string(string(action[0][0])) == "@" {
			msg := string(action[0][1:])
			switch string(msg) {
			case "getlist":
				var list string
				for _, client := range clientlist {
					list = list + client + "\n"
				}
				responses <- Response{Request: Request{[]byte(list), request.Conn, true}}
				fmt.Println("ACTION MADE")
				continue
			case "changename":
				if val, ok := clientlist[request.Conn]; ok {
					oldname := val
					clientlist[request.Conn] = string(action[1])
					fmt.Println(strings.ToUpper(oldname) + " CHANGED NAME TO " + strings.ToUpper(clientlist[request.Conn]))
					responses <- Response{Request: Request{[]byte(strings.ToUpper(oldname) + " CHANGED NAME TO " + strings.ToUpper(clientlist[request.Conn])), request.Conn, false}}
					continue
				}
				clientlist[request.Conn] = string(action[1])
				fmt.Println(strings.ToUpper(clientlist[request.Conn]) + " CONNECTED")
				continue
			case "logout":
				responses <- Response{Request: Request{[]byte(strings.ToUpper(clientlist[request.Conn]) + " IS OFFLINE"), request.Conn, false}}
				continue
			}
		}
		responses <- Response{Request: Request{[]byte(clientlist[request.Conn] + ": " + string(request.Message)), request.Conn, false}}
		fmt.Println("MESSAGE SENT")
	}
}

func removeClientFromList(remChan chan Request, wg3 sync.WaitGroup, clientlist map[net.Conn]string) {
	defer wg3.Done()
	for chann := range remChan {
		if string(chann.Message) == "logout" {
			for conn := range clientlist {
				if conn == chann.Conn {
					delete(clientlist, conn)
					conn.Close()
					break
				}
			}
			continue
		}

	}
}
