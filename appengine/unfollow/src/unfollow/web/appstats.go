package web

import (
    "unfollow/settings"
    "github.com/mjibson/appstats"
    "net/http"
)

func init() {
    if !settings.DEBUG {
        appstats.ShouldRecord = func(request *http.Request) bool {
            return false
        }
    }
}
