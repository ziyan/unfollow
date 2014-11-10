package web

import (
    "appengine"
    "unfollow/utils/cache"
    "unfollow/utils/db"
    "unfollow/web/templates"
    "github.com/gorilla/mux"
    "github.com/mjibson/appstats"
    "net/http"
)

type Handler struct {
    Context   appengine.Context
    Request   *http.Request
    Response  http.ResponseWriter
    Variables map[string]string
    Data      interface{}
    Cache     *cache.Cache
    Database  *db.Database
}

type HandlerFunction func(*Handler) error

func MakeHandler(function HandlerFunction) appstats.Handler {
    return appstats.NewHandler(func(context appengine.Context, response http.ResponseWriter, request *http.Request) {
        handler := new(Handler)

        handler.Context = context
        handler.Request = request
        handler.Response = response
        handler.Variables = mux.Vars(request)

        handler.Cache = cache.New(context)
        defer handler.Cache.Clear()

        handler.Database = db.New(context, handler.Cache)

        // set up panic handling
        defer RecoverError(context, request, response)

        // do the work
        err := function(handler)

        // handle error
        HandleError(context, request, response, err)
    })
}

func (handler *Handler) Render(data interface{}, filenames ...string) error {
    handler.Data = data
    return templates.Render(handler.Response, handler, filenames...)
}

func (handler *Handler) RenderToString(data interface{}, filenames ...string) (string, error) {
    handler.Data = data
    return templates.RenderToString(handler, filenames...)
}
