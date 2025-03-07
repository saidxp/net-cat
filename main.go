package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"test/Helpers"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	if len(os.Args) != 1 {
		var err error
		db, err = sql.Open("sqlite3", "message.db")
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
			Con: make(map[string]map[string]*helpers.Link),
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
	 
	// idont know why this it's work just with loop while he is just one message who is pass in channel!
	// now i understand  whi i use loop with channel becs ti need to handle multi message who is come form client 
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
		Ser.Mu.Lock()
		Thetargetone := Ser.Con[m.Groube]
		Ser.Mu.Unlock()
		for name, link := range Thetargetone {
			// her i skiip the user who is send the message !!
			if name == m.Login {
				continue
			}
			// her i was face the problem of not send to all client !!
			// the problem of new line ~ !!!
			//_, err := link.Conn.Write([]byte("\n" + result + "\n"))
			_, err := fmt.Fprintf(link.Conn, "\n"+result+"\n")
			if err != nil {
				fmt.Println("Error writing to ", m.Login)
			}
		}
		for name, link := range Thetargetone {
			fmt.Fprintf(link.Conn, "%s[%s]:", timestamp, name)
		}
	}
}

func handleconection(client net.Conn, auth *helpers.Authentication, limiteure chan struct{}) {
	defer client.Close()
	con := bufio.NewReader(client)
	s, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println("Error to Read data")
	}
	var name string
	var g string
	var gg string
	client.Write([]byte(s))
	// I will check if it is a valide name who dont conatin special char !!
	// Check and check until the user entre correct Name!!
	for {
		fmt.Fprintf(client, "[ENTER YOUR NAME]: ")
		str, err := con.ReadString('\n')
		if err != nil {
			fmt.Println("Error to read data from client !\n")
			fmt.Fprintf(client, "Error: Failed to read your input. Please try again.\n")
		}
		// her i will handle the name of user to not be special char !!
		if helpers.Check(str) {
			// -- > i trim the '\n'
			name = strings.TrimSpace(str)
			// here i check if the name is enregister at map !!
			if helpers.Checkmap(name, auth) {
				break
			} else {
				fmt.Fprintf(client, "Invalid name. ther duplicate name try again with diff name !!`\n")
			}
		} else {
			fmt.Fprintf(client, "Invalid name. Please avoid special characters\n")
		}
	}
	// var gg string
	// Her i enregister the users and There connection into map !!
	// i need to be sure the name of the >client not duplicate in this case of this Mini chat !!
	// i
	print(name+"\n")
	for {
		t := "Welcome u are login now please Entre Wich any groube u want ... 1 or 2"
		l := "Hello " + name
		Welcome := fmt.Sprintf("%s %s : ", l, t)
		fmt.Fprintf(client, Welcome)
		g, err = con.ReadString('\n')
		fmt.Printf("the messag %s\n",g)
		if err != nil {
			fmt.Println("Eroorrr to read the the stdin from client")
		}
		gg = strings.TrimSpace(g)
		fmt.Printf("AFTER THE TRIM SPACE %s\n",gg)
		if gg == "1" || gg == "2" {
			fmt.Printf("I BREAK HEREA THE LOOP\n")
			break
		} else {
			fmt.Fprintf(client, "Failed try again the name of groube must be between >> look under |")
		}
	}
	fmt.Printf("the name of groube %s\n",gg)
	fmt.Println(">>>",g)
	 
	// after the client entre the name and groube i need to enregister them !!
	// i face problem her !!
	// i lock all thread to accesse to this shared map ! write !
	// if i dont check if the groube it's exeist or not there is aprobelm is when user log the entire map is deleted and it's initialzation new one !
	auth.Mu.Lock()
	if _, exists := auth.Con[gg]; !exists {
		auth.Con[gg] = make(map[string]*helpers.Link) // Only initialize if group doesnt exist
	}
	fmt.Printf("the name of Groube : %s the members of groube %s\n", gg, name)
	auth.Con[gg][name] = &helpers.Link{Conn: client} // Add user to the group
	auth.Mu.Unlock()

	// >>>>>>>>>>>>>>>>>>>>>>>
	times := time.Now().Format("[2006-01-02 15:04:05]")
	result := fmt.Sprintf("%s[%s]:", times, name)
	// >>>>>>>>>>>>>>>>>>>>>>>
	// here im still need to spesific the chat in spesific Groube !!
	conversion, err := helpers.Getmessages(db, g)
	if err != nil {
		fmt.Println("erro to get message from data base")
	}
	for _, m := range conversion {
		fmt.Fprintf(client, m+"\n")
	}
	fmt.Fprintf(client, result)
	// Her i send Welcome message to every one login into server !!
	Thetargetone := auth.Con[gg]
	fmt.Println("Group Members:", auth.Con[gg])
	// her i specific the groube of client where he is enregister !!
	for n, link := range Thetargetone {
		if n != name {
			link.Conn.Write([]byte("\n" + name + " has joined our chat..." + "\n"))
		}
	}
	//
	for n, link := range Thetargetone {
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
			theone  := auth.Con[gg]

			for _, link := range theone {
				link.Conn.Write([]byte(s))
			}
			for n, link := range theone {
				times := time.Now().Format("[2006-01-02 15:04:05]")
				result := fmt.Sprintf("%s[%s]:", "\n"+times, n)
				fmt.Fprintf(link.Conn, result)
			}
			// I Move the reservation from channel!!
			<-limiteure
			// i delete the conection of the client who is left !
			delete(auth.Con[g], name)
			// TODO: send message to all client if the client it's left the chatroom!
			break
		}
		fmt.Printf("string : %q\n", chat)
		fmt.Printf("%d", len(chat))
		// her i will check the shape of message !!
		if len(chat) > 1 && helpers.CheckMessage(chat) {
			cha := strings.TrimSpace(chat)
			fmt.Printf("string : %q\n", chat)
			auth.Msg <- helpers.Message{
				Login:   name,
				Content: cha,
				Groube:  gg,
			}
		} else {
			fmt.Fprintf(client, "Invalid message it's contain special character!! \"\n\" pleas send a message with what human understand\n")
			times := time.Now().Format("[2006-01-02 15:04:05]")
			result := fmt.Sprintf("%s[%s]:", times, name)
			fmt.Fprintf(client, result)
			continue
		}
	}
}
