package user

import (
    "github.com/ziyan/oauth"
    "unfollow/ajax"
    "unfollow/utils/twitter"
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

var _ = ajax.Get("user:search", "/user/search", func(view *web.View) (interface{}, error) {
    data := &struct {
        Query string `json:"query"`
    }{}

    if err := ajax.Decode(view, data); err != nil {
        return nil, err
    }

    accessToken, err := oauth.DecodeToken(view.Session.User.AccessToken)
    if err != nil {
        return nil, err
    }

    t := twitter.New(view.Context, accessToken)
    tweets, err := t.Search(data.Query)
    if err != nil {
        return nil, err
    }

    return struct {
        Tweets []*twitter.Tweet `json:"tweets"`
    }{tweets}, nil
})

var _ = ajax.Get("user:mentions", "/user/mentions", func(view *web.View) (interface{}, error) {
    accessToken, err := oauth.DecodeToken(view.Session.User.AccessToken)
    if err != nil {
        return nil, err
    }

    t := twitter.New(view.Context, accessToken)
    tweets, err := t.Mentions()
    if err != nil {
        return nil, err
    }

    return struct {
        Tweets []*twitter.Tweet `json:"tweets"`
    }{tweets}, nil
})

var _ = ajax.Get("user:followers", "/user/followers", func(view *web.View) (interface{}, error) {
    accessToken, err := oauth.DecodeToken(view.Session.User.AccessToken)
    if err != nil {
        return nil, err
    }

    t := twitter.New(view.Context, accessToken)
    users, err := t.Followers(view.Session.User.ID())
    if err != nil {
        return nil, err
    }

    return struct {
        Users []*twitter.User `json:"users"`
    }{users}, nil
})

var _ = ajax.Get("user:friends", "/user/friends", func(view *web.View) (interface{}, error) {
    accessToken, err := oauth.DecodeToken(view.Session.User.AccessToken)
    if err != nil {
        return nil, err
    }

    t := twitter.New(view.Context, accessToken)
    users, err := t.Friends(view.Session.User.ID())
    if err != nil {
        return nil, err
    }

    return struct {
        Users []*twitter.User `json:"users"`
    }{users}, nil
})
