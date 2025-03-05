package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strings"
	"test/Helpers"
	"time"

	_ "github.com/mattn/go-sqlite3"
	 
)

var db *sql.DB

func main() {
	if len(os.Args) != 1 {
		var err error
		db , err = sql.Open("sqlite3", "message.db")
		if err != nil {
			fmt.Println("Error connect to db")
		}
		defer db.Close()
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (content TEXT NOT NULL);`)
		if err != nil {
			fmt.Println("opps not create table")
		}
		port := os.Args[1]
		// her i reseverd 10 struct inside struct of chan *bufer for stock and limit the number of client who it's connected at the same time !!
		limiteure := make(chan struct{}, 10)
		socket, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("error to connect")
			return
		}
		// her i face the problem of Initializasion
		Auth := &helpers.Authentication{
			Con: make(map[string]*helpers.Link),
			Msg: make(chan helpers.Message),
			
		}
		go Serveclient(Auth)
		for {
			accept, err := socket.Accept()
			// i add this to *buffer inside struct chan to reserve place "empty struct"
			limiteure <- struct{}{}
			if err != nil {
				fmt.Println("Error to connect with client")
			}
			go handleconection(accept, Auth, limiteure)
		}
	} else {
		fmt.Println("Invalid Argument!!")
	}
}

func Serveclient(Ser *helpers.Authentication) {
	Ser.Mu.Lock()
	// idont know why this it's work just with loop while he is just one message who is pass in channel!
	// m := <-ser.message
	for m := range Ser.Msg {
		timestamp := time.Now().Format("[2006-01-02 15:04:05]") //[YYYY-MM-DD HH:MM:SS]
		// Format message !!
		result := fmt.Sprintf("%s[%s]:%s", timestamp, m.Login, m.Content)
		fmt.Println(result)
		_, err := db.Exec("INSERT INTO messages (content) VALUES (?);", result)
		if err != nil {
			fmt.Println("Not insert data into the table")
		}
		for name, link := range Ser.Con {
			// her i skiip the user who is send the message !!
			if name == m.Login {
				continue
			}
			// her i was face the problem of not send to all client !!
			// the problem of new line ~ !!!
			//_, err := link.Conn.Write([]byte("\n" + result + "\n"))
			_, err:= fmt.Fprintf(link.Conn,"\n" + result + "\n")

			if err != nil {
				fmt.Println("Error writing to ", m.Login)
			}
		}
		for name, link := range Ser.Con {
			fmt.Fprintf(link.Conn, "%s[%s]:", timestamp, name)
		}
	}
	Ser.Mu.Unlock()
}

func handleconection(client net.Conn, auth *helpers.Authentication, limiteure chan struct{}) {

	defer client.Close()
	con := bufio.NewReader(client)
	s, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println("Error to Read data")
	}
	var name string
	client.Write([]byte(s))
	// I will check if it is a valide name who dont conatin special char !!
	// Check and check until the user entre correct Name!!
	for {
		fmt.Fprintf(client, "[ENTER YOUR NAME]: ")
		str, err := con.ReadString('\n')
		if err != nil {
			fmt.Println("Error to read data from client !\n")
		}
		// her i will handle the name of user to not be special char !!
		if helpers.Check(str) {
			// -- > i trim the '\n' 
			name = strings.TrimSpace(str)
			// here i check if the name is enregister at map !!
			if helpers.Checkmap(name, auth) {
				break
			}else {
				fmt.Fprintf(client, "Invalid name. ther duplicate name try again with diff name !!`\n")
			}
		} else {
			fmt.Fprintf(client, "Invalid name. Please avoid special characters\n")
		}
	}
	// Her i enregister the users and There connection into map !!
	// i need to be sure the name of the client not duplicate in this case of this Mini chat !!
	 
		auth.Con[name] = &helpers.Link{Conn: client}
		times := time.Now().Format("[2006-01-02 15:04:05]")
		result := fmt.Sprintf("%s[%s]:", times, name)
		conversion , err := helpers.Getmessages(db)
		if err != nil {
			fmt.Println("erro to get message from data base")
		}
		for _, m := range conversion {
			fmt.Fprintf(client, m + "\n")
		}
		fmt.Fprintf(client, result)
		// Her i send Welcome message to every one login into server !!
		for n, link := range auth.Con {
			if n != name {
				link.Conn.Write([]byte("\n" + name + " has joined our chat..." + "\n"))
			}
		}

		for n, link := range auth.Con {
			if n != name {
				times := time.Now().Format("[2006-01-02 15:04:05]")
				result := fmt.Sprintf("%s[%s]:", times, n)
				fmt.Fprintf(link.Conn, result)
			}
		}
		for {
			chat, err := con.ReadString('\n')
			if err != nil {
				s := "\n" + name + " has left our chat..."
				for _, link := range auth.Con {
					link.Conn.Write([]byte(s))
				}
				for n, link := range auth.Con {
					times := time.Now().Format("[2006-01-02 15:04:05]")
					result := fmt.Sprintf("%s[%s]:", "\n"+times, n)
					fmt.Fprintf(link.Conn, result)
				}
				// I Move the reservation from channel!!
				<-limiteure
				// i delete the conection of the client who is left !
				delete(auth.Con, name)
				// TODO: send message to all client if the client it's left the chatroom!
				break
			}
			fmt.Printf("string : %q\n", chat) 
			fmt.Printf("%d",len(chat))
			// her i will check the shape of message !!
			if len(chat) > 1 && helpers.CheckMessage(chat){
			cha := strings.TrimSpace(chat)
			fmt.Printf("string : %q\n", chat) 
			auth.Msg <- helpers.Message{
				Login:   name,
				Content: cha,
			}
			}else {
				fmt.Fprintf(client, "Invalid message it's contain special character!! \"\n\" pleas send a message with what human understand\n")
				times := time.Now().Format("[2006-01-02 15:04:05]")
				result := fmt.Sprintf("%s[%s]:", times, name)
				fmt.Fprintf(client, result)
				continue
			}
		}
	}

 

