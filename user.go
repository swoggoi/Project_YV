package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int
	Username string
	Name     string
	Password string
}

func saveUser(db *sql.DB, user *User) error {
	_, err := db.Exec(`
        INSERT INTO users (id, username, name, password)
        VALUES ($1, $2, $3, $4)`,
		user.ID, user.Username, user.Name, user.Password)
	return err
}

func updatePassword(db *sql.DB, username, newPassword string) error {
	_, err := db.Exec(`
        UPDATE users SET password = $1 WHERE username = $2`,
		newPassword, username)
	return err
}

func updateUsername(db *sql.DB, oldUsername, newUsername string) error {
	_, err := db.Exec(`
        UPDATE users SET username = $1 WHERE username = $2`,
		newUsername, oldUsername)
	return err
}

func (u *User) ChangePassword(db *sql.DB) bool {
	var NewPassword string
	fmt.Println("Введите новый пароль:")
	fmt.Scan(&NewPassword)

	if u.Password == NewPassword {
		fmt.Println("Пароли одинаковые")
		return false
	}

	u.Password = NewPassword

	if err := updatePassword(db, u.Username, NewPassword); err != nil {
		fmt.Println("Ошибка записи в БД:", err)
	}

	fmt.Println("Пароль успешно изменён!")
	return true
}

func (u *User) ChangeUsername(db *sql.DB) bool {
	var NewUsername string
	fmt.Println("Введите новый username:")
	fmt.Scan(&NewUsername)

	if u.Username == NewUsername {
		fmt.Println("Username совпадает со старым")
		return false
	}

	old := u.Username
	u.Username = NewUsername

	if err := updateUsername(db, old, NewUsername); err != nil {
		fmt.Println("Ошибка записи в БД:", err)
	}

	fmt.Println("Username успешно изменён!")
	return true
}

func getUser(db *sql.DB, username string) (*User, error) {
	var u User
	err := db.QueryRow(`
        SELECT id, username, name, password 
        FROM users WHERE username = $1`,
		username).Scan(&u.ID, &u.Username, &u.Name, &u.Password)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
