package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//check if database holds matching credentials, and if so, hand out a JWT
func (db *DatabaseHandler) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	usernameSlice := r.Form["username"]
	passwordSlice := r.Form["password"]
	username := usernameSlice[0]
	password := passwordSlice[0]
	//generate hash from password
	passwordHash, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error: Hash Failed", http.StatusInternalServerError)
	}
	//query database for the given information
	fmt.Printf("Querying for user %v with password hash %v\n", username, passwordHash)
	// var sb strings.Builder
	// sb.WriteString("SELECT * FROM users WHERE username = \'")
	// sb.WriteString(username)
	queryString := fmt.Sprintf("SELECT * FROM users WHERE username = '%v' AND password = '%v'", username, passwordHash)
	rows, err := db.database.Query(queryString)
	if err != nil {
		http.Error(w, "Internal Server Error: Database Query Failed", http.StatusInternalServerError)
	}
	columns, _ := rows.Columns()
	//check for any matches for that user
	if len(columns) > 2 {
		//If that user does not exist, redirect to the homepage
		fmt.Printf("Account could not be found for user: %s\n", username)
		http.Redirect(w, r, "index.html", http.StatusSeeOther)
	}

	//TODO: Generate a JWT and send to user
	//Redirect to the private homepage
	http.Redirect(w, r, "/private/privatehome.html", http.StatusSeeOther)
}

//adds a new user to the database, if it does not already exist
func (db *DatabaseHandler) postCreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	usernameSlice := r.Form["username"]
	passwordSlice := r.Form["password"]
	username := usernameSlice[0]
	password := passwordSlice[0]
	//generate hash from password
	passwordHash, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error: Hash Failed", http.StatusInternalServerError)
	}
	fmt.Printf("Adding user %v with password %v", username, passwordHash)

	stmt, err := db.database.Prepare("INSERT INTO users (username, password) VALUES  ( ?, ? )")
	if err != nil {
		http.Error(w, "Internal Server Error: SQL Statement Preparation Failed", http.StatusInternalServerError)
	}

	_, err = stmt.Exec(username, passwordHash)
	if err != nil {
		if err.Code == 19 {
			http.Error(w, "Internal Server Error: Account Already Exists", http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal Server Error: SQL Exec Failed", http.StatusInternalServerError)
		}
	}

	http.Redirect(w, r, "login.html", http.StatusSeeOther)
}

//Returns the string value of the hash of the given password
//if a bcrypt error occurs, that error is returned
func hashPassword(str string) (string, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(str), 4)
	if err != nil {
		return "", err
	}
	return string(passwordHashBytes), nil
}
