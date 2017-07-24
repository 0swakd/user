package main
// TODO : CRUD handlers could be a standalone package

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/mux"
)


func UserRead(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var userId string
    userId = vars["userId"]

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
        return
    }

    /* TODO print and return error */
    t, e := StorageCreateUser(user)
    if e != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict) // cannot be treated
        if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusConflict, Text: e.Error()}); err != nil {
            panic(err)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
    var userId string
    var user User

    vars := mux.Vars(r)
    userId = vars["userId"]

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
        return
    }

    /* TODO print and return error */
    t, e := StorageUpdateUser(user, userId)
    if e != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict) // cannot be treated
        if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusConflict, Text: e.Error()}); err != nil {
            panic(err)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}

func UserSearch(w http.ResponseWriter, r *http.Request) {
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
        return
    }

    /* TODO print and return error */
    t, e := StorageSearchUser(user)
    if e != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNotFound) // cannot be treated
        if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusConflict, Text: e.Error()}); err != nil {
            panic(err)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var userId string
    userId = vars["userId"]

    e := StorageDeleteUser(userId)
    if e == nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(""); err != nil {
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

