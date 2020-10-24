package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db = Conn()
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users(id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, firstname VARCHAR(255) NOT NULL, lastname VARCHAR(255) NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("taisakushimasu"),
		newrelic.ConfigLicense("b9da13af214c7aabced2f1efc831f0bc0ebfNRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	//http.HandleFunc("/users", users)
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/users", users))

	err = http.ListenAndServe(os.Getenv("LISTEN_PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Conn() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("HOST"), os.Getenv("MYSQL_SCHEMA")))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func users(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}
	var users []User
	users = append(users, user)

	stmt, err := db.Prepare("INSERT INTO users(firstname, lastname) VALUES(?, ?)")
	if err != nil {
		log.Print(err)
		return
	}

	_, err = stmt.Exec(user.FirstName, user.LastName)
	if err != nil {
		log.Print(err)
		return
	}

	json.NewEncoder(w).Encode(users)
}
