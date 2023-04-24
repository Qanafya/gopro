package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type accounts struct {
	id       int
	users    string
	email    string
	password string
}
type comments struct {
	id         int
	product_id int
	name       string
	comment    string
}
type Laptop struct {
	Id    int
	Name  string
	Star  float32
	Price float32
	Photo string
}
type User struct {
	id       int
	users    string
	email    string
	password string
}
type Comment struct {
	id         int
	product_id int
	name       string
	comment    string
}
type Demo struct {
	id     int
	price  int
	amount int
}

var e int
var ee int
var db *sql.DB

func main() {
	r := mux.NewRouter()
	fmt.Println("started")
	ee = 0
	db, err := sql.Open("mysql", "root:EROMA35292@tcp(localhost:3306)/first")
	if err != nil {
		fmt.Println("error at the connecting db")
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("connected")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "template.html", nil)
	})

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM laptops")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		laptops := make([]Laptop, 0)
		for rows.Next() {
			var laptop Laptop
			err := rows.Scan(&laptop.Id, &laptop.Name, &laptop.Star, &laptop.Price, &laptop.Photo)
			if err != nil {
				log.Fatal(err)
			}
			laptops = append(laptops, laptop)
		}
		t, err := template.ParseFiles("templates/products.html")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(laptops)
		t.Execute(w, laptops)
	})
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		//id := r.FormValue("id")
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

		d, _ := db.Query("select * from users")
		defer d.Close()
		var count []accounts
		for d.Next() {
			var x accounts
			err = d.Scan(&x.id, &x.users, &x.email, &x.password)
			if err != nil {
				panic(err)
			}
			count = append(count, x)
		}
		/////
		id := len(count) + 1

		fmt.Println("e equals: ", e, "  length: ", len(account))
		if e%2 == 1 && len(account) == 0 {
			fmt.Println("inserting id: ", id, " username: ", users, " email: ", email, " password: ", password)
			rows, _ := db.Query("insert into users(id, users, email, password) values (?, ?, ?, ?)", id, users, email, password)
			defer rows.Close()
			tpl.ExecuteTemplate(w, "template.html", "Successfully registered")
		} else {
			tpl.ExecuteTemplate(w, "template.html", "Username already taken")
		}
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
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
	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<table>")
		for _, d := range account {
			hashn := len(d.password)
			hash := strings.Repeat("*", hashn)

			fmt.Fprintf(w, "<tr>id: %d<br>username: %s<br>email: %s<br>password: %s<br><br></tr>", d.id, d.users, d.email, hash)
		}
		fmt.Fprint(w, "</table>")
	})
	r.HandleFunc("/commed/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inserting values")
		product_id := r.FormValue("to_id")
		name := r.FormValue("name")
		comment := r.FormValue("comment")
		ee = ee + 1
		d, _ := db.Query("select * from comments")
		defer d.Close()
		var count []comments
		for d.Next() {
			var x comments
			err = d.Scan(&x.id, &x.product_id, &x.name, &x.comment)
			if err != nil {
				panic(err)
			}
			count = append(count, x)
		}
		id := len(count) + 1

		if ee%2 == 1 {
			fmt.Println("e equals: ", e, "  length: ", len(comment))
			fmt.Println("inserting id: ", id, " product_id: ", product_id, " name: ", name, " comment: ", comment)
			rows, _ := db.Query("insert into comments(id, product_id, name, comment) values (?, ?, ?, ?)", id, product_id, name, comment)
			defer rows.Close()
		}
		tpl.ExecuteTemplate(w, "template.html", "Successfully commented")
	})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM laptops")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		laptops := make([]Laptop, 0)
		for rows.Next() {
			var laptop Laptop
			err := rows.Scan(&laptop.Id, &laptop.Name, &laptop.Star, &laptop.Price, &laptop.Photo)
			if err != nil {
				log.Fatal(err)
			}
			laptops = append(laptops, laptop)
		}
		_, err = template.ParseFiles("templates/products.html")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(laptops)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<table>")
		for _, d := range laptops {
			l := "http://localhost:8000/detail/" + strconv.Itoa(d.Id)
			//l := "/detail/1"
			//img := d.photo
			fmt.Fprintf(w, "<tr>id: %d<br>Model: <a href=%s>%s</a><br>Price: %f<br>Rating: %f<br>Photo src: <img src=%s></img><br><br></tr>", d.Id, l, d.Name, d.Price, d.Star, d.Photo)
		}
		fmt.Fprint(w, "</table>")

		/*rows, err := db.Query("SELECT * FROM demo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		demos := []Demo{}
		for rows.Next() {
			var demo Demo
			err := rows.Scan(&demo.id, &demo.price, &demo.amount)
			if err != nil {
				log.Fatal(err)
			}
			demos = append(demos, demo)
		}
		tpl.ExecuteTemplate(w, "test.html", demos)
		fmt.Println(demos)*/
	})
	r.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		laptop1 := Laptop{
			Id:    1,
			Name:  "HP",
			Star:  2,
			Price: 800,
			Photo: "asd",
		}
		laptop2 := Laptop{
			Id:    2,
			Name:  "Dell",
			Star:  3,
			Price: 200,
			Photo: "asd",
		}
		laptop := []Laptop{
			laptop1,
			laptop2,
		}
		tpl.ExecuteTemplate(w, "test2.html", laptop)
	})

	r.HandleFunc("/detail/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		rows, err := db.Query("SELECT * FROM laptops where id=?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		laptops := make([]Laptop, 0)
		for rows.Next() {
			var laptop Laptop
			err := rows.Scan(&laptop.Id, &laptop.Name, &laptop.Star, &laptop.Price, &laptop.Photo)
			if err != nil {
				log.Fatal(err)
			}
			laptops = append(laptops, laptop)
		}
		fmt.Println(laptops)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<table>")
		for _, d := range laptops {
			l := "/detail/" + strconv.Itoa(d.Id)
			//l := "/detail/1"
			fmt.Fprintf(w, "<tr>id: %d<br>Model: <a href=%s>%s</a><br>Price: %f<br>Rating: %f<br>Photo src: <img src=%s></img><br><br></tr>", d.Id, l, d.Name, d.Price, d.Star, d.Photo)
		}
		fmt.Fprint(w, "</table>")
	})

	r.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "product-detail.html", nil)
	})
	log.Fatal(http.ListenAndServe(":8000", r))

}
