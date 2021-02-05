package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var cookie http.Cookie

var m = make(map[*websocket.Conn]bool)

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		prefix_string := string(cookie.Value) + " says: "
		prefix := []byte(prefix_string)
		p1 := append(prefix[:], p[:]...)
		for c, _ := range m {

			c.WriteMessage(messageType, p1)
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
	log.Println(cookie.Value)
	http.ServeFile(w, r, "home.html")
	//fmt.Fprint(w, "Home")
}

func serveChat(w http.ResponseWriter, r *http.Request) {
	if cookie.Value == "" {
		http.Redirect(w, r, "/", 302)
	}
	http.ServeFile(w, r, "socket.html")

}

func registerUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("userName")
	cookie.Name = "userName"
	cookie.Value = string(name)
	cookie.Expires = time.Now().Add(60 * time.Second)
	cookie.HttpOnly = true
	http.SetCookie(w, &cookie)
	log.Println("Cookie set")
	http.Redirect(w, r, "/room", 301)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws", serveSocket)
	router.HandleFunc("/room", serveChat)
	router.HandleFunc("/register", registerUser)
	http.ListenAndServe(":8080", router)
}
