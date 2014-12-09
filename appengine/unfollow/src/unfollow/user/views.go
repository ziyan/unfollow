package user

import (
    "net/http"
    "strconv"
    "strings"
    "unfollow/models"
    "unfollow/urls"
    "unfollow/web"
)

func LoginView(view *web.View) error {
    next := view.Request.URL.Query().Get("next")
    if next != "" && !strings.HasPrefix(next, "/") {
        next = ""
    }
    return view.Render(&struct{ Next string }{next}, "html.html", "base.html", "user/login.html")
}

func LogonView(view *web.View) error {
    if err := handleLogon(view); err != nil {
        return err
    }

    http.Redirect(view.Response, view.Request, urls.Reverse("page:network").String(), http.StatusTemporaryRedirect)
    return nil
}

func LogoutView(view *web.View) error {
    view.Session.Abandon(view.Context, view.Cache, view.Database, view.Response)
    http.Redirect(view.Response, view.Request, urls.Reverse("page:welcome").String(), http.StatusTemporaryRedirect)
    return nil
}

func WelcomeView(view *web.View) error {
    return nil
}

func UserView(view *web.View) error {
    var user *models.User = nil
    var err error = nil

    if view.Variables["username"] != "" {
        user, err = models.GetUserByUsername(view.Database, view.Variables["username"])
        if err != nil {
            return err
        }
    }
    if view.Variables["id"] != "" {
        id, err := strconv.ParseInt(view.Variables["id"], 10, 64)
        if err != nil {
            return err
        }

        user, err = models.GetUserByID(view.Database, id)
        if err != nil {
            return err
        }
    }

    if user == nil {
        return web.ErrNotFound
    }

    // if the user has a username
    if user.Username != "" && view.Variables["username"] != user.Username {
        http.Redirect(view.Response, view.Request, user.URL().String(), http.StatusTemporaryRedirect)
        return nil
    }

    data := struct {
        User *models.User
    }{user}
    view.Title = user.Name
    return view.Render(data, "html.html", "base.html", "user/user.html")
}
