package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	// process "test/ProcessFunction"
	"test/Helpers"
	"test/ProcessFunction"

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
		//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (content TEXT NOT NULL);`)
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
        		id INTEGER PRIMARY KEY AUTOINCREMENT,
        		groubname TEXT NOT NULL,
        		content TEXT NOT NULL
    		);
		`)
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
		go connection.Server(Auth, db)
		for {
			accept, err := socket.Accept()
			if err != nil {
				fmt.Println("Error to connect with client")
			}
			// i add this to *buffer inside struct chan to reserve place "empty struct"
			limiteure <- struct{}{}

			go connection.Handleconection(accept, Auth, limiteure, db)
		}
	} else {
		fmt.Println("Invalid Argument!!")
	}
}
