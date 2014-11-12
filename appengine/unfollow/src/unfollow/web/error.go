package web

import (
    "appengine"
    "errors"
    "net/http"
    "unfollow/settings"
    "unfollow/urls"
)

var (
    ErrNotFound      = errors.New("web: not found")
    ErrForbidden     = errors.New("web: forbidden")
    ErrBadRequest    = errors.New("web: bad request")
    ErrInternal      = errors.New("web: internal server error")
    ErrLoginRequired = errors.New("web: login required")
)

func HandleError(context appengine.Context, request *http.Request, response http.ResponseWriter, err error) {
    if err == nil {
        return
    }

    switch err {
    case ErrLoginRequired:
        http.Redirect(response, request, urls.Reverse("user:login", "?", "next", request.URL.Path).String(), http.StatusTemporaryRedirect)
    case ErrNotFound:
        http.NotFound(response, request)
    case ErrForbidden:
        http.Error(response, err.Error(), http.StatusForbidden)
    case ErrBadRequest:
        http.Error(response, err.Error(), http.StatusBadRequest)
    default:
        http.Error(response, err.Error(), http.StatusInternalServerError)
        context.Errorf("web: error: %q", err)
    }
}

func RecoverError(context appengine.Context, request *http.Request, response http.ResponseWriter) {
    if settings.DEBUG {
        return
    }

    err := recover()
    if err == nil {
        return
    }

    http.Error(response, "panic", http.StatusInternalServerError)
    context.Criticalf("web: panic: %q", err)
}
