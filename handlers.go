package main
// TODO : CRUD handlers/server/RPA could be a standalone package
// we could have a structure that only link him what to do with every action
// Or not
// Could be really handy to make improvements on all the RPAs

// TODO : RPA => move handlers.go, routes.go, router.go, error.go, logger.go in RPA package (rpa-service)
// TODO : User => move user.go in User package (Yes I know, it's already in, do this with places) (rpa-entity + user)
// TODO : Storage => evaluate the possibility to agnosticize storage in a RPAStock package (rpa-storage)
// TODO : Storage => config files and such things could be set on load
// TODO : Main => could be splitted in previous mentionned packages needing config files
// TODO : RPA => Starting to look like CRUD war

// TODO : Main => Actions focused : config storage, define actions, config handlers, config routes, starting service, eating popcorn

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/mux"
)

// TODO : RPA Standalone package could have a type struct with function
// TODO : RPA Standalone package could allow user to specify functions in this type
// TODO : Leaving the question of anonymous functions with anonymous types
// TODO : using interface{} as an anonymous (void *) type, don't know if its leggit
// TODO : Use a not defined yet type instead ? "Entry" for example, with a cast in the caller package
// TODO : maaah we are not going anywhere, time to bed
// TODO : hmmmmmm interfaces looks great might check this
//type Action struct {
//    Read        func (string) (interface{}, error)
//    Create      func (interface{}) (interface{}, error)
//    Update      func (interface{}, string) (interface{}, error)
//    Search      func (interface{}) (interface{}, error)
//    Delete      func (string) (interface{}, error)
//}
//var actions Actions
// TODO : RPA standalone package could have an initialisation function (would init be possible ?, config ? called during initialisation ?

// TODO : RPA standalone package could have this kind of modification
// TODO : using interface{} as an anonymous (void *) type, don't know if its leggit
//func Read(w http.ResponseWriter, r *http.Request, fr (string) (interface{}, error)) {
func UserRead(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var userId string
    userId = vars["userId"]

// TODO : RPA standalone package could have this kind of modification
//    entry, e := fr(userId)
//    entry, e := actions.Read(userId)
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

