package connection

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"test/Helpers"
)

func checkgroubeandreturn(client net.Conn, name string, auth *helpers.Authentication) (string, error) {
	var g string
	var gg string
	var err error

	con := bufio.NewReader(client)
	for {
		t := "Welcome u are login now please Entre Wich any groube u want ... 1 or 2 : "
		l := "Hello " + name
		Welcome := fmt.Sprintf("%s : %s ", l, t)
		fmt.Fprintf(client, Welcome)

		g, err = con.ReadString('\n')
		if err != nil {
			fmt.Println("Eroorrr to read the FILE stdin from client")
			return "", err
		}
		gg = strings.TrimSpace(g)
		// her i will check if this name it's deplicate in groube who it's enregister
		if helpers.Checkmap(name, auth, gg) {
			if gg == "1" || gg == "2" {
				// break infinit loop
				break
			} else {
				fmt.Fprintf(client, "Failed try again the name of groube must be between >> 1 and 2\n")
			}
		} else {
			fmt.Fprintf(client, "Failed try again After u entre the name of groube there is a Duplicate Name in the groube u try to joined it\n")
		}
	}
	return gg, nil
}
