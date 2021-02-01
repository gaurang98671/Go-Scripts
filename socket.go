package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var m = make(map[*websocket.Conn]bool)

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(p)
		if string(p) == "hello" {
			conn.WriteMessage(messageType, []byte("hi"))
		}
		for c, _ := range m {

			c.WriteMessage(messageType, p)
		}

	}

}

func serveSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected")
	m[ws] = true
	fmt.Println(len(m))

	reader(ws)

}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws", serveSocket)
	http.ListenAndServe(":8080", router)
}
