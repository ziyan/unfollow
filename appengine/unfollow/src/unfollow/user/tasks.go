package user

import (
    "appengine/datastore"
    "appengine/taskqueue"
    "github.com/ziyan/oauth"
    "unfollow/models"
    "unfollow/task"
    "unfollow/urls"
    "unfollow/utils/twitter"
    "unfollow/web"
)

var _ = task.Handle("user:pool", "/user/pool", func(handler *web.Handler) (interface{}, error) {

    // build query
    query := datastore.NewQuery(models.USER_KIND).Order("__key__").Limit(10)
    cursor := handler.Request.URL.Query().Get("cursor")
    if cursor != "" {
        handler.Context.Infof("user: cursor = %s", cursor)
        start, err := datastore.DecodeCursor(cursor)
        if err != nil {
            return nil, err
        }
        query = query.Start(start)
    }

    schedule := false

    // go over the list of users
    iterator := query.Run(handler.Context)
    for {
        user := &models.User{}
        key, err := iterator.Next(user)
        if err == datastore.Done {
            break
        }
        if err != nil {
            return nil, err
        }
        user.SetKey(key)
        user.SetOk(true)

        // decode access token
        accessToken, err := oauth.DecodeToken(user.AccessToken)
        if err != nil {
            return nil, err
        }

        // pool access token
        if err := twitter.PoolAccessToken(handler.Context, accessToken); err != nil {
            return nil, err
        }

        schedule = true
    }

    if !schedule {
        return nil, nil
    }

    // figure out end cursor
    end, err := iterator.Cursor()
    if err != nil {
        return nil, err
    }

    // schedule callback for next batch
    url := urls.Reverse("task:user:pool", "?", "cursor", end.String())
    if _, err := taskqueue.Add(handler.Context, &taskqueue.Task{Path: url.Path + "?" + url.RawQuery, Method: "POST"}, "default"); err != nil {
        return nil, err
    }

    return nil, nil
})
