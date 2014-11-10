package ajax

import (
    "unfollow/urls"
)

var subrouter = urls.Routers[0].PathPrefix("/ajax/").Methods("GET", "POST").Subrouter()
