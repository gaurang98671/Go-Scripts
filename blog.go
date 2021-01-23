package main

import (
	"database/sql"
	"fmt"
	"strconv"

	//"fmt"
	"html/template"
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
	ID         string
	Title      string
	Content    template.HTML
	Comments   []Comment
	RawContent string
	Date       string
	GUIDE      string
}

type Comment struct {
	Id      int
	Name    string
	Email   string
	Comment string
}

type JSONResponse struct {
	Fields map[string]string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title, page_content, page_date, page_guide FROM blog ORDER BY page_date DESC")

	if err != nil {
		log.Println("Unable to fetch records")
		log.Println(err)
	}
	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUIDE)

		thisPage.RawContent = truncate(thisPage)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}

	t, _ := template.ParseFiles("templates/home.html")
	t.Execute(w, Pages)

	//fmt.Fprint(w, Pages)

}

func getPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guide := vars["guide"]
	thisPage := Page{}
	log.Println(guide)
	err := database.QueryRow("SELECT id,page_title, page_content, page_date FROM blog WHERE page_guide=?", guide).Scan(&thisPage.ID, &thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	log.Println("Page id is: ", thisPage.ID)
	if err != nil {
		log.Println("Page not found")
		log.Println(err)
		http.ServeFile(w, r, "404.html")
	}
	thisPage.GUIDE = guide
	thisPage.Content = template.HTML(thisPage.RawContent)

	comments, err := database.Query("SELECT name, email, comment from comments where page_id=?", thisPage.ID)
	if err != nil {
		log.Println(err)
	}

	for comments.Next() {
		var thisComment Comment
		comments.Scan(&thisComment.Name, &thisComment.Email, &thisComment.Comment)

		thisPage.Comments = append(thisPage.Comments, thisComment)
	}

	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func APICommentsPost(w http.ResponseWriter, r *http.Request) {

	var apiCommentAdded bool
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error)
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	comment := r.FormValue("comment")
	guide := r.FormValue("guide")
	page_id := r.FormValue("page_id")
	log.Println("Name:", name, "Email:", email, "Comment", comment, "guide", guide, "Page id: ", page_id)
	res, err := database.Exec("INSERT INTO comments(page_id, name, email, comment, comment_guid) VALUES(?,?,?,?,?)", page_id, name, email, comment, guide)

	if err != nil {
		log.Println("Failed to insert")
		log.Println(err)
		fmt.Fprint(w, "Failed to insert comment")
	} else {
		log.Println("Inserted")
		id, err := res.LastInsertId()
		log.Println(id)
		log.Println(err)
		if err != nil {
			apiCommentAdded = false
		} else {
			apiCommentAdded = true
		}

		//var resp JSONResponse

		commentAddedBool := strconv.FormatBool(apiCommentAdded)
		log.Println(commentAddedBool)
		var m map[string]string
		m = make(map[string]string)
		m["id"] = string(id)
		m["added"] = commentAddedBool
		log.Println(m)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, m)

	}

}

func truncate(page Page) string {
	if len(page.RawContent) > 150 {
		content := page.Content
		truncated_content := []rune(content)
		new_truncated_content := string(truncated_content[0:150]) + "...."
		return new_truncated_content
	}
	return page.RawContent
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
	router.HandleFunc("/page/{guide}", getPage)
	router.HandleFunc("/home", serveHome)
	router.HandleFunc("/api/comments", APICommentsPost)
	http.ListenAndServe(Port, router)
}
