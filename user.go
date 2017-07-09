package user

import "time"

type User struct {
    id          int
    name        string
    surname     string
    email       string
    birthdate   time.Time
}

type Users []User

