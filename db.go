package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	connStr := "host=10.172.156.137 port=5432 user=postgres password=123123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
func updatePassword(db *sql.DB, username, newPassword string) error {
	_, err := db.Exec(`UPDATE users SET password = $1 WHERE username = $2`, newPassword, username)
	return err
}
