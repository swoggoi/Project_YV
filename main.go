package main

import (
	"database/sql"
	"fmt"
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
	// дальше — твое меню (без изменений)
	for {
		clearConsole()
		fmt.Printf("│ %-26s   │\n", HelloUser()+" "+NewUser.Name)
		fmt.Printf("│ Username: @%-16s  │\n", NewUser.Username)
		fmt.Printf("│ ID: %-22d   │\n", NewUser.ID)
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
			var IdEnter string
			fmt.Println("Введите id пользователя:")
			fmt.Scan(&IdEnter)

		case 0:
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неверный выбор")
		}
	}

}
