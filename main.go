package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func users(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}
	var users []User
	users = append(users, user)
	json.NewEncoder(w).Encode(users)
}
func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/sandbox-docker-compose-mysql")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users(id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, firstname VARCHAR(255) NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/users", users)
	err = http.ListenAndServe(":8002", nil)
	if err != nil {
		log.Fatal(err)
	}
}
