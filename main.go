package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type User struct {
	Username string
	Name     string
	Password string
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

func (u *User) ChangePassword() bool {
	var NewPassword string
	fmt.Println("Введите новый пароль:")
	fmt.Scan(&NewPassword)
	if u.Password == NewPassword {
		fmt.Println("Пароли одинаковые")
		return false
	}
	u.Password = NewPassword
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

func (u *User) ChangeUsername() bool {
	var NewUsername string
	fmt.Println("Введите новый username:")
	fmt.Scan(&NewUsername)
	if u.Username == NewUsername {
		fmt.Println("user'ы одинаковые")
		return false
	}
	u.Username = NewUsername
	fmt.Println("Вы успешно сменили username!")
	return true
}

func main() {
	NewUser := User{
		Username: "",
		Name:     "",
		Password: "",
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
			NewUser.ChangePassword()
			fmt.Println("Нажмите Enter для продолжения...")
			fmt.Scanln()
		case 3:
			NewUser.ChangeUsername()
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
