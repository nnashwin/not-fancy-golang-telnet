package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	flag.Parse()
	port := ":" + flag.Arg(0)

	if port == ":" {
		port := ":9001"
	}

	log.Println("Telnet Server!!")

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
}
