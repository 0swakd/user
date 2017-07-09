package main

import (
    "encoding/json"
//    "fmt"
    "io"
    "io/ioutil"
    "net/http"
//    "strconv"

    "github.com/gorilla/mux"
)


func UserShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var userId string
    //var userId int
    //var err error
    userId = vars["userId"]
    //if userId, err = strconv.Atoi(vars["userId"]); err != nil {
    //    panic(err)
    //}
    //user := StorageFindUser(userId)
    user, e := StorageFindUser(userId)
    if e {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(user); err != nil {
            panic(err)
        }
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
        panic(err)
    }

}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"New User"}' http://localhost:8080/todos
*/
func UserCreate(w http.ResponseWriter, r *http.Request) {
    var user User
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    t, _ := StorageCreateUser(user)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}

