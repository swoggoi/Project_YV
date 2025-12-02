package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func getUser(db *sql.DB, username string) (*User, error) {
	var u User
	err := db.QueryRow(`
        SELECT id, username, name, password 
        FROM users WHERE username = $1`, username).
		Scan(&u.ID, &u.Username, &u.Name, &u.Password)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	db := initDB()
	defer db.Close()

	var NewUser User
	MainMenu()

	var choose int
	fmt.Scan(&choose)

	switch choose {
	case 1:
		user := login(db)
		if user != nil {
			NewUser = *user
		}
	case 2:
		user := register(db)
		if user != nil {
			NewUser = *user
		}
	case 0:
		fmt.Println("Выход...")
		return
	default:
		fmt.Println("Неверный выбор")
		return
	}

	if NewUser.ID == 0 {
		fmt.Println("Ошибка входа или регистрации. Попробуйте снова.")
		return
	}

	// Основное меню
	for {
		clearConsole()
		fmt.Printf(" %-26s   \n", HelloUser()+" "+NewUser.Name)
		fmt.Printf(" Username: @%-16s  \n", NewUser.Username)
		fmt.Printf(" ID: %-22d   \n", NewUser.ID)
		UserMenu()

		var choice int
		fmt.Print("Выберите пункт: ")
		fmt.Scan(&choice)

		switch choice {

		case 1:
			NewUser.ChangePassword(db)

		case 2:
			NewUser.ChangeUsername(db)

		case 3:
			NewUser.ChangeName(db)

		case 4:

			fmt.Print("Введите ID пользователя: ")
			var targetID int
			fmt.Scan(&targetID)

			if targetID == NewUser.ID {
				fmt.Println("Нельзя писать самому себе.")
				break
			}

			partner, err := findUserByID(db, targetID)
			if err != nil {
				fmt.Println("Ошибка поиска:", err)
				break
			}
			if partner == nil {
				fmt.Println("Пользователь не найден.")
				break
			}

			clearConsole()
			fmt.Printf("Чат с %s (@%s)\n", partner.Name, partner.Username)
			fmt.Println("История сообщений:\n")

			showChatHistory(db, NewUser.ID, partner.ID)

			fmt.Println("\nВведите сообщение (или 'exit' для выхода):")
			startChat(db, &NewUser, partner.ID)

		case 0:
			fmt.Println("Выход...")
			return

		default:
			fmt.Println("Неверный выбор")
		}
	}
}
