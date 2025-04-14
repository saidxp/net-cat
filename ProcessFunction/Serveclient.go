package connection

import (
	"fmt"
	"io"
	"test/Helpers"
	"time"
)

func Server(Ser *helpers.Authentication ) {

	// idont know why this it's work just with loop while he is just one message who is pass in channel!
	// now i understand  whi i use loop with channel becs it need to handle multi message who is come form client 
	// m := <-ser.message
	for m := range Ser.Msg {
		timestamp := time.Now().Format("[2006-01-02 15:04:05]") //[YYYY-MM-DD HH:MM:SS]
		// Format message !!
		result := fmt.Sprintf("%s[%s]:%s", timestamp, m.Login, m.Content)
		namegroube := m.Groube
		Ser.Mu.Lock()
		n := Ser.Log[namegroube]
		io.WriteString(n, result + "\n")
		Ser.Mu.Unlock()
		// The target connection!
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
			_, err := fmt.Fprintf(link.Conn, "\n"+result+"\n")
			if err != nil {
				fmt.Println("Error writing to %s", m.Login)
			}
		}
		for name, link := range Thetargetone {
			fmt.Fprintf(link.Conn, "%s[%s]:", timestamp, name)
		}
	}
}
