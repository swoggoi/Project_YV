package main

import (
	"database/sql"
	"fmt"

	"github.com/fatih/color"
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
	db := initDB()
	defer db.Close()

	var NewUser User
	color.Red("═══════════════════════════════════════════════")
	fmt.Println("|                                             |")
	fmt.Println("|            ", HelloUser(), "                  |")
	fmt.Println("|                                             |")
	fmt.Println("|1 - войти                                    |")
	fmt.Println("|2 - зарегестрироваться                       |")
	fmt.Println("|                                             |")
	fmt.Println("|                                             |")
	color.Red("═══════════════════════════════════════════════")

	var choose int
	fmt.Println("Введите пункт:")
	fmt.Scan(&choose)

	switch choose {
	case 1:
		handleLoginOrRegister(db, &NewUser)
	case 2:
		handleLoginOrRegister(db, &NewUser)
	case 0:
		fmt.Println("Выход...")
		return
	default:
		fmt.Println("Неверный выбор")
	}

	// дальше — твое меню (без изменений)
	for {
		clearConsole()

		fmt.Printf("│ %-26s   │\n", HelloUser()+" "+NewUser.Name)
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
		case 2:
			NewUser.ChangePassword(db)
		case 3:
			NewUser.ChangeUsername(db)
		case 4:
			// чат
		case 0:
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неверный выбор")
		}
	}

	///
}
