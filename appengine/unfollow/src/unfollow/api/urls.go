package api

import (
    "unfollow/urls"
)

var subrouter = urls.Routers[0].PathPrefix("/api").Methods("GET", "POST").Subrouter()
