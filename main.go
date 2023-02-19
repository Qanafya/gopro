package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type Welcome struct {
	Sale string
	Time string
}

func main() {
	welcome := Welcome{"Sale Begins Now", time.Now().Format(time.Stamp)}
	template := template.Must(template.ParseFiles("template/template.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if sale := r.FormValue("sale"); sale != "" {
			welcome.Sale = sale
		}
		if err := template.ExecuteTemplate(w, "template.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8000", nil))
}

type User struct {
	Username string
	Password string
	Email    string
}

var users = []User{}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		u := User{
			Username: username,
			Password: password,
			Email:    email,
		}
		users = append(users, u)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "template.html", nil)
}
