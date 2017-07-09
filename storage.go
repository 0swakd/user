package main

import (
    "fmt"
    "time"
)

//var currentId int

var users Users
// var users = map[string]User

// Give us a seed data
func init() {
    StorageCreateUser(User{
        id: 1,
        name: "George",
        surname: "Doe",
        email: "george@doe.done",
        birthdate: time.Now(),
    })
}

// Will be usefull when we'll have distributed databases
func StorageFindUser(key string) (User, bool) {
    //return users[id];
    user, exist := users[key];

    if exist {
        return user, true
    }

    // return empty User if not found
    return User{}, false
}

//this is bad, I don't think it passes race condtions
func StorageCreateUser(u User) (User, bool) {

    key := UserKey(u);

    /* One day we'll have to ensure atomicity here */
    user, e := StorageFindUser(key);

    if e {
        return user, false
    }

    users[key] = u
    /* atomicity until here */

    return u, true
}

func StorageDestroyUser(key string) error {
    /* One day we'll have to ensure atomicity here */
    _, e := StorageFindUser(key);

    if !e {
        return fmt.Errorf("Could not find User with id of %s to delete", key)
    }

    delete(users, key)

    return nil
}
