package main

import (
	"flag"
	"os"
	"os/signal"
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var path = flag.String("path", "/Users/rahul/Documents/test.txt", "path of file to tail")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/tail"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan bool);

	go func() {
		err = c.WriteMessage(websocket.TextMessage, []byte(*path))
		if err != nil {
			done <- true
			log.Println("error : ", err)
			return
		}
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				done <- true
			}
			switch messageType {
			case websocket.TextMessage:
				fmt.Printf("%s", message)
			}
		}
	}()

	select {
	case <-interrupt:
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	case <- done:
	}

}
