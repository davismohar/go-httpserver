package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	//create server, it will respond to any request
	fmt.Printf("Starting server on port: 19212")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":19212", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	//If we have a GET request, serve a file in the server-root directory
	fmt.Printf("Recieved connection from %s\n", r.RemoteAddr)
	if r.Method == "GET" {
		getRequestHandler(w, r)
	} else {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	}
}

func getRequestHandler(w http.ResponseWriter, r *http.Request) {
	//check for illegal directory access (../)
	if strings.Contains(r.URL.Path, "../") {
		http.Error(w, "Forbidden", http.StatusForbidden)
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
