package main

import (
	"bufio"
	"log"
	"net"
)

func handleConnection(c net.Conn) {
	defer c.Close()

	log.Printf("Connection established from %v.\n", c.RemoteAddr())
}

func main() {
	log.Println("Telnet Server!!")

	l, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}
