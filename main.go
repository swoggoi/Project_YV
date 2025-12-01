package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func main() {
	db := initDB()
	defer db.Close()

	var NewUser User
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
		if err := saveUser(db, &NewUser); err != nil {
			fmt.Println("Ошибка сохранения:", err)
			return
		}
		fmt.Println("Регистрация успешна! Ваш ID:", NewUser.ID)

	case nil:
		fmt.Println("Введите пароль:")
		var inputPassword string
		fmt.Scan(&inputPassword)

		if inputPassword != userFromDB.Password {
			fmt.Println("Неверный пароль!")
			return
		}

		NewUser = *userFromDB
		fmt.Println("Вход выполнен!")

	default:
		fmt.Println("Ошибка БД:", err)
		return
	}

	// меню
	reader := bufio.NewReader(os.Stdin)
	for {
		clearConsole()
		fmt.Println("┌──────────────────────────────┐")
		fmt.Printf("│ %-26s   │\n", HelloUser()+NewUser.Name)
		fmt.Printf("│ Username: @%-16s  │\n", NewUser.Username)
		fmt.Printf("│ ID: %-22d   │\n", NewUser.ID)
		fmt.Println("├──────────────────────────────┤")
		fmt.Println("│ 1 — Показать пароль          │")
		fmt.Println("│ 2 — Сменить пароль           │")
		fmt.Println("│ 3 — Сменить username         │")
		fmt.Println("│ 4 — Чат                      │")
		fmt.Println("│ 0 — Выход                    │")
		fmt.Println("└──────────────────────────────┘")

		var choice int
		fmt.Print("Выберите пункт: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Ваш пароль:", NewUser.Password)
			fmt.Scanln()
			fmt.Scanln()
		case 2:
			NewUser.ChangePassword(db)
			fmt.Scanln()
			fmt.Scanln()
		case 3:
			NewUser.ChangeUsername(db)
			fmt.Scanln()
			fmt.Scanln()
		case 4:
			fmt.Println("Введите ID собеседника:")
			var otherID int
			fmt.Scan(&otherID)
			for {
				clearConsole()
				msgs, _ := getMessages(db, NewUser.ID, otherID)
				fmt.Println("Чат:")
				for _, m := range msgs {
					timeStr := m.CreatedAt.Format("15:04:05")
					if m.FromID == NewUser.ID {
						fmt.Printf("[Вы | %s]: %s\n", timeStr, m.Text)
					} else {
						fmt.Printf("[Собеседник | %s]: %s\n", timeStr, m.Text)
					}
				}
				fmt.Println("Введите сообщение (или /exit):")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if text == "/exit" {
					break
				}
				if text != "" {
					sendMessage(db, NewUser.ID, otherID, text)
				}
			}
		case 0:
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неверный выбор")
			fmt.Scanln()
			fmt.Scanln()
		}
	}
}
