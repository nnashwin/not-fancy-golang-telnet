package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	conn         net.Conn
	userId       string
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

		t := time.Now().Format(time.RFC822)
		fmt.Printf("%+v %+v: %+v", t, c.userId, line)
		ch <- fmt.Sprintf("%+v %+v: %+v", t, c.userId, line)
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

func getUserName(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "Welcome to the Not-so-fancy chat\n")
	io.WriteString(c, "Please input your name\n")
	nick, _ := bufc.ReadString('\n')

	// return a slick of the string which does not include the \n byte at the end of the string
	nickSlice := nick[:len(nick)-2]

	return nickSlice
}

func main() {
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Sprintf("The telnet server has started at %v!", time.Now().Format(time.RFC822))
	if _, err := f.Write([]byte(fmt.Sprintf("The telnet server has started at %v!", time.Now().Format(time.RFC822)))); err != nil {
		log.Fatal(err)
	}

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
	defer f.Close()
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
	// no space after the verb because it seems to add a \n
	msgCChan <- fmt.Sprintf("Howdy ho %s! Welcome to the Telnet Chat!\n", client.userId)

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
			fmt.Printf("%v New client has joined the channel: %v\n", time.Now().Format(time.RFC822), client.userId)
			clients[client.conn] = client.ch

		case client := <-rmCChan:
			fmt.Printf("%v Client disconnects: %v\n", time.Now().Format(time.RFC822), client.conn)
			delete(clients, client.conn)
		}
	}
}
