package connection

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	helpers "test/Helpers"
)

func Handleconection(client net.Conn, auth *helpers.Authentication, limiteure chan struct{}, logfile *os.File) {
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

	t := time.Now().Format("[2006-01-02 15:04:05]")

	ffff := fmt.Sprintf("Info: the user `%q` IP : %s is connected successfuly to server AT %s", name, client.RemoteAddr().String(), t)
	logfile.WriteString(ffff + "\n")
	if err != nil {
		fmt.Println("Error to read the name from user")
	}
	print(name + "\n")
	// her i need to handle the name of groube
	gg, err = checkgroubeandreturn(client, name, auth)
	if err != nil {
		fmt.Print("ERRO TO READ THE NAME OF GROUBE")
	}
	fmt.Printf("the name of groube %s\n", gg)
	// after the client entre the name and groube i need to enregister them !!
	// i face problem her !!
	// i lock all thread to accesse to this shared map ! write !
	// if i dont check if the groube it's exeist or not there is aprobelm is when user login the entire map is deleted and it's initialzation new one !
	auth.Mu.Lock()
	if _, exists := auth.Con[gg]; !exists {
		auth.Con[gg] = make(map[string]*helpers.Link) // Only initialize if group doesnt exist
	}
	// i will creat a file by name of groube and i will check if is exist or no !
	// BUT IF IS execist i need to check !
	if !helpers.Exists(gg + ".txt") {
		f, err := os.Create(gg + ".txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error \"At Creation File\"occurred:\"%v\n\"", err)
			// kill Process >> File Not Create
			os.Exit(1)
		}
		// Here i add a name of goube to map with her File for enregister the message of this groube !
		fmt.Print("--- AT Enregister the file at map")
		auth.Log[gg] = f
	}
	// --------->>>>
	fmt.Printf("The Name of Groube : %s The Members of Groube %s\n", gg, name)
	auth.Con[gg][name] = &helpers.Link{Conn: client} // Add user to the group
	bbbb := fmt.Sprintf("User %s is joined at Groube %s", name, gg)
	logfile.WriteString(bbbb + "\n")
	auth.Mu.Unlock()
	// >>>>>>>>>>>>>>>>>>>>>>>
	// >>>>>>>>>>>>>>>>>>>>>>>
	times := time.Now().Format("[2006-01-02 15:04:05]")
	result := fmt.Sprintf("%s[%s]:", times, name)
	// >>>>>>>>>>>>>>>>>>>>>>>
	logmessage := auth.Log[gg]
	logmessage.Seek(0, io.SeekStart)
	fl, err := io.ReadAll(logmessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %v\n", err)
	}

	client.Write(fl)
	// here im still need to spesific the chat in spesific Groube !!
	/*conversion, err := helpers.Getmessages(db, gg)
	if err != nil {
		fmt.Println("erro to get message from data base")
	}
	for _, m := range conversion {
		fmt.Fprintf(client, m+"\n")
	}
	*/
	fmt.Fprintf(client, result)
	// Her i send Welcome message to every one login into server with spesific the target groube !!
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
			t := time.Now().Format("[2006-01-02 15:04:05]")
			ffff := fmt.Sprintf("info : The user %s is log out successfuly AT : %s", name, t)
			logfile.WriteString(ffff + "\n")
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
