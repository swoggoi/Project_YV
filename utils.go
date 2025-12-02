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
	return rand.Intn(90000000) + 10000000
}

func GenerateUniqueID(db *sql.DB) int {
	for {
		id := IdGenerator()

		var exists int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&exists)

		if exists == 0 {
			return id
		}
	}
}

func clearConsole() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
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

	existing, err := findUserByUsername(db, username)
	if err != nil {
		fmt.Println("Ошибка при проверке:", err)
		return nil
	}
	if existing != nil {
		fmt.Println("Пользователь уже существует.")
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Ошибка хэширования:", err)
		return nil
	}

var user User
user.ID = GenerateUniqueID(db)

err = db.QueryRow(`
    INSERT INTO users (id, username, password, name)
    VALUES ($1, $2, $3, $4)
    RETURNING id, username, password, name
`, user.ID, username, string(hashed), username).Scan(
    &user.ID, &user.Username, &user.Password, &user.Name,
)

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
	fmt.Println("║", HelloUser(), "                                              ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  1 — 🔐 Войти                                              ║")
	fmt.Println("║  2 — 📝 Зарегистрироваться                                 ║")
	fmt.Println("║  0 — 🚪 Выход                                              ║")
	fmt.Println("║                                                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println("Введите пункт:")
}

func UserMenu() {
    fmt.Println("┌──────────────────────────────┐")
    fmt.Println("│ 1 — Сменить пароль           │")
    fmt.Println("│ 2 — Сменить username         │")
    fmt.Println("│ 3 — Сменить имя              │")
    fmt.Println("│ 4 — Войти в чат по id        │")
    fmt.Println("│ 0 — Выход                    │")
    fmt.Println("└──────────────────────────────┘")
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

func showChatHistory(db *sql.DB, userID, partnerID int) {
	rows, err := db.Query(`
        SELECT from_id, text, created_at
        FROM messages
        WHERE (from_id = $1 AND to_id = $2)
           OR (from_id = $2 AND to_id = $1)
        ORDER BY created_at
    `, userID, partnerID)
	if err != nil {
		fmt.Println("Ошибка чтения истории:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var fromID int
		var text string
		var createdAt time.Time
		rows.Scan(&fromID, &text, &createdAt)

		sender := "Собеседник"
		if fromID == userID {
			sender = "Вы"
		}

		fmt.Printf("[%s] %s: %s\n",
			createdAt.Format("15:04"),
			sender,
			text,
		)
	}
}

func startChat(db *sql.DB, currentUser *User, partnerID int) {
	for {
		clearConsole()

		partner, _ := findUserByID(db, partnerID)
		fmt.Printf("Чат с %s (@%s)\n\n", partner.Name, partner.Username)

		showChatHistory(db, currentUser.ID, partnerID)

		fmt.Println("\nВведите сообщение (или 'exit' для выхода):")
		fmt.Print("Вы: ")

		text := readLine()
		if text == "exit" {
			break
		}

		_, err := db.Exec(`
            INSERT INTO messages (from_id, to_id, text)
            VALUES ($1, $2, $3)
        `, currentUser.ID, partnerID, text)

		if err != nil {
			fmt.Println("Ошибка отправки:", err)
			time.Sleep(1 * time.Second)
			continue
		}
	}
}
