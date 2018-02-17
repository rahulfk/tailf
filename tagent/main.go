package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"os"
	"github.com/fsnotify/fsnotify"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func tail(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("path read error:", err)
		return
	}
	filename := string(message)
	log.Printf("path : %s", filename)


	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		file, err := os.Open(filename)
		if err != nil {
			log.Println("file open error : ", err)
			c.WriteMessage(websocket.TextMessage, []byte("file not found"))
			done <- true
			return
		}

		fileInfo, err := os.Stat(filename)
		offset := fileInfo.Size() - 1
		file.Close()

		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Write == fsnotify.Write {
					f, err := os.Open(filename)
					if err != nil {
						log.Println("file open error : ", err)
						c.WriteMessage(websocket.TextMessage, []byte("file not found"))
						f.Close()
						done <- true
						return
					}
					buffer := make([]byte, 1024)
					numRead, err := f.ReadAt(buffer, offset)
					for numRead > 0 {
						data := buffer[:numRead]
						err = c.WriteMessage(websocket.TextMessage, data)
						if err != nil {
							f.Close()
							done <- true
							return
						}
						offset = offset + int64(numRead)
						buffer := make([]byte, 1024)
						numRead, err = f.ReadAt(buffer, offset)
					}
					f.Close()
				}
			case <-watcher.Errors:
				done <- true
				return
			}
		}


	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	select {
	case <-done:
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/tail", tail)
	log.Fatal(http.ListenAndServe(*addr, nil))
}