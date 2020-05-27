package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	var db DatabaseHandler
	db.database = database
	//Creates a new table users with keys id, firstname, lastname, and passHash
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, passHash INT)")
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

//Handles all requests receieved by the webserver
func (db *DatabaseHandler) handler(w http.ResponseWriter, r *http.Request) {
	//If we have a GET request, serve a file in the server-root directory
	fmt.Printf("Recieved connection from %s\n", r.RemoteAddr)
	if r.Method == "GET" {
		//TODO: check for GET/data calls to access db
		//TODO: check for GET/api calls to login
		getFileRequestHandler(w, r)
	} else if r.Method == "POST" {
		db.postHandler(w, r)
	} else {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	}
}

/**
* Handles all GET requests for files.
* TODO: Check for JWT token for access to private dir
**/
func getFileRequestHandler(w http.ResponseWriter, r *http.Request) {
	//check for illegal directory access (../)
	if strings.Contains(r.URL.Path, "../") {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	//check if we are trying to access a file in the private directory
	if strings.Contains(r.URL.Path, "private/") {
		//TODO: JWT Token authentication
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	//Send file contents, if it exists
	filePath := "server-root/" + r.URL.Path //all files are stored in server-root dir
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", fileContents)

}

//check if database holds matching credentials, and if so, hand out a JWT
func (db *DatabaseHandler) postHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	username := r.Form["username"]
	password := r.Form["password"]
	var queryString string
	queryString = fmt.Sprintf(queryString, "SELECT * FROM users WHERE username = %s AND password = %s", username, password)
	rows, err := db.database.Query(queryString)
	if err != nil {
		http.Error(w, "Internal Server Error: Database Query Failed", http.StatusInternalServerError)
	}
	columns, _ := rows.Columns()
	//check for any matches for that user
	if len(columns) == 0 {
		//If that user does not exist, redirect to the homepage
		http.Redirect(w, r, "index.html", http.StatusSeeOther)
	}

	//TODO: Generate a JWT and send to user
	//Redirect to the private homepage
	http.Redirect(w, r, "/private/privatehome.html", http.StatusSeeOther)
}
