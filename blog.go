package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

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
	http.Redirect(w, r, "/home", 301)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title, page_content, page_date FROM blog ORDER BY page_date DESC")

	if err != nil {
		log.Println("Unable to fetch records")
		log.Println(err)
	}

	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
		Pages = append(Pages, thisPage)
	}

	fmt.Fprint(w, Pages)

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
		http.ServeFile(w, r, "404.html")
	}

	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func main() {

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
	router.HandleFunc("/home", serveHome)
	http.ListenAndServe(Port, router)
}
