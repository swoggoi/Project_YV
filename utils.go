package main

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

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
		return "Доброе утро! "
	case hour >= 12 && hour <= 18:
		return "Добрый день! "
	case hour >= 19 && hour <= 23:
		return "Добрый вечер! "
	default:
		return "Доброй ночи! "
	}
}
