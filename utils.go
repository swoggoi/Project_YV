package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func IdGenerator() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(90000000) + 10000000
}

func clearConsole() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func HelloUser() string {
	hour := time.Now().Hour()
	switch {
	case hour >= 9 && hour <= 11:
		return "Доброе утро!"
	case hour >= 12 && hour <= 18:
		return "Добрый день!"
	case hour >= 19 && hour <= 23:
		return "Добрый вечер!"
	default:
		return "Доброй ночи!"
	}
}
func register(db *sql.DB) *User {
	fmt.Print("Введите новый username: ")
	username := readLine()

	fmt.Print("Введите пароль: ")
	password := readLine()

	// Проверка: существует ли уже такой пользователь
	existing, err := findUserByUsername(db, username)
	if err != nil {
		fmt.Println("Ошибка при проверке:", err)
		return nil
	}
	if existing != nil {
		fmt.Println("Пользователь уже существует.")
		return nil
	}

	// Хэшируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Ошибка хэширования:", err)
		return nil
	}

	// Создаём пользователя
	var user User
	err = db.QueryRow(`
        INSERT INTO users (username, password, name)
        VALUES ($1, $2, $3)
        RETURNING id, username, password, name
    `, username, string(hashed), username).Scan(&user.ID, &user.Username, &user.Password, &user.Name)

	if err != nil {
		fmt.Println("Ошибка регистрации:", err)
		return nil
	}

	fmt.Println("Регистрация успешна!")
	return &user
}

func MainMenu() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                            ║")
	fmt.Println("║", HelloUser(), "                                             ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  1 — 🔐 Войти                                              ║")
	fmt.Println("║  2 — 📝 Зарегистрироваться                                 ║")
	fmt.Println("║  0 — 🚪 Выход                                              ║")
	fmt.Println("║                                                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println("Введите пункт:")
}

func UserMenu() {
	fmt.Println("├──────────────────────────────┤")
	fmt.Println("│ 1 — Сменить пароль           │")
	fmt.Println("│ 2 — Сменить username         │")
	fmt.Println("│ 3 — Сменить имя              │")
	fmt.Println("│ 4 — Войти в чат по id        │")
	fmt.Println("│ 0 — Выход                    │")
	fmt.Println("└──────────────────────────────┘")
}
func isValidIDString(id string) bool {
	if len(id) != 8 {
		return false
	}
	for _, ch := range id {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
func findUserByID(db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow(`SELECT id, username, password, name FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Username, &u.Password, &u.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func startChat(db *sql.DB, currentUser *User) {
	fmt.Print("Введите ID собеседника: ")
	var targetID int
	fmt.Scan(&targetID)

	if targetID == currentUser.ID {
		fmt.Println("Нельзя писать самому себе.")
		return
	}

	partner, err := findUserByID(db, targetID)
	if err != nil {
		fmt.Println("Ошибка при поиске пользователя:", err)
		return
	}
	if partner == nil {
		fmt.Println("Пользователь с таким ID не найден.")
		return
	}

	fmt.Printf("Чат с %s (@%s)\n", partner.Name, partner.Username)
	fmt.Println("Введите сообщение (или 'exit' для выхода):")

	for {
		fmt.Print("Вы: ")
		text := readLine()
		if text == "exit" {
			break
		}

		// Сохраняем сообщение
		_, err := db.Exec(`
            INSERT INTO messages (user_id, receiver_id, content)
            VALUES ($1, $2, $3)
        `, currentUser.ID, partner.ID, text)
		if err != nil {
			fmt.Println("Ошибка отправки:", err)
			continue
		}

		fmt.Println("Сообщение отправлено.")
	}
}
func showChatHistory(db *sql.DB, userID, partnerID int) {
	rows, err := db.Query(`
        SELECT user_id, content, created_at
        FROM messages
        WHERE (user_id = $1 AND receiver_id = $2)
           OR (user_id = $2 AND receiver_id = $1)
        ORDER BY created_at
    `, userID, partnerID)
	if err != nil {
		fmt.Println("Ошибка чтения истории:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var senderID int
		var content string
		var createdAt string
		rows.Scan(&senderID, &content, &createdAt)

		prefix := "Собеседник"
		if senderID == userID {
			prefix = "Вы"
		}
		fmt.Printf("[%s] %s: %s\n", createdAt, prefix, content)
	}
}
