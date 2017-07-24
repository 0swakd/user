package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

/* TODO : Get rid of action in url */
/* TODO : Put a version in the url */
/* TODO : Possibilite de venir ajouter des routes a la liste depuis un package externe (monitoring, managing) */
var routes = Routes{
    Route{
        "UserCreate",
        "POST",
        "/user/create",
        UserCreate,
    },
    Route{
        "UserRead",
        "GET",
        "/user/id/{userId}",
        UserRead,
    },
    Route{
        "UserUpdate",
        "POST",
        "/user/id/{userId}",
        UserUpdate,
    },
    Route{
        "UserDelete",
        "DELETE",
        "/user/id/{userId}",
        UserDelete,
    },
    Route{
        "UserSearch",
        "POST",
        "/user/search",
        UserSearch,
    },
}

