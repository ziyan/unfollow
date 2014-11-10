package user

import (
    "unfollow/models"
    "unfollow/urls"
    "unfollow/utils/security"
    "unfollow/web"
    "strings"
)

const (
    sessionStateKey    = "state"
    sessionCallbackKey = "callback"
)

func createLogonUrl(view *web.View, next string) (string, error) {
    callback := ""
    if next != "" && strings.HasPrefix(next, "/") {
        callback = urls.ReverseByRequest(view.Request, "user:logon", "?", "next", next).String()
    } else {
        callback = urls.ReverseByRequest(view.Request, "user:logon").String()
    }

    state := security.GenerateRandomHexString(16)

    // save stuff to session
    view.Session.Values[sessionStateKey] = state
    view.Session.Values[sessionCallbackKey] = callback
    view.Session.Save()

    return "", nil
}

func handleLogon(view *web.View) error {
    if err := view.Request.ParseForm(); err != nil {
        return err
    }

    if view.Request.FormValue("state") != view.Session.Values[sessionStateKey] {
        return web.ErrBadRequest
    }

    code := view.Request.FormValue("code")
    if code == "" {
        // access denied
        return web.ErrBadRequest
    }

    // create or update the user
    u, err := models.PutUser(view.Database, models.UserKey(view.Database, 1), &models.User{
        Name:        "",
        Bio:         "",
        Username:    "",
        Email:       "",
        AccessToken: "",
    })
    if err != nil {
        return err
    }

    // login the user
    view.Session.Abandon(view.Context, view.Cache, view.Database, view.Response)
    view.Session.User = u
    view.Session.Save()

    return nil
}
