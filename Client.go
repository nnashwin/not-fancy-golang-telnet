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
	blockedUsers map[string]int
	currentRm    string
}

func (c *Client) ReadLines(ch chan<- Message, commCCh chan<- Command) {
	buffc := bufio.NewReader(c.conn)
	for {
		line, err := buffc.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}

		if line[:1] == "/" {
			fmt.Println("Command and not message\n")
			comm := Command{sender: c, input: line[1:]}
			commCCh <- comm
		} else {
			// adds readable timestamp to each message
			t := time.Now().Format(time.RFC822)
			fmt.Printf("%+v %+v: %+v", t, c.userId, line)
			message := Message{sender: c.userId, text: fmt.Sprintf("%+v %+v: %+v", t, c.userId, line)}
			ch <- message
		}

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
