package connection

import (
	"fmt"
	"time"
	"database/sql"
	"test/Helpers"
)


func Server(Ser *helpers.Authentication, db *sql.DB) {
	 
	// idont know why this it's work just with loop while he is just one message who is pass in channel!
	// now i understand  whi i use loop with channel becs ti need to handle multi message who is come form client 
	// m := <-ser.message
	for m := range Ser.Msg {
		timestamp := time.Now().Format("[2006-01-02 15:04:05]") //[YYYY-MM-DD HH:MM:SS]
		// Format message !!
		result := fmt.Sprintf("%s[%s]:%s", timestamp, m.Login, m.Content)
		namegroube := m.Groube
		fmt.Println(namegroube)
		fmt.Println(result)
		_, err := db.Exec("INSERT INTO messages (groubname, content) VALUES (?, ?);", namegroube ,result)
		if err != nil {
			fmt.Println("Not insert data into the table")
		}
		Ser.Mu.Lock()
		Thetargetone := Ser.Con[m.Groube]
		Ser.Mu.Unlock()
		for name, link := range Thetargetone {
			// her i skip the user who is send the message !!
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
