package user

import (
    "unfollow/ajax"
    "unfollow/web"
)

var _ = ajax.Get("user:login", "/user/login", func(view *web.View) (interface{}, error) {
    data := &struct {
        Next string `json:"next"`
    }{}

    if err := ajax.Decode(view, data); err != nil {
        return nil, err
    }

    redirect, err := createLogonUrl(view, data.Next)
    if err != nil {
        return nil, err
    }

    return struct {
        Redirect string `json:"redirect"`
    }{redirect}, nil
})

