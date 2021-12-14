package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

type ShortLink struct {
	id   int
	Link string
}

var db sql.DB

func main() {
	db, err := sql.Open("sqlite3", "./shortener.db")
	checkErr(err)
	defer db.Close()

	var tablename string
	sqlInitCheck := "SELECT tbl_name FROM sqlite_master WHERE type='table' AND name='links';"
	err = db.QueryRow(sqlInitCheck).Scan(&tablename)
	checkErr(err)

	if err == sql.ErrNoRows {
		sqlInit := "CREATE TABLE links (id TEXT not null primary key, link TEXT);"
		_, err = db.Exec(sqlInit)
		checkErr(err)
		sqlIndex := "CREATE UNIQUE INDEX idx_links_id ON links(id);"
		_, err = db.Exec(sqlIndex)
		checkErr(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/createlink", linkCreator)
	r.HandleFunc("/{shorturl}", rerouteHandler)
	r.HandleFunc("/", pageHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}
	log.Fatal(srv.ListenAndServe())
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("createlink.gohtml")
	t.Execute(w, nil)
}

func linkCreator(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	longurl := r.Form["longurl"]
	fmt.Printf("übereichte lange URL: %v\n", longurl)

	sqlInsert, err := db.Prepare("INSERT INTO links (link) VALUES (?);")
	checkErr(err)
	result, err := sqlInsert.Exec(longurl)
	checkErr(err)

	returnedId, err := result.LastInsertId()
	checkErr(err)
	fmt.Printf("Zurückgegebene ID für die URL %v : %v\n", longurl, returnedId)

	fertigerLink := fmt.Sprintf("%v/%v", r.URL, returnedId)
	t, _ := template.ParseFiles("createlink.gohtml")
	link := ShortLink{int(returnedId), fertigerLink}
	t.Execute(w, link)
	fmt.Println("Template wurde ausgeführt")
}

func rerouteHandler(w http.ResponseWriter, r *http.Request) {

	path := mux.Vars(r)
	sqlQuery, err := db.Prepare("SELECT link FROM links WHERE id = ?;")
	checkErr(err)

	result, err := sqlQuery.Query(path["shorturl"])
	checkErr(err)

	var dbResult string
	err = result.Scan(&dbResult)
	checkErr(err)
	result.Close()

	http.Redirect(w, r, dbResult, 302)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
