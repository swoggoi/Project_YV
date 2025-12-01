package main

import (
	"database/sql"
	"time"
)

type Message struct {
	FromID    int
	Text      string
	CreatedAt time.Time
}

func sendMessage(db *sql.DB, fromID, toID int, text string) error {
	_, err := db.Exec(`
        INSERT INTO messages (from_id, to_id, text)
        VALUES ($1, $2, $3)`,
		fromID, toID, text)
	return err
}

func getMessages(db *sql.DB, userID, otherID int) ([]Message, error) {
	rows, err := db.Query(`
        SELECT from_id, text, created_at FROM messages
        WHERE (from_id = $1 AND to_id = $2)
           OR (from_id = $2 AND to_id = $1)
        ORDER BY created_at`,
		userID, otherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		rows.Scan(&m.FromID, &m.Text, &m.CreatedAt)
		msgs = append(msgs, m)
	}

	return msgs, nil
}
