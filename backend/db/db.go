package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	psqlInfo := os.Getenv("DB_URL")
	log.Println(psqlInfo)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	createTable(db)
	// seedTodos()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database!!!")
}

func GetDB() *sql.DB {
	return db
}

func createTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        completed BOOLEAN NOT NULL DEFAULT FALSE
    );`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("eeeerrrroooorrrr")
		log.Fatal(err)
	}
	fmt.Println("Table created successfully")
}

func seedTodos() {
	query := `INSERT INTO todos VALUES (1, 'First', 'false')`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Seed successfully")
}
