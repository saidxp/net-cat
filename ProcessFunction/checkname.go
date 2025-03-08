package connection

import (
	"fmt"
	"net"
	"bufio"
	"test/Helpers"
	"strings"
)

func Checknameandreturn(client net.Conn) (string, error) {
	con := bufio.NewReader(client)
	var name string
	for {
	fmt.Fprintf(client, "[ENTER YOUR NAME]: ")
	str, err := con.ReadString('\n')
	if err != nil {
		fmt.Println("Error to read data from client !\n")
		fmt.Fprintf(client, "Error: Failed to read your input. Please try again.\n")
		return "" , err
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
	return name, nil
}