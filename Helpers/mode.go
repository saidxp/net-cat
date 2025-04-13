package helpers

import "sync"
import "net"
import "os"

type Authentication struct {
	Con map[string]map[string]*Link
	Mu  sync.Mutex
	// chanel to passe the data trought multi processe of data in concurency of go routine !!
	Msg chan Message
	Log map[string]*os.File
}

type Link struct {
	Conn net.Conn
}
type Message struct {
	Login   string
	Content string
	Groube string
}
