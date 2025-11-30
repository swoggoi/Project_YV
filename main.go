package main

import (
    "bufio"
    "database/sql"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "runtime"
    "strings"
    "time"

    _ "github.com/lib/pq"
)

type User struct {
    ID       int
    Username string
    Name     string
    Password string
}

type Message struct {
    FromID    int
    Text      string
    CreatedAt time.Time
}

func IdGenerator() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(90000000) + 10000000
}

func initDB() *sql.DB {
    connStr := "host=100.113.95.4 port=5432 user=postgres password=123123 dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    return db
}
func saveUser(db *sql.DB, user *User) error {
    _, err := db.Exec(`
        INSERT INTO users (id, username, name, password)
        VALUES ($1, $2, $3, $4)`,
        user.ID, user.Username, user.Name, user.Password)
    return err
}

func updatePassword(db *sql.DB, username, newPassword string) error {
    _, err := db.Exec(`
        UPDATE users SET password = $1 WHERE username = $2`,
        newPassword, username)
    return err
}

func updateUsername(db *sql.DB, oldUsername, newUsername string) error {
    _, err := db.Exec(`
        UPDATE users SET username = $1 WHERE username = $2`,
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

    fmt.Println("Пароль успешно изменён!")
    return true
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

func (u *User) ChangeUsername(db *sql.DB) bool {
    var NewUsername string
    fmt.Println("Введите новый username:")
    fmt.Scan(&NewUsername)

    if u.Username == NewUsername {
        fmt.Println("Username совпадает со старым")
        return false
    }

    old := u.Username
    u.Username = NewUsername

    if err := updateUsername(db, old, NewUsername); err != nil {
        fmt.Println("Ошибка записи в БД:", err)
    }

    fmt.Println("Username успешно изменён!")
    return true
}

func getUser(db *sql.DB, username string) (*User, error) {
    var u User
    err := db.QueryRow(`
        SELECT id, username, name, password 
        FROM users WHERE username = $1`,
        username).Scan(&u.ID, &u.Username, &u.Name, &u.Password)

    if err != nil {
        return nil, err
    }

    return &u, nil
}

func sendMessage(db *sql.DB, fromID, toID int, text string) error {
    _, err := db.Exec(`
        INSERT INTO messages (from_id, to_id, text)
        VALUES ($1, $2, $3)`,
        fromID, toID, text)
    return err
}

func getMessages(db *sql.DB, userID, otherID int) ([]Message, error) {
    rows, err := db.Query(`
        SELECT from_id, text, created_at FROM messages
        WHERE (from_id = $1 AND to_id = $2)
           OR (from_id = $2 AND to_id = $1)
        ORDER BY created_at`,
        userID, otherID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var msgs []Message
    for rows.Next() {
        var m Message
        rows.Scan(&m.FromID, &m.Text, &m.CreatedAt)
        msgs = append(msgs, m)
    }

    return msgs, nil
}

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

    reader := bufio.NewReader(os.Stdin)

    for {
        clearConsole()

        fmt.Println("┌──────────────────────────────────────────┐")
        fmt.Printf("│ %-40s │\n", HelloUser()+NewUser.Name)
        fmt.Printf("│ Username: @%-30s│\n", NewUser.Username)
        fmt.Printf("│ ID: %-36d │\n", NewUser.ID)
        fmt.Println("├──────────────────────────────────────────┤")
        fmt.Println("│ 1 — Показать пароль                      │")
        fmt.Println("│ 2 — Сменить пароль                       │")
        fmt.Println("│ 3 — Сменить username                     │")
        fmt.Println("│ 4 — Подключиться к чату по ID            │")
        fmt.Println("│ 0 — Выход                                │")
        fmt.Println("└──────────────────────────────────────────┘")

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
                fmt.Println("──────────────────────────────")

                for _, m := range msgs {
                    timeStr := m.CreatedAt.Format("15:04:05")

                    if m.FromID == NewUser.ID {
                        fmt.Printf("[Вы | %s]: %s\n", timeStr, m.Text)
                    } else {
                        fmt.Printf("[Собеседник | %s]: %s\n", timeStr, m.Text)
                    }
                }

                fmt.Println("──────────────────────────────")
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
