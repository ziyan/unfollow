package user

import (
    "unfollow/api"
    "unfollow/web"
)

var _ = api.Get("user", "/user", func(handler *web.Handler) (interface{}, error) {
    return nil, nil
})
