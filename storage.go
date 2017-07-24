package main

import (
    "fmt"
    "os"
    "bufio"
    "encoding/csv"
    "io"
    "time"
    "strings"
)

// TODO : We should save the key in user so that it could be not structure dependant
type StorageUser struct {
    Id          int         `json:"-"`
    Name        string      `json:"name"`
    Surname     string      `json:"surname"`
    Email       string      `json:"email"`
    Password    string      `json:"-"`
    Birthdate   time.Time   `json:"-"`
    Activity    time.Time   `json:"-"`
    Salt        string      `json:"key"` /* Here comes the ugly part, we are so doomed */
}

// TODO : split response between exact match and partial match
type StorageSearchMatch struct {
    Partial     []string
    Exact       []string
}

type StorageSearchResult struct {
    Name        StorageSearchMatch      `json:"name"`
    Surname     StorageSearchMatch      `json:"surname"`
    Email       StorageSearchMatch      `json:"email"`
}

//type StorageSearchResult struct {
//    Name        []string      `json:"name"`
//    Surname     []string      `json:"surname"`
//    Email       []string      `json:"email"`
//}

// TODO retirer fmt et les faire gerer par le heandler ou autre

//type Users map[string]User
//var users = make(Users)

// A decomposer si besoin
var users = make(map[string]User)

func init() {
    fmt.Fprintf(os.Stderr, "Initializing storage\n")
}

// Will be usefull when we'll have distributed databases
func StorageFindUser(key string) (StorageUser, bool) {
    user, exist := users[key]

    if exist {
        fmt.Println(user)
        user.Salt = key
        return StorageUser(user), true
    }

    // return empty User if not found
    // TODO anonymiser l'adresse mail quand le retour n'est pas pour le user actuel, ce sera fait par l'api centrale
    return StorageUser(User{}), false
}



func StorageCreateUser(u User) (StorageUser, error) {

    StorageUserFill(&u)
    // TODO verbose
    fmt.Println(u)

    key, err := UserKey(u);
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error : %s\n", err.Error())
        u.Salt = ""
        return StorageUser(u), err
    }

    /* One day we'll have to ensure atomicity here */
    for ukey, user := range users {
        if key == ukey {
            fmt.Fprintf(os.Stderr, "User already exist for key %s\n", key)
            fmt.Println(u)
            user.Salt = key
            return StorageUser(user), fmt.Errorf("User already exist for key %s", key)
        }

        if err, valid := StorageUnicityValid(user, u); valid == false {
            fmt.Fprintf(os.Stderr, "Error : %s\n", err.Error())
            fmt.Println(u)
            user.Salt = key
            //return user, fmt.Errorf("User already exist for key %s", err)
            return StorageUser(user), err
        }
    }

    // TODO verbose 
    fmt.Println(u)

    users[key] = u

    u.Salt = key

    return StorageUser(u), nil
}

func StorageDeleteUser(key string) error {
    /* One day we'll have to ensure atomicity here */
    _, found := StorageFindUser(key);

    if !found {
        return fmt.Errorf("Could not find User with id of %s to delete", key)
    }

    delete(users, key)

    return nil
}

func StorageLoadUsers(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return fmt.Errorf("Error reading %s", file)
    }
    defer f.Close()

    r := csv.NewReader(bufio.NewReader(f))
    for {
        record, err := r.Read()
        if record == nil && err == io.EOF {
            return nil
        }

        if err != nil {
            return fmt.Errorf("Error reading line")
        }

        u, _/*err*/ := UserFromTable(record)
        /*if err {
            return err
        }*/

        _, err = StorageCreateUser(u)

        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to load user\n")
            fmt.Println(u)
            fmt.Fprintf(os.Stderr, "Resuming without loading it\n")
        }
    }


    return nil
}

func StorageSaveUsers(file string) error {
    f, err := os.OpenFile(file, os.O_WRONLY, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error writing\n")
        return fmt.Errorf("Error writing %s", file)
    }
    defer f.Close()

    w := csv.NewWriter(f)
    defer w.Flush()

    for _, user := range users {
        // TODO verbose
        fmt.Println(user)
        r := UserToTable(user)

        err := w.Write(r)
        if err != nil {
            /* Put that in a raw file for debug */
            fmt.Fprintf(os.Stderr, "Failed to save user\n")
            fmt.Println(user)
        }
    }

    return nil
}

func StorageFirstFreeId() int {
    var id = 0

    for u := range users {
        if id <= users[u].Id {
            id = users[u].Id + 1
        }
    }

    return id
}


func StorageUserFill(u *User) error {
    // TODO verbose
    fmt.Println(*u)

    if (*u).Id == 0 {
        (*u).Id = StorageFirstFreeId()
    }

    if (*u).Salt == "" {
        (*u).Salt = UserSalt(*u, "")
    }

    // TODO verbose
    fmt.Println(*u)
    return nil
}

/* Might exist a more flexible way to do this... */
func StorageUnicityValid(u1 User, u2 User) (error, bool) {

    if u1.Id == u2.Id {
        return fmt.Errorf("Id already exist"), false
    }

    if u1.Email == u2.Email {
        return fmt.Errorf("Email already exist"), false
    }

    return nil, true
}

func StorageUpdateUser(u User, key string) (StorageUser, error) {
    user, exist := users[key]

    if !exist {
        //return User{}, fmt.Errorf("User (%s) not found", key)
        return StorageUser{}, fmt.Errorf("User (%s) not found", key)
    }

    if u.Name != "" {
        user.Name = u.Name
    }

    if u.Surname != "" {
        user.Surname = u.Surname
    }

    if u.Password != "" {
        user.Password = u.Password
    }

    if u.Email != "" {
        user.Email = u.Email
    }

    users[key] = user

    user.Salt = key

    return StorageUser(user), nil
}

/* TODO : Can be splitted in index RPAs and managed by Entry RPAs */
// TODO : split response between exact match and partial match
func StorageSearchUser(u User) (StorageSearchResult, error) {
    var result StorageSearchResult

    for key, user := range users {
        //if u.Name != "" && user.Name == u.Name {
        if len(u.Name) > 4 && strings.Contains(strings.ToUpper(user.Name), strings.ToUpper(u.Name)) {
            if user.Name == u.Name {
                result.Name.Exact = append(result.Name.Exact, key)
            } else {
                result.Name.Partial = append(result.Name.Partial, key)
            }
            //result.Name = append(result.Name, key)
        }
        //if u.Surname != "" && user.Surname == u.Surname {
        if len(u.Surname) > 4 && strings.Contains(strings.ToUpper(user.Surname), strings.ToUpper(u.Surname)) {
            if user.Surname == u.Surname {
                result.Surname.Exact = append(result.Surname.Exact, key)
            } else {
                result.Surname.Partial = append(result.Surname.Partial, key)
            }
            //result.Surname = append(result.Surname, key)
        }
        // TODO : Email lookup should/could be an exact match search
        // TODO : Entry RPAs should anonymise them when searched
        //if u.Email != "" && user.Email == u.Email {
        if len(u.Email) > 4 && strings.Contains(strings.ToUpper(user.Email), strings.ToUpper(u.Email)) {
            if user.Email == u.Email {
                result.Email.Exact = append(result.Email.Exact, key)
            } else {
                result.Email.Partial = append(result.Email.Partial, key)
            }
            //result.Email = append(result.Email, key)
        }
    }

    return result, nil
}
