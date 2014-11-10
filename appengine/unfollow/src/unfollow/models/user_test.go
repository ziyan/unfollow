package models

import (
    "appengine/aetest"
    "unfollow/utils/cache"
    "unfollow/utils/db"
    "testing"
)

func TestPutUser(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    cache := cache.New(context)
    defer cache.Clear()

    db := db.New(context, cache)

    data := &User{Name: "Name", Bio: "Bio", Username: "username", Avatar: "", AccessToken: "access_token"}
    user1, err := PutUser(db, UserKey(db, 1), data)
    if err != nil {
        t.Fatal(err)
    }

    if !UserKey(db, 1).Equal(user1.Key()) {
        t.Fatal(user1.Key())
    }

    if user1.Username != "username" {
        t.Fatal(user1.Username)
    }

    user, err := GetUserByUsername(db, "username")
    if err != nil {
        t.Fatal(err)
    }

    if user == nil {
        t.Fatal(user)
    }

    if !UserKey(db, 1).Equal(user.Key()) {
        t.Fatal(user.Key())
    }

    user2, err := PutUser(db, UserKey(db, 2), data)
    if err != nil {
        t.Fatal(err)
    }

    if !UserKey(db, 2).Equal(user2.Key()) {
        t.Fatal(user2.Key())
    }

    if user2.Username != "username" {
        t.Fatal(user2.Username)
    }

    user, err = GetUserByUsername(db, "username")
    if err != nil {
        t.Fatal(err)
    }

    if user == nil {
        t.Fatal(user)
    }

    if !UserKey(db, 2).Equal(user.Key()) {
        t.Fatal(user.Key())
    }

    user, err = GetUserByID(db, 1)
    if err != nil {
        t.Fatal(err)
    }

    if user.Username != "" {
        t.Fatal(user.Username)
    }
}
