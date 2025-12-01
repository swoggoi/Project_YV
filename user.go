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

func handleLoginOrRegister(db *sql.DB, u *User) {
	var username, password string
	fmt.Println("Введите username:")
	fmt.Scan(&username)
	fmt.Println("Введите пароль:")
	fmt.Scan(&password)

	existingUser, err := getUser(db, username)
	if err == nil {
		// пользователь найден — проверяем пароль
		if CheckPassword(existingUser.Password, password) {
			*u = *existingUser
			fmt.Println("Успешный вход!")
		} else {
			fmt.Println("Неверный пароль")
			os.Exit(0)
		}
	} else {
		// пользователь не найден — регистрируем
		u.ID = IdGenerator()
		u.Username = username
		u.Name = username
		u.Password = password

		if err := saveUser(db, u); err != nil {
			fmt.Println("Ошибка регистрации:", err)
			os.Exit(0)
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

func (u *User) ChangeName(db *sql.DB) bool {
	var newName string
	fmt.Println("Введите новое имя:")
	fmt.Scan(&newName)

	if newName == "" {
		fmt.Println("Пустое имя!")
		return false
	}

	// проверим, что такого username ещё нет
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
	fmt.Print("Введите username: ")
	username := readLine()
	fmt.Print("Введите пароль: ")
	password := readLine()

	user, err := findUserByUsername(db, username)
	if err != nil {
		fmt.Println("Ошибка при поиске пользователя:", err)
		return nil
	}

	if user == nil {
		fmt.Println("Пользователь не найден.")
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Неверный пароль.")
		return nil
	}

	fmt.Println("Успешный вход!")
	return user
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

func registerUser(db *sql.DB, username, password string) (*User, error) {
	// Проверка: есть ли уже такой пользователь
	existing, err := findUserByUsername(db, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("пользователь уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var u User
	err = db.QueryRow(`
        INSERT INTO users (username, password)
        VALUES ($1, $2)
        RETURNING id, username, password, name, created_at
    `, username, string(hashedPassword)).Scan(&u.ID, &u.Username, &u.Password, &u.Name, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
