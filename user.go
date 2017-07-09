package main

import "time"

type User struct {
    id          int
    name        string
    surname     string
    email       string
    birthdate   time.Time
}

//type Users []User
type Users map[string]User

func UserKey(u User) string {
    return "keyisakey"
}
