package csrf

import (
    "appengine"
    "unfollow/settings"
    "unfollow/utils/security"
    "unfollow/web/sessions"
    "errors"
    "net/http"
)

var ErrTokenInvalid = errors.New("csrf: invalid token")

const (
    COOKIE_NAME = "csrftoken"
    COOKIE_PATH = "/"
    HEADER_NAME = "X-Csrf-Token"
)

func Do(context appengine.Context, request *http.Request, response http.ResponseWriter, session *sessions.Session) error {

    if session.CSRFToken == "" {
        session.CSRFToken = security.GenerateRandomID()
        session.Save()
    }

    if cookie, _ := request.Cookie(COOKIE_NAME); cookie == nil || cookie.Value != session.CSRFToken {
        http.SetCookie(response, &http.Cookie{
            Name:   COOKIE_NAME,
            Value:  session.CSRFToken,
            Path:   COOKIE_PATH,
            Secure: settings.SECURE,
        })
    }

    if request.Method == "POST" || request.Method == "PUT" {
        headers := request.Header[HEADER_NAME]
        if len(headers) == 0 || headers[0] != session.CSRFToken {
            http.Error(response, "CSRF Token Invalid", http.StatusForbidden)
            return ErrTokenInvalid
        }
    }

    return nil
}
