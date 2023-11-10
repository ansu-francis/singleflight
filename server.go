package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/singleflight"
	"log"
	"net/http"
	"os"
)

var (
	db *sql.DB
	g  singleflight.Group
)

// Intialize and populate sqliteDB user table with dummy data
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "user.db")
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat("user.db"); err == nil {
		fmt.Println("DB present already. Skipping...")
	} else if errors.Is(err, os.ErrNotExist) {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT NOT NULL PRIMARY KEY, name TEXT)")
		if err != nil {
			panic(err)
		}
		populateDB()
	} else {
		panic(err)
	}
}

// Populate user DB with dummy data
func populateDB() {
	for i := 0; i < 100; i++ {
		idStr := fmt.Sprint(i)
		_, err := db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", idStr, "user-"+idStr)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Populated DB with data records")
}

func main() {
	initDB()

	http.HandleFunc("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		result, err := processRequest(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, result)
	})

	http.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		response, err, _ := g.Do(id, func() (interface{}, error) {
			result, err := processRequest(id)
			if err != nil {
				return nil, err
			}
			return result, nil
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, response)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// Perform DB query
func processRequest(id string) (string, error) {
	fmt.Printf("[DEBUG] Processing request for id %s..\n", id)
	var name string
	row := db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	err := row.Scan(&name)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return name, nil
}
