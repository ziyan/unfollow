package task

import (
    "encoding/json"
    "unfollow/web"
)

type Function func(handler *web.Handler) (interface{}, error)

func Handle(name, path string, function Function) Function {
    subrouter.Handle(path, web.MakeHandler(wrap(function))).Methods("GET", "POST").Name("task:" + name)
    return function
}

func wrap(function Function) web.HandlerFunction {
    return func(handler *web.Handler) error {
        data, err := function(handler)
        if err != nil {
            return err
        }

        if data != nil {
            handler.Response.Header().Set("Content-Type", "application/json")
            json.NewEncoder(handler.Response).Encode(data)
        }
        return nil
    }
}
