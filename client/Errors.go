package client

import (
	"log"
	"fmt"
)

func CheckFatalError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
func PlotError(err error) {
	if err.Error() == "EOF" {
		fmt.Println("SERVER DISCONNECTED")
		return
	}
	fmt.Println("Error founded: " + err.Error())
}
