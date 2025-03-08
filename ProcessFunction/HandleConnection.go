package connection

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"test/Helpers"
)

func Handleconection(client net.Conn, auth *helpers.Authentication, limiteure chan struct{}, db *sql.DB) {
	defer client.Close()
	con := bufio.NewReader(client)
	s, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println("Error to Read data")
	}
	var name string
	var gg string

	client.Write([]byte(s))
	// I will check if it is a valide name who dont conatin special char !!
	// Check and check until the user entre correct Name!!

	// her i check the valide name of user !
	name, err = Checknameandreturn(client)
	if err != nil {
		fmt.Println("Error to read the name from user")
	}
	/*
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
				break
			} else {
				fmt.Fprintf(client, "Invalid name. Please avoid special characters\n")
			}
		}
	*/
	// var gg string
	// Her i enregister the users and There connection into map !!
	// i need to be sure the name of the >client not duplicate in this case of this Mini chat !!
	// i
	print(name + "\n")
	// her i need to handle the name of groube

	gg, err = checkgroubeandreturn(client, name, auth)
	if err != nil {
		fmt.Print("ERRO TO READ THE NAME OF GROUBE")
	}

	/*for {
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
		// her i will check if this name it's deplicate in groube who it's enregister
		if helpers.Checkmap(name, auth, gg) {
		fmt.Printf("AFTER THE TRIM SPACE %s\n",gg)
		if gg == "1" || gg == "2" {
			fmt.Printf("I BREAK HERE THE LOOP\n")
			break
		} else {
			fmt.Fprintf(client, "Failed try again the name of groube must be between >> 1 and 2\n")
		}
		}else {
			fmt.Fprintf(client, "Failed try again After u entre the name of groube there is a Duplicate Name in the groube u try to joined it\n")

		}
	}
	*/
	fmt.Printf("the name of groube %s\n", gg)
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
	conversion, err := helpers.Getmessages(db, gg)
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
			theone := auth.Con[gg]

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
			delete(auth.Con[gg], name)
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
			fmt.Fprintf(client, "Invalid message it's contain special character!! \"\n\" please send a message with what human understand\n")
			times := time.Now().Format("[2006-01-02 15:04:05]")
			result := fmt.Sprintf("%s[%s]:", times, name)
			fmt.Fprintf(client, result)
			continue
		}
	}
}
