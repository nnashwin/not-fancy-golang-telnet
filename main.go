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
	userId       []byte
	ch           chan string
	blockedUsers []string
	currentRm    string
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

func getUserName(c net.Conn, bufc *bufio.Reader) []byte {
	io.WriteString(c, "Welcome to the Not-so-fancy chat")
	io.WriteString(c, "Please input your name")
	nick, _ := bufc.ReadString()
	fmt.Println(nick)
	return []byte(nick)
}

func main() {
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

		go handleConnection(conn, msgClientChan, addClientChan, rmClientChan)
	}
}

func handleConnection(c net.Conn, msgCChan chan<- string, addCChan chan<- Client, rmCChan chan<- Client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	client := Client{
		conn:      c,
		userId:    getUserName(c, bufc),
		ch:        make(chan string),
		currentRm: "home",
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
			for _, ch := range clients {
				go func(mesch chan<- string) {
					mesch <- msg
				}(ch)
			}

		case client := <-addCChan:
			fmt.Printf("New client has joined the channel: %v\n", client.conn)
			clients[client.conn] = client.ch

		case client := <-rmCChan:
			fmt.Printf("Client disconnects: %v\n", client.conn)
			delete(clients, client.conn)
		}
	}
}
