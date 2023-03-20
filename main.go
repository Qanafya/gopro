package main

import (
	"html/template"
	"net/http"
	"users"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("template/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/log-in", logInUser)
	http.HandleFunc("/signup", signUpUser)
	http.HandleFunc("/signsuccess", signsuccess)
	http.HandleFunc("/signfail", signfail)
	http.HandleFunc("/loginsuccess", loginsuccess)
	http.HandleFunc("/loginfail", loginfail)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "template.html", nil)
}
func login(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "template.html", nil)
}
func signsuccess(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "template.html", "signed up successfully")
}
func signfail(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "template.html", "Failed to sign up")
}
func loginsuccess(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "template.html", "logged in successfully")
}
func loginfail(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "template.html", "Failed to login")
}
func getUser(r *http.Request) users.User {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	return users.User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

func signUpUser(w http.ResponseWriter, r *http.Request) {
	newUser := getUser(r)
	err := users.DefaultUserService.CreateUser(newUser)
	if err != nil {
		http.Redirect(w, r, "/signfail", http.StatusSeeOther)
		//t.ExecuteTemplate(w, filename, "User registration failure")
		return
	}
	http.Redirect(w, r, "/signsuccess", http.StatusSeeOther)
	return
}
func logInUser(w http.ResponseWriter, r *http.Request) {
	newUser := getUser(r)
	ok := users.DefaultUserService.VerifyUser(newUser)
	if !ok {
		http.Redirect(w, r, "/loginfail", http.StatusSeeOther)
		//t.ExecuteTemplate(w, filename, "User registration failure")
		return
	}
	http.Redirect(w, r, "/loginsuccess", http.StatusSeeOther)
	return
}
