package api

import (
    "encoding/json"
    "unfollow/web"
)

type Function func(handler *web.Handler) (interface{}, error)

func Get(name, path string, function Function) Function {
    subrouter.Handle(path, web.MakeHandler(wrap(function))).Methods("GET").Name("api:" + name)
    return function
}

func Post(name, path string, function Function) Function {
    subrouter.Handle(path, web.MakeHandler(wrap(function))).Methods("POST").Name("api:" + name)
    return function
}

func wrap(function Function) web.HandlerFunction {
    return func(handler *web.Handler) error {
        data, err := function(handler)
        if err != nil {
            return err
        }

        if data != nil {
            handler.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
            json.NewEncoder(handler.Response).Encode(data)
        }
        return nil
    }
}
