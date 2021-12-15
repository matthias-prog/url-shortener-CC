package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type ShortLink struct {
	id   int
	Link string
}

var Datenbank *sql.DB

func main() {
	var err error
	Datenbank, err = sql.Open("sqlite3", "./shortener.db")
	checkErr(err)
	//defer Datenbank.Close()

	var tablename string
	sqlInitCheck := "SELECT tbl_name FROM sqlite_master WHERE type='table' AND name='links';"
	err = Datenbank.QueryRow(sqlInitCheck).Scan(&tablename)

	if err == sql.ErrNoRows {
		sqlInit := "CREATE TABLE links (id INTEGER not null primary key autoincrement, link TEXT);"
		_, err = Datenbank.Exec(sqlInit)
		checkErr(err)
		sqlIndex := "CREATE UNIQUE INDEX idx_links_id ON links(id);"
		_, err = Datenbank.Exec(sqlIndex)
		checkErr(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/createlink", linkCreator)
	r.HandleFunc("/{shorturl:[A-Za-z0-9]+}", rerouteHandler)
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
	longurl := strings.Join(r.Form["longurl"], "")
	fmt.Printf("übereichte lange URL: %v\n", longurl)

	sqlInsert, err := Datenbank.Prepare("INSERT INTO links (link) VALUES (?);")
	checkErr(err)
	result, err := sqlInsert.Exec(longurl)
	checkErr(err)

	returnedId, err := result.LastInsertId()
	checkErr(err)
	fmt.Printf("Zurückgegebene ID für die URL %v : %v\n", longurl, returnedId)

	fertigerLink := fmt.Sprintf("%v/%v", r.Host, returnedId)
	fmt.Printf("Der fertige Link %v\n", fertigerLink)
	fmt.Printf("Der Hostname: %v \n", r.Host)
	t, _ := template.ParseFiles("createlink.gohtml")
	link := ShortLink{int(returnedId), fertigerLink}
	t.Execute(w, link)
	fmt.Println("Template wurde ausgeführt")
}

func rerouteHandler(w http.ResponseWriter, r *http.Request) {

	path := mux.Vars(r)
	sqlQuery, err := Datenbank.Prepare("SELECT link FROM links WHERE id = ?;")
	checkErr(err)

	var result string
	err = sqlQuery.QueryRow(path["shorturl"]).Scan(&result)
	checkErr(err)

	http.Redirect(w, r, result, 302)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
