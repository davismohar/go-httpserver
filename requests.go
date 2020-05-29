package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Handles all requests receieved by the webserver
func (db *DatabaseHandler) handler(w http.ResponseWriter, r *http.Request) {
	//If we have a GET request, serve a file in the server-root directory
	fmt.Printf("Recieved connection from %s\n", r.RemoteAddr)
	if r.Method == "GET" {
		//TODO: check for GET/data calls to access db
		//TODO: check for GET/api calls to login
		getFileRequestHandler(w, r)
	} else if r.Method == "POST" {
		if r.URL.Path == "/api/login" {
			db.postLoginHandler(w, r)
		} else if r.URL.Path == "/api/createAccount" {
			db.postCreateAccountHandler(w, r)
		}
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
