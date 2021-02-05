package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var cookie http.Cookie

func setCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "user_name", Value: "gaurang", Expires: time.Now().Add(60 * time.Minute), HttpOnly: true}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "Cookie set: "+cookie.Name)
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/setCookie", setCookie)
	router.HandleFunc("/getCookie", getCookie)
	http.ListenAndServe(":8080", router)
}
