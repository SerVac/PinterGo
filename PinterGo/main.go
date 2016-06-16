package main

import (
	"log"
	"net"
	"os"
	"fmt"
	"bufio"
"strings"
)

const (
	service_name = "myservice"
	service_description = "My Echo Service"

	port = ":8080"
)

var stdlog, errlog *log.Logger

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

func main() {
	/*listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Possibly was a problem with the port binding")
		os.Exit(1)
	}


	listen := make(chan net.Conn, 10)
	go acceptConnection(listener, listen)*/

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
}

func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	// forever loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error!")
			continue
		}
		fmt.Println("server work ok!")
		listen <- conn
	}
}