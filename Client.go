package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
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

		// adds readable timestamp to each message
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
