package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"time"

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
	color.Red("|                                             |")
	color.Black("|            ", HelloUser(), "                   |")
	color.Red("|                                             |")
	color.Red("|1 - войти                                    |")
	color.Red("|2 - зарегестрироваться                       |")
	color.Red("|                                             |")
	color.Red("|                                             |")
	color.Red("═══════════════════════════════════════════════")

	var choose int
	fmt.Println("Введите пункт:")
	fmt.Scan(&choose)

	switch choose {
	case 1:
		handleLoginOrRegister(db, &NewUser)
	case 0:
		fmt.Println("Выход...")
		return
	default:
		fmt.Println("Неверный выбор")
	}

	// дальше — твое меню (без изменений)
	reader := bufio.NewReader(os.Stdin)
	for {
		clearConsole()
		// ...
		// чат/смена пароля/смена username
		// ...
		_ = reader // чтобы пример был самодостаточный
		time.Sleep(10 * time.Millisecond)
		return
	}
}
