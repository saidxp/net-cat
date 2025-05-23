package main

import (
	"fmt"
	"net"
	"os"
	// process "test/ProcessFunction"
	"test/Helpers"
	"test/ProcessFunction"
)

func main() {
	
	port := "8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	// her i reseverd 10 struct inside struct of chan *bufer for stock and limit the number of client who it's connected at the same time !!
	limiteure := make(chan struct{}, 10)
	socket, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("error to connect")
		return
	}
	// Her i will create a file log to enregister state and interaction of Server !
	filelog, err := os.Create("log.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error to creat log file")
	}
	defer filelog.Close()
	// her i face the problem of Initializasion
	Auth := &helpers.Authentication{
		Con: make(map[string]map[string]*helpers.Link),
		Msg: make(chan helpers.Message),
		Log: make(map[string]*os.File),
	}
	filelog.WriteString("Server is starting..." + "\n")
	fmt.Println("Listening on the port", ":"+port)
	
	go connection.Server(Auth)
	for {
		accept, err := socket.Accept()
		if err != nil {
			fmt.Println("Error to connect with client")
		}
		if len(limiteure) == cap(limiteure) {
			accept.Write([]byte("Server is currently full. Please wait until a spot is available...\n"))
		}
		// i add this to *buffer inside struct chan to reserve place "empty struct"
		limiteure <- struct{}{}

		go connection.Handleconection(accept, Auth, limiteure, filelog)
	}
}
