package server

import "fmt"

func PlotError(err error) {
	if err.Error() == "EOF" {
		fmt.Println("CLIENT DISCONNECTED")
		return
	}
	fmt.Println("Error founded: " + err.Error())
}
