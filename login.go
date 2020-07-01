package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	//Generate a JWT and send to user
	now := time.Now().Unix()
	expires := now + 3600
	claims := jwtClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now,
			ExpiresAt: expires,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSigningKey)
	if err != nil {
		http.Error(w, "Internal Server Error: Generating JWT Failed", http.StatusInternalServerError)
	}
	var cookie http.Cookie
	cookie.Name = "auth-token"
	cookie.Value = tokenString
	cookie.Path = "/private"

	http.SetCookie(w, &cookie)

	//Redirect to the private homepage
	fmt.Fprintf(w, "<p>Login sucessful.\n<a href=\"/private/privatehome.html\">Private Home</a></p>")
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
	//adds entry into database with their username and password hash
	stmt, err := db.database.Prepare("INSERT INTO users (username, password) VALUES  ( ?, ? )")
	if err != nil {
		http.Error(w, "Internal Server Error: SQL Statement Preparation Failed", http.StatusInternalServerError)
	}

	_, err = stmt.Exec(username, passwordHash)
	if err != nil {
		http.Error(w, "Internal Server Error: Account Already Exists", http.StatusInternalServerError)
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
