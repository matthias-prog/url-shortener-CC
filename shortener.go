package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ShortLink struct {
	Link string
	id   int
}

func main() {
	fmt.Println("Hello World")
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/createlink", linkCreator)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("createlink.html")
	t.Execute(w, nil)
}

func linkCreator(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("createlink.html")
	link := ShortLink{"https://google.de", 1}
	t.Execute(w, link)
	fmt.Println("Template wurde ausgeführt")
}
