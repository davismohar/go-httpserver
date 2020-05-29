package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

//DatabaseHandler is used to pass the reference to the open database to the http handle functions
type DatabaseHandler struct {
	database *sql.DB
}

func main() {
	//open database
	database, err := sql.Open("sqlite3", "./httpserver.db")
	defer database.Close()
	db := new(DatabaseHandler)
	db.database = database
	//Creates a new table users with keys id, firstname, lastname, and passHash
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT)")
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error opening database")
	}
	fmt.Printf("Opened database\n")
	//create server, it will respond to any request
	fmt.Printf("Starting server on port: 19212\n")
	http.HandleFunc("/", db.handler)
	log.Fatal(http.ListenAndServe(":19212", nil))
}
