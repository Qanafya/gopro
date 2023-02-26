package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	_ "https://github.com/go-sql-driver/mysql"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type Welcome struct {
	Sale string
	Time string
}

type User struct {
	Username string
	Password string
	Email    string
}

type Item struct {
	Name        string
	Description string
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

func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Read the form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Create a new user from the form data
			user := &User{
				Username: r.Form.Get("username"),
				Password: r.Form.Get("password"),
				Email:    r.Form.Get("email"),
			}

			// Insert the user into the database
			_, err = db.Exec(`
				INSERT INTO users (username, password, email)
				VALUES (?, ?, ?)
			`, user.Username, user.Password, user.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Render the success message
			tmpl := template.Must(template.ParseFiles("index2.html"))
			err = tmpl.Execute(w, "Registration successful")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// Render the registration form
		tmpl := template.Must(template.ParseFiles("register.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Read the form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Retrieve the user from the database
			username := r.Form.Get("username")
			password := r.Form.Get("password")
			row := db.QueryRow(`
				SELECT username, password FROM users
				WHERE username = ? AND password = ?
			`, username, password)

			// Check if the user exists
			var dbUsername, dbPassword string
			err = row.Scan(&dbUsername, &dbPassword)
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Render the success message
			tmpl := template.Must(template.ParseFiles("index2.html"))
			return
		}
	}
}*/

var users = []User{}

func createDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(255) NOT NULL PRIMARY KEY,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			name VARCHAR(255) NOT NULL PRIMARY KEY,
			description TEXT
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

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
