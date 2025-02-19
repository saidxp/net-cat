package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Authentication struct {
	con map[string]*Link
	mu  sync.RWMutex
	// chanel to passe the data trought multi processe of data in concurency of go routine !!
	msg chan Message
}

type Link struct {
	conn net.Conn
}
type Message struct {
	Login   string
	content string
}

func main() {
	port := os.Args[1]
	// her i reseverd 10 struct inside struct of chan *bufer for stock and limit the number of client who it's connected at the same time !!
	limiteure := make(chan struct{}, 10)
	socket, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("error to connect")
		return
	}
	// her i face the problem of Initializasion
	auth := &Authentication{
		con: make(map[string]*Link),
		msg: make(chan Message),
	}
	go auth.serveclient()
	for {
		accept, err := socket.Accept()
		// i add this to *buffer inside struct chan to reserve place "empty struct"
		limiteure <- struct{}{}
		if err != nil {
			fmt.Println("Error to connect with client")
		}
		go handleconection(accept, auth, limiteure)
	}
}

func (ser *Authentication) serveclient() {
	ser.mu.Lock()
	// idont know why this it's work just with loop while he is just one message who is pass in channel!
	// m := <-ser.message
	for m := range ser.msg {
		timestamp := time.Now().Format("[2006-01-02 15:04:05]") //[YYYY-MM-DD HH:MM:SS]
		// Format message !!
		result := fmt.Sprintf("%s [%s]: %s", timestamp, m.Login, m.content)
		for name, link := range ser.con {
			// her i skiip the user who is send the message !!
			if name == m.Login {
				continue
			}
			// her i was face the problem of not send to all client !!
			// the problem of new line ~ !!!
			_, err := link.conn.Write([]byte("\n" + result + "\n"))
			if err != nil {
				fmt.Println("Error writing to ", m.Login)
			}
		}
		for name, link := range ser.con {
			fmt.Fprintf(link.conn, "%s [%s]: ", timestamp, name)
		}
	}
	ser.mu.Unlock()
}

func handleconection(client net.Conn, auth *Authentication, limiteure chan struct{}) {
	defer client.Close()
	con := bufio.NewReader(client)
	s, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println("Error to Read data")
	}
	client.Write([]byte(s))
	fmt.Fprintf(client, "[ENTER YOUR NAME]: ")
	str, err := con.ReadString('\n')
	if err != nil {
		fmt.Println("Error to read data from client !")
	}
	name := strings.TrimSpace(str)
	// Her i enregister the users and There connection into map !!
	auth.con[name] = &Link{conn: client}
	times := time.Now().Format("[2006-01-02 15:04:05]")
	result := fmt.Sprintf("%s [%s]: ", times, name)
	fmt.Fprintf(client, result)
	// Her i send Welcome message to every one login into server !!
	for n, link := range auth.con {
		if n != name {
			link.conn.Write([]byte("\n" + name + " has joined our chat..." + "\n"))
		}
	}
	for n, link := range auth.con {
		if n != name {
			times := time.Now().Format("[2006-01-02 15:04:05]")
			result := fmt.Sprintf("%s [%s]: ", times, n)
			fmt.Fprintf(link.conn, result)
		}
	}
	for {
		chat, err := con.ReadString('\n')
		if err != nil {
			s := "\n" + name + " has left our chat..."
			for _, link := range auth.con {
				link.conn.Write([]byte(s))
			}
			for n, link := range auth.con {
				times := time.Now().Format("[2006-01-02 15:04:05]")
				result := fmt.Sprintf("%s [%s]: ", "\n"+times, n)
				fmt.Fprintf(link.conn, result)
			}
			// I Move the reservation from channel!!
			<-limiteure
			// i delete the conection of the client who is left !
			delete(auth.con, name)
			// TODO: send message to all client if the client it's left the chatroom!
			break
		}
		cha := strings.TrimSpace(chat)
		auth.msg <- Message{
			Login:   name,
			content: cha,
		}
	}
}
