package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	Name      string
	Password  string
	CreatedAt string
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (u *User) ChangeUsername(db *sql.DB) bool {
	var newUsername string
	fmt.Println("Введите новый username:")
	fmt.Scan(&newUsername)

	if newUsername == "" {
		fmt.Println("Пустой username")
		return false
	}

	// проверим, что такого username ещё нет
	var exists int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = $1`, newUsername).Scan(&exists)
	if err != nil {
		fmt.Println("Ошибка проверки:", err)
		return false
	}
	if exists > 0 {
		fmt.Println("Такой username уже занят")
		return false
	}

	// обновляем в БД
	_, err = db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, newUsername, u.ID)
	if err != nil {
		fmt.Println("Ошибка обновления:", err)
		return false
	}

	u.Username = newUsername
	fmt.Println("Username успешно изменён!")
	return true
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

func (u *User) ChangeName(db *sql.DB) bool {
	var newName string
	fmt.Println("Введите новое имя:")
	fmt.Scan(&newName)

	if newName == "" {
		fmt.Println("Пустое имя!")
		return false
	}

	//такого username ещё нет
	var exists int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE name = $1`, newName).Scan(&exists)
	if err != nil {
		fmt.Println("Ошибка проверки:", err)
		return false
	}
	if exists > 0 {
		fmt.Println("Такой username уже занят")
		return false
	}

	// обновляем в БД
	_, err = db.Exec(`UPDATE users SET name = $1 WHERE id = $2`, newName, u.ID)
	if err != nil {
		fmt.Println("Ошибка обновления:", err)
		return false
	}

	u.Name = newName
	fmt.Println("Username успешно изменён!")
	return true
}

func login(db *sql.DB) *User {
	for {
		fmt.Print("Введите username: ")
		username := readLine()
		if username == "" {
			fmt.Println("Вы ввели пустоту!!!")
			continue
		}

		fmt.Print("Введите пароль: ")
		password := readLine()
		if password == "" {
			fmt.Println("Вы ввели пустоту!!!")
			continue
		}

		user, err := findUserByUsername(db, username)
		if err != nil {
			fmt.Println("Ошибка при поиске пользователя:", err)
			continue
		}

		if user == nil {
			fmt.Println("Пользователь не найден.")
			continue
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			fmt.Println("Неверный пароль.")
			continue
		}

		fmt.Println("Успешный вход!")
		return user
	}
}

func findUserByUsername(db *sql.DB, username string) (*User, error) {
	var u User
	err := db.QueryRow(`SELECT id, username, password, name FROM users WHERE username = $1`, username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Name)
	if err == sql.ErrNoRows {
		return nil, nil // пользователь не найден
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
