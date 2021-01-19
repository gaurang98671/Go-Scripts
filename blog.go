package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	DBHost     = "127.0.0.1"
	DBPort     = ":3306"
	DBUser     = "root"
	DBPassword = "1962"
	DBDbase    = "blogs"
	Port       = ":8080"
)

var database *sql.DB

type Page struct {
	Title   string
	Content string
	Date    string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "homePage.html")
}

func getPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	thisPage := Page{}
	log.Println(id)
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM blog WHERE id=?", id).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
	if err != nil {
		log.Println("Page not found")
		log.Println(err)
	}
	fmt.Fprint(w, thisPage.Title)
}

func main() {
	//dbConn := fmt.Sprint("mysql", "root:1962@/blogs")
	db, err := sql.Open("mysql", "root:1962@/blogs")
	if err != nil {
		log.Println("Failed to connect")
		log.Println(err.Error)
	}
	log.Println("Database connected")

	database = db
	router := mux.NewRouter()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/page/{id}", getPage)
	http.ListenAndServe(Port, router)
}
