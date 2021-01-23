package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func setCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "user_name", Value: "gaurang", Expires: time.Now().Add(60 * time.Second), HttpOnly: true}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "Cookie set: "+cookie.Name)
}
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/setCookie", setCookie)
	http.ListenAndServe(":8080", router)
}
