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

type Laptop struct {
	Id    int
	Name  string
	Star  float32
	Price float32
	Photo string
}
type Detail struct {
	Id    int
	Name  string
	Desc  string
	Foto1 string
	Foto2 string
	Foto3 string
	Foto4 string
	Price int
}
type Detailf struct {
	Rating float32
	Id     int
	Name   string
	Desc   string
	Price  int
	Foto1  string
	Foto2  string
	Foto3  string
	Foto4  string
}
type Detaila struct {
	Amount int
	Id     int
	Name   string
	Desc   string
	Price  int
	Foto1  string
	Foto2  string
	Foto3  string
	Foto4  string
}
type Comment struct {
	Id         int
	Product_id int
	Name       string
	Comment    string
}
type Rating struct {
	Id         int
	Product_id int
	Rating     int
}
type Rat struct {
	Rating float32
}
type Purchase struct {
	Id         int
	Product_id int
	Amount     int
	Active     int
}

var e int
var ee int
var e1 int
var db *sql.DB

func main() {
	r := mux.NewRouter()
	fmt.Println("started")
	ee = 0
	e1 = 0
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
		rows, err := db.Query("select a.rat, d.Id, d.Name, d.Desc, d.Foto1, d.Foto2, d.Foto3, d.Foto4, d.Price from detail d join (select product_id, avg(rating) as rat from rating group by product_id) a on d.id = a.product_id;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		details := make([]Detailf, 0)
		for rows.Next() {
			var detail Detailf
			err := rows.Scan(&detail.Rating, &detail.Id, &detail.Name, &detail.Desc, &detail.Foto1, &detail.Foto2, &detail.Foto3, &detail.Foto4, &detail.Price)
			if err != nil {
				log.Fatal(err)
			}
			details = append(details, detail)
		}
		t, err := template.ParseFiles("templates/products.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, details)
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
		if e%2 == 1 {
			if len(account) == 0 {
				fmt.Println("inserting id: ", id, " username: ", users, " email: ", email, " password: ", password)
				rows, _ := db.Query("insert into users(id, users, email, password) values (?, ?, ?, ?)", id, users, email, password)
				defer rows.Close()
				tpl.ExecuteTemplate(w, "template.html", "Successfully registered")
			} else {
				tpl.ExecuteTemplate(w, "template.html", "Username already taken")
			}
		} else {
			tpl.ExecuteTemplate(w, "template.html", nil)
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
		var u = "http://localhost:8000/products"

		if len(account) > 0 {
			fmt.Println(account[0].password, " and ", password)
			if account[0].password == password {
				http.Redirect(w, r, u, http.StatusSeeOther)
			} else {
				tpl.ExecuteTemplate(w, "template.html", "password is not correct")
			}
		} else {
			tpl.ExecuteTemplate(w, "template.html", "user not found")
		}
	})
	r.HandleFunc("/filter/", func(w http.ResponseWriter, r *http.Request) {
		price_min := r.FormValue("price")
		price_max := r.FormValue("price2")
		rating_min := r.FormValue("rating")
		rating_max := r.FormValue("rating2")
		rows, err := db.Query("select a.rat, d.Id, d.Name, d.Desc, d.Foto1, d.Foto2, d.Foto3, d.Foto4, d.Price from detail d join (select product_id, avg(rating) as rat from rating group by product_id) a on d.id = a.product_id where a.rat > ? and a.rat < ? and d.Price > ? and d.Price < ?;", rating_min, rating_max, price_min, price_max)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		details := make([]Detailf, 0)
		for rows.Next() {
			var detail Detailf
			err := rows.Scan(&detail.Rating, &detail.Id, &detail.Name, &detail.Desc, &detail.Foto1, &detail.Foto2, &detail.Foto3, &detail.Foto4, &detail.Price)
			if err != nil {
				log.Fatal(err)
			}
			details = append(details, detail)
		}
		t, err := template.ParseFiles("templates/products.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, details)
	})
	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		soz := r.FormValue("search")
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

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<table>")
		for _, d := range account {
			hashn := len(d.password)
			hash := strings.Repeat("*", hashn)

			fmt.Fprintf(w, "<tr>id: %d<br>username: %s<br>email: %s<br>password: %s<br><br></tr>", d.id, d.users, d.email, hash)
		}
		fmt.Fprint(w, "</table>")
	})
	r.HandleFunc("/rate/", func(w http.ResponseWriter, r *http.Request) {
		product_id := r.FormValue("product_id")
		rating := r.FormValue("rating")
		e1 = e1 + 1
		d, _ := db.Query("select * from rating")
		defer d.Close()
		var count []Rating
		for d.Next() {
			var x Rating
			err = d.Scan(&x.Id, &x.Product_id, &x.Rating)
			if err != nil {
				panic(err)
			}
			count = append(count, x)
		}
		id := len(count) + 1

		//if e1%2 == 1 {
		rows, _ := db.Query("insert into rating(Id, Product_id, Rating) values (?, ?, ?)", id, product_id, rating)
		defer rows.Close()
		//}
		var u = "http://localhost:8000/detail/" + product_id
		http.Redirect(w, r, u, http.StatusSeeOther)
	})
	r.HandleFunc("/commed/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inserting values")
		product_id := r.FormValue("to_id")
		name := r.FormValue("name")
		comment := r.FormValue("comment")
		ee = ee + 1
		d, _ := db.Query("select * from comments")
		defer d.Close()
		var count []Comment
		for d.Next() {
			var x Comment
			err = d.Scan(&x.Id, &x.Product_id, &x.Name, &x.Comment)
			if err != nil {
				panic(err)
			}
			count = append(count, x)
		}
		id := len(count) + 1

		rows, _ := db.Query("insert into comments(id, product_id, name, comment) values (?, ?, ?, ?)", id, product_id, name, comment)
		defer rows.Close()
		var u = "http://localhost:8000/detail/" + product_id
		http.Redirect(w, r, u, http.StatusSeeOther)

		/*rows, err := db.Query("SELECT * FROM detail")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		details := make([]Detail, 0)
		for rows.Next() {
			var detail Detail
			err := rows.Scan(&detail.Id, &detail.Name, &detail.Desc, &detail.Foto1, &detail.Foto2, &detail.Foto3, &detail.Foto4, &detail.Price)
			if err != nil {
				log.Fatal(err)
			}
			details = append(details, detail)
		}
		t, err := template.ParseFiles("templates/products.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, details)*/
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
	r.HandleFunc("/sell/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "Sell.html", nil)
	})
	r.HandleFunc("/filtering/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "login.html", nil)
	})
	r.HandleFunc("/cart/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("select p.amount, d.Id, d.Name, d.Desc, d.Foto1, d.Foto2, d.Foto3, d.foto4, d.Price from purchase p join detail d on p.product_id=d.Id where p.active=1;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		details := make([]Detaila, 0)
		for rows.Next() {
			var detail Detaila
			err := rows.Scan(&detail.Amount, &detail.Id, &detail.Name, &detail.Desc, &detail.Foto1, &detail.Foto2, &detail.Foto3, &detail.Foto4, &detail.Price)
			if err != nil {
				log.Fatal(err)
			}
			details = append(details, detail)
		}
		t, err := template.ParseFiles("templates/cart.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, details)
	})
	r.HandleFunc("/buy/", func(w http.ResponseWriter, r *http.Request) {
		product_id := r.FormValue("id")
		amount := r.FormValue("amount")
		d, _ := db.Query("select * from purchase")
		defer d.Close()
		count := make([]Purchase, 0)
		for d.Next() {
			var x Purchase
			err = d.Scan(&x.Id, &x.Product_id, &x.Amount, &x.Active)
			if err != nil {
				panic(err)
			}
			count = append(count, x)
		}
		id := len(count) + 1
		row, _ := db.Query("INSERT INTO `first`.`purchase` (`id`, `product_id`, `amount`, `active`) VALUES (?, ?, ?, ?);", id, product_id, amount, 1)
		defer row.Close()

		var u = "http://localhost:8000/cart/"
		http.Redirect(w, r, u, http.StatusSeeOther)
	})
	r.HandleFunc("/selling/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inserting values")
		Name := r.FormValue("name")
		Desc := r.FormValue("desc")
		Foto1 := r.FormValue("foto1")
		Foto2 := r.FormValue("foto2")
		Foto3 := r.FormValue("foto3")
		Foto4 := r.FormValue("foto4")
		Price := r.FormValue("price")
		ee = ee + 1
		d, _ := db.Query("select * from detail")
		defer d.Close()
		count := make([]Detail, 0)
		for d.Next() {
			var x Detail
			err = d.Scan(&x.Id, &x.Name, &x.Desc, &x.Foto1, &x.Foto2, &x.Foto3, &x.Foto4, &x.Price)
			if err != nil {
				panic(err)
			}
			count = append(count, x)
		}
		id := len(count) + 1

		rowd, _ := db.Query("INSERT INTO `first`.`detail` (`Id`, `Name`, `Desc`, `Foto1`, `Foto2`, `Foto3`, `foto4`, `Price`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);", id, Name, Desc, Foto1, Foto2, Foto3, Foto4, Price)
		defer rowd.Close()

		var u = "http://localhost:8000/products"
		http.Redirect(w, r, u, http.StatusSeeOther)
	})
	r.HandleFunc("/dele/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		rows, _ := db.Query("update purchase set active=0 where product_id=?", id)
		defer rows.Close()

		var u = "http://localhost:8000/cart/"
		http.Redirect(w, r, u, http.StatusSeeOther)
	})
	r.HandleFunc("/detail/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		rows, err := db.Query("SELECT * FROM detail where id=?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		details := make([]Detail, 0)
		for rows.Next() {
			var detail Detail
			err := rows.Scan(&detail.Id, &detail.Name, &detail.Desc, &detail.Foto1, &detail.Foto2, &detail.Foto3, &detail.Foto4, &detail.Price)
			if err != nil {
				log.Fatal(err)
			}
			details = append(details, detail)
		}

		rowsc, err := db.Query("SELECT * FROM comments where product_id=?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		comments := make([]Comment, 0)
		for rowsc.Next() {
			var comment Comment
			err := rowsc.Scan(&comment.Id, &comment.Product_id, &comment.Name, &comment.Comment)
			if err != nil {
				log.Fatal(err)
			}
			comments = append(comments, comment)
		}

		rowsr, err := db.Query("SELECT case when avg(rating) is not null then avg(rating) when avg(rating) is null then 0 end as rate FROM rating where product_id=?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		rats := make([]Rat, 0)
		for rowsr.Next() {
			var rat Rat
			err := rowsr.Scan(&rat.Rating)
			if err != nil {
				log.Fatal(err)
			}
			rats = append(rats, rat)
		}

		t, err := template.ParseFiles("templates/product-detail.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, map[string]interface{}{"details": details, "comments": comments, "rats": rats})
	})

	r.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "product-detail.html", nil)
	})
	log.Fatal(http.ListenAndServe(":8000", r))
}
