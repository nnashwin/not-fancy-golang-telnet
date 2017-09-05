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

func main() {
	config, err := loadConfigFile(configJson)
	if err != nil {
		fmt.Printf("The config file at %v was not found\n", configJson)
		log.Fatal(err)
	}

	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("The log file %v was not found\n", config.LogFile)
		log.Fatal(err)
	}

	writeToLog := openLogFile(f)

	l, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		writeToLog(err.Error())
		log.Fatal(err)
	}

	startMsg := fmt.Sprintf("The telnet server has started at %v!\n", time.Now().Format(time.RFC822))
	writeToLog(startMsg)

	addClientChan := make(chan Client)
	msgClientChan := make(chan string)
	rmClientChan := make(chan Client)

	go handleMsgs(msgClientChan, addClientChan, rmClientChan, writeToLog)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			writeToLog(err.Error())
			// continues the server so one error with a connection will not crash the entire chat server
			continue
		}

		go handleConnection(conn, msgClientChan, addClientChan, rmClientChan)
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
		if _, err := file.Write([]byte(str)); err != nil {
			log.Fatal(err)
		}
	}
}

func handleConnection(c net.Conn, msgCChan chan<- string, addCChan chan<- Client, rmCChan chan<- Client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	client := Client{
		conn:   c,
		userId: getUserName(c, bufc),
		ch:     make(chan string),
		// Puts everyone in the same room for now.  Will be able to switch in the future
		currentRm: "home",
	}

	addCChan <- client
	msgCChan <- fmt.Sprintf("Howdy ho %s! Welcome to the Telnet Chat!\n", client.userId)

	go client.ReadLines(msgCChan)
	client.WriteLines(client.ch)
}

func handleMsgs(msgCChan <-chan string, addCChan <-chan Client, rmCChan <-chan Client, writeToLog func(string)) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case msg := <-msgCChan:
			writeToLog(msg)
			for _, ch := range clients {
				go func(mesch chan<- string) {
					mesch <- msg
				}(ch)
			}

		case client := <-addCChan:
			joinMsg := fmt.Sprintf("%v New client has joined the channel: %v\n", time.Now().Format(time.RFC822), client.userId)
			writeToLog(joinMsg)
			clients[client.conn] = client.ch

		case client := <-rmCChan:
			leaveMsg := fmt.Sprintf("%v Client disconnects: %v\n", time.Now().Format(time.RFC822), client.userId)
			writeToLog(leaveMsg)
			delete(clients, client.conn)
		}
	}
}
