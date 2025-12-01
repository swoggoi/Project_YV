package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=123123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
