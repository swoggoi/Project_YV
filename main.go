package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"database/sql"
    "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
	Name     string
	Password string
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "/home/pop/бд/nashbd.letter")
	if err != nil {
		panic(err)
	}
	return db
}

func saveUser(db *sql.DB, user User) error {
	_, err := db.Exec(`
        INSERT INTO users (username, name, password)
        VALUES (?, ?, ?)`,
		user.Username, user.Name, user.Password)
	return err
}

func updatePassword(db *sql.DB, username, newPassword string) error {
	_, err := db.Exec(`
        UPDATE users SET password = ? WHERE username = ?`,
		newPassword, username)
	return err
}

func updateUsername(db *sql.DB, oldUsername, newUsername string) error {
	_, err := db.Exec(`
        UPDATE users SET username = ? WHERE username = ?`,
		newUsername, oldUsername)
	return err
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

func (u *User) ChangePassword(db *sql.DB) bool {
	var NewPassword string
	fmt.Println("Введите новый пароль:")
	fmt.Scan(&NewPassword)
	if u.Password == NewPassword {
		fmt.Println("Пароли одинаковые")
		return false
	}
	u.Password = NewPassword
	if err := updatePassword(db, u.Username, NewPassword); err != nil {
		fmt.Println("Ошибка записи в БД:", err)
	}
	fmt.Println("Вы успешно сменили пароль!")
	return true
}

func (u *User) SetNewName() bool {
	var NewName string
	fmt.Println("Введите новое имя:")
	fmt.Scan(&NewName)
	if u.Name == NewName {
		fmt.Println("Имена совпадают")
		return false
	}
	u.Name = NewName
	return true
}

func (u *User) ChangeUsername(db *sql.DB) bool {
	var NewUsername string
	fmt.Println("Введите новый username:")
	fmt.Scan(&NewUsername)
	if u.Username == NewUsername {
		fmt.Println("user'ы одинаковые")
		return false
	}
	old := u.Username
	u.Username = NewUsername
	if err := updateUsername(db, old, NewUsername); err != nil {
		fmt.Println("Ошибка записи в БД:", err)
	}
	fmt.Println("Вы успешно сменили username!")
	return true
}

func main() {
	db := initDB()
	defer db.Close()

	NewUser := User{
		Username: "@sanya",
		Name:     "Sasha",
		Password: "Arina1978!",
	}
	for {
		clearConsole()
		fmt.Println("---------------------------")
		fmt.Printf("| %-25s|\n", "ПРИВЕТСТВУЕМ ВАС "+NewUser.Name)
		fmt.Printf("| %-25s|\n", "ваш username: "+NewUser.Username)
		fmt.Println("____________________________")
		fmt.Printf("|%-25s |\n", "Здравствуйте! "+NewUser.Name)
		fmt.Printf("|%-25s |\n", "Ваш username: @"+NewUser.Username)
		fmt.Println("| 1 - показать пароль      |")
		fmt.Println("| 2 - сменить пароль       |")
		fmt.Println("| 3 - сменить username     |")
		fmt.Println("| 4 - Сменить имя          |")
		fmt.Println("| 5 - Зайди в чат по id    |")
		fmt.Println("| 0 - выход                |")
		fmt.Println("|__________________________|")

		var choice int
		fmt.Print("Выберите пункт: ")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			fmt.Println("Ваш пароль:", NewUser.Password)
			fmt.Println("Нажмите Enter для продолжения...")
			fmt.Scanln()
		case 2:
			NewUser.ChangePassword(db)
			fmt.Println("Нажмите Enter для продолжения...")
			fmt.Scanln()
		case 3:
			NewUser.ChangeUsername(db)
			fmt.Println("Нажмите Enter для продолжения...")
			fmt.Scanln()
		case 4:
			NewUser.SetNewName()
			fmt.Println("Нажмите Enter для продолжения...")
			fmt.Scanln()
		//case 5:
		//
		//
		//
		case 0:
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неверный выбор!")
			fmt.Scanln()


		
		}
	}
}
