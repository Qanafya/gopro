package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"reflect"
	"strings"
)

var tpl *template.Template

//var db *sql.DB
//var err error

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type accounts struct {
	id       int
	users    string
	email    string
	password string
}

var e int

func main() {
	fmt.Println("started")
	db, err := sql.Open("mysql", "root:EROMA35292@tcp(localhost:3306)/first")
	if err != nil {
		fmt.Println("error at the connecting db")
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("connected")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "template.html", nil)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		users := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		s, _ := db.Query("select * from users where users = ?", users)
		e = e + 1
		defer s.Close()
		var account []accounts
		for s.Next() {
			var a accounts
			err = s.Scan(&a.id, &a.users, &a.email, &a.password)
			if err != nil {
				panic(err)
			}
			account = append(account, a)
		}
		/////
		fmt.Println("e equals: ", e, "  length: ", len(account))
		if e%2 == 1 && len(account) == 0 {
			fmt.Println("inserting id: ", id, " username: ", users, " email: ", email, " password: ", password)
			rows, _ := db.Query("insert into users(id, users, email, password) values (?, ?, ?, ?)", id, users, email, password)
			defer rows.Close()
		}
		tpl.ExecuteTemplate(w, "template.html", "Successfully registered")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		users := r.FormValue("username")
		password := r.FormValue("password")
		rows, _ := db.Query("select * from users where users = ?", users)
		defer rows.Close()
		var account []accounts
		for rows.Next() {
			var a accounts
			err = rows.Scan(&a.id, &a.users, &a.email, &a.password)
			if err != nil {
				panic(err)
			}
			account = append(account, a)
		}
		if len(account) > 0 {
			fmt.Println(account[0].password, " and ", password)
			if account[0].password == password {
				str := "Welcome to the site, " + users
				tpl.ExecuteTemplate(w, "login.html", str)
			} else {
				tpl.ExecuteTemplate(w, "login.html", "password is not correct")
			}
		} else {
			tpl.ExecuteTemplate(w, "login.html", "user not found")
		}
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		soz := r.FormValue("search")
		fmt.Println("This is soz:", soz)
		rows, _ := db.Query("select * from users where users = ?", soz)
		defer rows.Close()
		var account []accounts
		for rows.Next() {
			var a accounts
			err = rows.Scan(&a.id, &a.users, &a.email, &a.password)
			if err != nil {
				panic(err)
			}
			account = append(account, a)
		}
		fmt.Println("This is rows: ", account)
		//tpl.ExecuteTemplate(w, "template.html", account)

		// Render data as HTML table
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<table>")
		for _, d := range account {
			hashn := len(d.password)
			hash := strings.Repeat("*", hashn)
			fmt.Fprintf(w, "<tr>id: %d<br>username: %s<br>email: %s<br>password: %s<br><br></tr>", d.id, d.users, d.email, hash)
		}
		fmt.Fprint(w, "</table>")
	})
	http.ListenAndServe("localhost:8000", nil)

}
func getField(v *accounts, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.Int())
}
