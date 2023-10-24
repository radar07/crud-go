package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var p product

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Fatalln("Invalid request payload!")
	}

	defer r.Body.Close()
}

func init() {
	err := godotenv.Load()
	fatalErr(err)
}

func initDB() {
	var (
		host     = os.Getenv("DB_HOST")
		port     = 5432
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	fatalErr(err)

	res, err := db.Query("select * from products")
	fatalErr(err)
	defer db.Close()

	for res.Next() {
		var (
			id    int
			name  string
			price float64
		)
		err := res.Scan(&id, &name, &price)
		fatalErr(err)
		fmt.Println(id, name, price)
	}
}

func main() {
	initDB()
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/products", createHandler).Methods("POST")

	// log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("Let's go!!")
}

func fatalErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
