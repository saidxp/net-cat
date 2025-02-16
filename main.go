package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)
type yew struct {
	str string
}
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
	Login string
	content string
}

func main() {
	port := os.Args[1]
	socket, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("error to connect")
		return
	}
	// her i face the problem of Initializasion
	auth := &Authentication{
		con: make(map[string]*Link) ,
		msg: make(chan Message),
	}
	go auth.serveclient()
	for {
		accept, err := socket.Accept()
		if err != nil {
			fmt.Println("ello fuck")
		}
		go handleconection(accept, auth)
	}
}

func (ser *Authentication) serveclient() {
		 
		for m := range ser.msg {
		ser.mu.RLock()
		// her i Merge the name of client and her message !!
		result := m.Login + m.content
		for _, link := range ser.con {
			// her i was face the problem of not send to all client !!
			_, err := link.conn.Write([]byte(result + "\n"))
			if err != nil {
				fmt.Println("Error writing to ",m.Login)
			}
			}
		}
		ser.mu.RUnlock()
}

func handleconection(client net.Conn, auth *Authentication) {
	w := yew{
		str: 'Welcome to TCP-Chat!
		_nnnn_
	   dGGGGMMb
	  @p~qp~~qMb
	  M|@||@) M|
	  @,----.JM|
	 JS^\\__/  qKL
	dZP        qKRb
   dZP          qKKb
  fZP            SMMb
  HZM            MMMM
  FqM            MMMM
__| ".        |\dS"qML
|    `.       | `' \\Zq
_)      \\.___.,|     .'
\\____   )MMMMMP|   .'
	`-'       `--'
[ENTER YOUR NAME]:',
	}
	defer client.Close()
	con := bufio.NewReader(client)
	client.Write([]byte(w.str))
	str, err := con.ReadString('\n')
	if err != nil {
		fmt.Println("hello dumb")
	}
	name := strings.TrimSpace(str)
	// her i enregister the users and her connection into map !!
	auth.con[name] = &Link{conn: client}
	// her i send Welcome message to every one login into server !!
	auth.msg <- Message {
		Login: "SUDO",
		content: "You welcome to groube of chat" + name ,
	}
	for {
		chat , err := con.ReadString('\n')
		if err != nil {
			fmt.Println("there is no connection the user left groube")
		}
		cha := strings.TrimSpace(chat)
		auth.msg <- Message{
			Login: name,
			content: cha,
		}
	}
}
