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

// JSON from which to load the config
const configJson = "config.json"

type Message struct {
	sender string
	text   string
}

func main() {
	config, err := loadConfigFile(configJson)
	if err != nil {
		fmt.Printf("The config file at %v was not found\n", configJson)
		log.Fatal(err)
	}

	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("The log file at %v was not found and is unable to be created\n", config.LogFile)
		log.Fatal(err)
	}

	writeToLog := openLogFile(f)

	l, err := net.Listen("tcp", config.Ip+":"+config.Port)
	if err != nil {
		writeToLog(err.Error())
		log.Fatal(err)
	}

	startMsg := fmt.Sprintf("Server has started on %v:%v at %v!\n", config.Ip, config.Port, time.Now().Format(time.RFC822))
	writeToLog(startMsg)

	addClientChan := make(chan Client)
	msgClientChan := make(chan Message)
	commClientChan := make(chan Command)
	rmClientChan := make(chan Client)

	go handleMsgs(msgClientChan, addClientChan, rmClientChan, writeToLog)
	go handleCommands(commClientChan, writeToLog)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			writeToLog(err.Error())
			// continues the server so one error with a connection will not crash the entire chat server
			continue
		}

		go handleConnection(conn, msgClientChan, addClientChan, rmClientChan, commClientChan)
	}
	defer f.Close()
}

func getUserName(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "Welcome to the Not-so-fancy chat\n")
	io.WriteString(c, "Please input your name\n")
	nick, _ := bufc.ReadString('\n')

	// return a slick of the string which does not include the \n byte at the end of the string
	nickSlice := nick[:len(nick)-2]

	return nickSlice
}

// partially applied os.File parameter.  creates another function which can be called with the string to write
func openLogFile(file *os.File) func(string) {
	return func(str string) {
		fmt.Println(str)
		if _, err := file.Write([]byte(str)); err != nil {
			log.Fatal(err)
		}
	}
}

func handleConnection(c net.Conn, msgCChan chan<- Message, addCChan chan<- Client, rmCChan chan<- Client, commCChan chan<- Command) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	client := Client{
		conn:         c,
		userId:       getUserName(c, bufc),
		ch:           make(chan string),
		blockedUsers: make(map[string]int),
		// Puts everyone in the same room for now.  Will be able to switch in the future
		currentRm: "home",
	}

	addCChan <- client

	defer func() {
		fmt.Printf("Closed connection from %v", c.RemoteAddr())
		rmCChan <- client
	}()

	msgCChan <- Message{sender: "host", text: fmt.Sprintf("Howdy ho %s! Welcome to the Telnet Chat!\n", client.userId)}
	go client.ReadLines(msgCChan, commCChan)
	client.WriteLines(client.ch)
}

func handleMsgs(msgCChan <-chan Message, addCChan <-chan Client, rmCChan <-chan Client, logFunc func(string)) {
	clients := make(map[net.Conn]Client)
	for {
		select {
		case msg := <-msgCChan:
			logFunc(msg.text)
			for _, client := range clients {
				if client.blockedUsers[msg.sender] == 0 {
					go func(mesch chan<- string) {
						mesch <- msg.text
					}(client.ch)
				}
			}

		case client := <-addCChan:
			joinMsg := fmt.Sprintf("%v New client has joined the channel: %v\n", time.Now().Format(time.RFC822), client.userId)
			// msgCChan <- Message{sender: "host", text: joinMsg}
			logFunc(joinMsg)
			clients[client.conn] = client

		case client := <-rmCChan:
			leaveMsg := fmt.Sprintf("%v Client disconnects: %v\n", time.Now().Format(time.RFC822), client.userId)
			// msgCChan <- Message{sender: "host", text: leaveMsg}
			logFunc(leaveMsg)
			delete(clients, client.conn)
		}
	}
}

func handleCommands(commandCCh <-chan Command, logFunc func(string)) {
	for {
		select {
		case command := <-commandCCh:
			fmt.Printf("%+v", command)
		}
	}
}
