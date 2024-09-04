package main

import (
	"dash-iot/dash-iot/auth"
	"database/sql"
	"io"
	"log"
	"net/http"
	"text/template"

	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	url := "postgres://username:password@localhost:5432/database_name"

	db, err := sql.Open("pgx", url)

	if err != nil {
		log.Fatal(err)
	}

	tmpl := make(map[string]*template.Template)

	tmpl["login"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/login.html"))

	mux := http.NewServeMux()

	mux.HandleFunc("/index", auth.AuthHandler(func(w http.ResponseWriter, r *http.Request, s auth.Session) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Hello world 2")
	}))
	mux.HandleFunc("/login", auth.Login(db, func(w http.ResponseWriter, r *http.Request) {
		tmpl["login"].Execute(w, nil)
	}))

	fileHandler := http.FileServer(http.Dir("./style"))

	mux.Handle("/style/", http.StripPrefix("/style", fileHandler))

	err = http.ListenAndServe("127.0.0.1:8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
