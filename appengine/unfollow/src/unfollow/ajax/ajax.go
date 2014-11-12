package ajax

import (
    "encoding/json"
    "strings"
    "unfollow/urls"
    "unfollow/web"
)

type Function func(view *web.View) (interface{}, error)

func Get(name, path string, function Function) Function {
    subrouter.Handle(path, web.MakeView(wrap(function))).Methods("GET").Name("ajax:" + name)
    return function
}

func Post(name, path string, function Function) Function {
    subrouter.Handle(path, web.MakeView(wrap(function))).Methods("POST").Name("ajax:" + name)
    return function
}

func wrap(function Function) web.ViewFunction {
    return func(view *web.View) error {
        data, err := function(view)
        if err == nil {
            return Encode(view, data)
        }

        if err == web.ErrLoginRequired {
            return Encode(view, &struct {
                Redirect string `json:"redirect"`
            }{urls.Reverse("user:login").String()})
        }

        return err
    }
}

func Encode(view *web.View, data interface{}) error {
    view.Response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(view.Response).Encode(data)
    return nil
}

func Decode(view *web.View, data interface{}) error {
    s := view.Request.FormValue("data")
    if s == "" {
        return nil
    }
    return json.NewDecoder(strings.NewReader(s)).Decode(data)
}
