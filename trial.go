package main

import (
	"log"
	"net/http"
	"text/template"
)

type ContactDetails struct {
	Login         string
	Password      string
	Email         string
	Success       bool
	StorageAccess string
}

var (
	tmpl = template.Must(template.ParseFiles("template/template.html"))
)

func handler(w http.ResponseWriter, req *http.Request) {

	if req.FormValue("password") == req.FormValue("verpass") {
		data := ContactDetails{
			Login:    req.FormValue("login"),
			Password: req.FormValue("password"),
			Email:    req.FormValue("email"),
		}
		data.Success = true
		tmpl.Execute(w, data)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":800", nil))
}
