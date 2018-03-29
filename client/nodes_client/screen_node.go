package nodes_client

import (
	"bufio"
	"os"
	"strings"
	"encoding/json"
	"chat-terminal/protocols"
	"fmt"
	"net"
)

type Screen struct {
	msg protocols.MessageProtocol
}

func NewScreen() *Screen {
	return &Screen{protocols.MessageProtocol{"", ""}}
}
func (screen *Screen) Read(conn net.Conn, request chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		bText, _ := reader.ReadString('\n')
		bText = strings.Replace(bText, "\n", "", -1)
		action := strings.Split(string(bText), " ")
		if string(action[0][0]) == "@" {
			switch string(action[0][1:]) {
			case "changename":
				if len(action) < 2 {
					action = append(action, "empty")
				}
				screen.msg = protocols.MessageProtocol{string(action[0][1:]), string(action[1])}
				break
			default:
				screen.msg = protocols.MessageProtocol{string(action[0][1:]), ""}
				break
			}
			b, _ := json.Marshal(screen.msg)
			request <- string(b)
			continue
		}
		screen.msg = protocols.MessageProtocol{"sendMessage", strings.Join(action, " ")}
		b, _ := json.Marshal(screen.msg)
		request <- string(b)
	}
}

func (screen *Screen) Write(conn net.Conn, responses chan string) {
	for response := range responses {
		var m protocols.MessageProtocol
		json.Unmarshal([]byte(response), &m)
		fmt.Println(m.Content)
	}
}
func (screen *Screen) GetConn() net.Conn{
	return nil
}

func (screen *Screen) GetName() string{
	return ""
}