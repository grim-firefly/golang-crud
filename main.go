package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

var temp *template.Template
var err error

// connecting to database
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	return db
}

type User struct {
	Id   int
	Name string
}

// index page
func index(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error())
	}
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, user)
	}
	// http.ServeFile(w, r, "static/index.html")
	temp.ExecuteTemplate(w, "index.html", users)
}

// create page
func create(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/create.html")

}

// create page
func edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := dbConn()
	var name string
	err := db.QueryRow("SELECT * FROM user WHERE id=?", id).Scan(&id, &name)
	if err != nil {
		panic(err)
	}
	fmt.Print(name, id)

}

// store
func store(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Please enter name", http.StatusBadRequest)
	}
	db := dbConn()
	_, err := db.Query("INSERT INTO user(name) VALUES(?)", name)

	if err != nil {
		panic(err)
	}

	db.Close()
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// delete function
func destroy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Invalid Action", http.StatusBadRequest)
	}
	db := dbConn()
	_, err := db.Query("DELETE FROM user WHERE id=?", id)

	if err != nil {
		panic(err)
	}

	db.Close()
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// main function

func main() {

	temp, err = template.ParseGlob("static/*.html")
	if err != nil {
		panic(err)
	}
	fmt.Println("Hello, World!")
	router := chi.NewRouter()
	router.Get("/", index)
	router.Get("/create", create)
	router.Post("/", store)
	router.Post("/delete", destroy)
	router.Get("/edit/{id}", edit)

	http.ListenAndServe(":3000", router)
}
