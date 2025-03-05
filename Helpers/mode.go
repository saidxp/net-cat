package helpers

import "sync"
import "net"

type Authentication struct {
	Con map[string]*Link
	Mu  sync.Mutex
	// chanel to passe the data trought multi processe of data in concurency of go routine !!
	Msg chan Message
}

type Link struct {
	Conn net.Conn
}
type Message struct {
	Login   string
	Content string
}
