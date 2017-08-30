package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	conn         net.Conn
	userId       string
	ch           chan string
	blockedUsers []string
}

func (c Client) ReadLines(ch chan<- string) {
	buffc := bufio.NewReader(c.conn)

	for {
		line, err := buffc.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		ch <- fmt.Sprintf("%s: %s", c.userId, line)
	}
}

func (c Client) WriteLines(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	clientCount := 0
	log.Println("Telnet Server!!")

	l, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}

	addClientChan := make(chan Client)
	msgClientChan := make(chan string)
	rmClientChan := make(chan Client)

	go handleMsgs(msgClientChan, addClientChan, rmClientChan)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		nick := "Client-" + string(clientCount)
		clientCount++
		go handleConnection(conn, msgClientChan, addClientChan, rmClientChan, nick)
	}
}

func handleConnection(c net.Conn, msgCChan chan<- string, addCChan chan<- Client, rmCChan chan<- Client, nick string) {
	defer c.Close()

	client := Client{
		conn:   c,
		userId: nick,
		ch:     make(chan string),
	}

	addCChan <- client
	io.WriteString(c, "Howdy Ho, Welcome to the Telnet Chat")
	msgCChan <- fmt.Sprintf("A new user has joined the CHAT!!")

	go client.ReadLines(msgCChan)
	client.WriteLines(client.ch)
}

func handleMsgs(msgCChan <-chan string, addCChan <-chan Client, rmCChan <-chan Client) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case msg := <-msgCChan:
			log.Println(msg)
			for _, ch := range clients {
				go func(mesch chan<- string) {
					mesch <- msg
				}(ch)
			}

		case client := <-addCChan:
			log.Printf("New client has joined the channel: %v\n", client.conn)
			clients[client.conn] = client.ch

		case client := <-rmCChan:
			log.Printf("Client disconnects: %v\n", client.conn)
			delete(clients, client.conn)
		}
	}
}
