package task

import (
    "unfollow/urls"
)

var subrouter = urls.Routers[0].PathPrefix("/task/").Methods("GET", "POST").Subrouter()
