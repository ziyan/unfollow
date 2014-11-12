package web

import (
    "github.com/mjibson/appstats"
    "net/http"
    "unfollow/settings"
)

func init() {
    if !settings.DEBUG {
        appstats.ShouldRecord = func(request *http.Request) bool {
            return false
        }
    }
}
