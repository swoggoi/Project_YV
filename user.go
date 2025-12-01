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

func handleLoginOrRegister(db *sql.DB, NewUser *User) {
	fmt.Println("Введите username:")
	fmt.Scan(&NewUser.Username)

	userFromDB, err := getUser(db, NewUser.Username)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Пользователь не найден. Регистрация.")
		fmt.Println("Введите имя:")
		fmt.Scan(&NewUser.Name)
		fmt.Println("Введите пароль:")
		fmt.Scan(&NewUser.Password)

		NewUser.ID = IdGenerator()
		if err := saveUser(db, NewUser); err != nil {
			fmt.Println("Ошибка сохранения:", err)
			return
		}
		fmt.Println("Регистрация успешна! Ваш ID:", NewUser.ID)

	case nil:
		fmt.Println("Введите пароль:")
		var inputPassword string
		fmt.Scan(&inputPassword)

		if !verifyLogin(userFromDB, inputPassword) {
			fmt.Println("Неверный пароль!")
			return
		}

		*NewUser = *userFromDB
		fmt.Println("Вход выполнен!")

	default:
		fmt.Println("Ошибка БД:", err)
		return
	}
}

func saveUser(db *sql.DB, user *User) error {
	if user.Password != "" && len(user.Password) < 4 || (len(user.Password) >= 2 && user.Password[:2] != "$2") {
		h, err := HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("hash error: %w", err)
		}
		user.Password = h
	}

	_, err := db.Exec(`
        INSERT INTO users (id, username, name, password)
        VALUES ($1, $2, $3, $4)`,
		user.ID, user.Username, user.Name, user.Password)
	return err
}

func verifyLogin(userFromDB *User, inputPassword string) bool {
	return CheckPassword(userFromDB.Password, inputPassword)
}

func (u *User) ChangePassword(db *sql.DB) bool {
	var newPassword string
	fmt.Println("Введите новый пароль:")
	fmt.Scan(&newPassword)

	if newPassword == "" {
		fmt.Println("Пустой пароль")
		return false
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		fmt.Println("Ошибка хэширования:", err)
		return false
	}
	u.Password = hash

	if err := updatePassword(db, u.Username, hash); err != nil {
		fmt.Println("Ошибка записи в БД:", err)
		return false
	}
	fmt.Println("Пароль успешно изменён!")
	return true
}
