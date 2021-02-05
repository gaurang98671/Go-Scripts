package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func serveRoot(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user_name")
	fmt.Fprint(w, cookie)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", serveRoot)
	http.ListenAndServe(":8080", router)
}
