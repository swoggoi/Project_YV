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

func handleLoginOrRegister(db *sql.DB, u *User) {
	var username, password string
	fmt.Println("Введите username:")
	fmt.Scan(&username)
	fmt.Println("Введите пароль:")
	fmt.Scan(&password)

	existingUser, err := getUser(db, username)
	if err == nil {
		// пользователь найден → проверяем пароль
		if CheckPassword(existingUser.Password, password) {
			*u = *existingUser
			fmt.Println("Успешный вход!")
		} else {
			fmt.Println("Неверный пароль")
		}
	} else {
		// пользователь не найден → регистрируем
		u.ID = IdGenerator()
		u.Username = username
		u.Name = username
		u.Password = password
		if err := saveUser(db, u); err != nil {
			fmt.Println("Ошибка регистрации:", err)
		} else {
			fmt.Println("Регистрация успешна!")
		}
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
	// userFromDB.Password — это ХЭШ
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
