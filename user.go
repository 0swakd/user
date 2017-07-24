package main
// TODO : Can become a standalone package

import (
    "time"
    "fmt"
    "crypto/md5"
    "encoding/hex"
    "strconv"
    "reflect"
)

/* Voir pour faire des structures de sortie un peu moins gores... */
type User struct {
    Id          int         `json:"-"`
    Name        string      `json:"name"`
    Surname     string      `json:"surname"`
    Email       string      `json:"email"`
    Password    string      `json:"password"` /* TODO meeeeeeh noooooo don't do that, do a func updatePassword or anything else but not that */
    Birthdate   time.Time   `json:"birthdate"`
    Activity    time.Time   `json:"activity"`
    Salt        string      `json:"-"`
}

func UserKey(u User) (string, error) {
    if u.Salt == "" {
        return "", fmt.Errorf("No user salt to make key")
    }
    h := md5.New()
    h.Write([]byte(u.Salt))
    return hex.EncodeToString(h.Sum(nil)), nil
}

func UserSalt(u User, localSalt string) string {
    h := md5.New()
    s := fmt.Sprintf("Salt;user;%s;%d", localSalt, u.Id)
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}

func UserToTable(u User) []string {
    v := reflect.ValueOf(u)

    record := make([]string, v.NumField())

    for i := 0; i < v.NumField(); i++ {
        record[i] = fmt.Sprintf("%v", v.Field(i).Interface())
    }

    return record
}

func UserFromTable(r []string) (User, error) {
    /* test length */
    i, _/*err*/ := strconv.Atoi(r[0])
    /*
    if err {
        return nil, err
    }*/
    t, _/*err*/ := time.Parse("20060102", r[5])
    /*if err {
        return nil, err
    }*/
    a, _/*err*/ := time.Parse("20060102T150405", r[6])
    /*if err {
        return nil, err
    }*/
    u := User{
        Id:             i,
        Name:           r[1],
        Surname:        r[2],
        Email:          r[3],
        Password:       r[4],
        Birthdate:      t,
        Activity:       a,
        Salt:           r[7],
    }

    return u, nil
}

