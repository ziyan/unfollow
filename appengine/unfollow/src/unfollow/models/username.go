package models

import (
    "appengine/datastore"
    "strings"
    "unfollow/utils/db"
)

const (
    USERNAME_KIND = "username"
)

type Username struct {
    key *datastore.Key `datastore:"-"`
    ok  bool           `datastore:"-"`

    UserKey *datastore.Key `datastore:"user_key,noindex"`
}

func (u *Username) Key() *datastore.Key {
    return u.key
}

func (u *Username) SetKey(key *datastore.Key) {
    u.key = key
}

func (u *Username) Ok() bool {
    return u.ok
}

func (u *Username) SetOk(ok bool) {
    u.ok = ok
}

func UsernameKey(db *db.Database, username string) *datastore.Key {
    return db.Key(USERNAME_KIND, strings.ToLower(username), 0, nil)
}
