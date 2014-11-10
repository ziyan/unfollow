package models

import (
    "appengine/datastore"
    "unfollow/urls"
    "unfollow/utils/db"
    "net/url"
    "strconv"
    "strings"
)

const (
    USER_KIND = "user"
)

type User struct {
    key *datastore.Key `datastore:"-"`
    ok  bool           `datastore:"-"`

    Name string `datastore:"name,noindex"`
    Bio  string `datastore:"bio,noindex"`

    Username string `datastore:"username,noindex"`
    Avatar   string `datastore:"avatar,noindex"`

    AccessToken string `datastore:"access_token,noindex"`
}

func (u *User) Key() *datastore.Key {
    return u.key
}

func (u *User) SetKey(key *datastore.Key) {
    u.key = key
}

func (u *User) Ok() bool {
    return u.ok
}

func (u *User) SetOk(ok bool) {
    u.ok = ok
}

func (u *User) ID() int64 {
    return u.key.IntID()
}

func (u *User) URL() *url.URL {
    if u.Username != "" {
        return urls.Reverse("user:user:username", "username", u.Username)
    }
    return urls.Reverse("user:user:id", "id", strconv.FormatInt(u.ID(), 10))
}

func UserKey(db *db.Database, id int64) *datastore.Key {
    return db.Key(USER_KIND, "", id, nil)
}

func GetUser(db *db.Database, key *datastore.Key) (*User, error) {
    user := &User{}
    user.SetKey(key)
    if err := db.Get(user, nil); err != nil {
        return nil, err
    }
    if !user.Ok() {
        return nil, nil
    }
    return user, nil
}

func GetUserByID(db *db.Database, id int64) (*User, error) {
    return GetUser(db, UserKey(db, id))
}

func GetUserByUsername(db *db.Database, username string) (*User, error) {
    // lookup username
    u := &Username{}
    u.SetKey(UsernameKey(db, username))
    if err := db.Get(u, nil); err != nil {
        return nil, err
    }
    if !u.Ok() {
        return nil, nil
    }

    // lookup user
    return GetUser(db, u.UserKey)
}

func PutUser(d *db.Database, key *datastore.Key, data *User) (*User, error) {
    user := &User{}
    user.SetKey(key)
    if err := d.Transaction(func(db *db.Database) error {
        // get existing user
        if err := db.Get(user, nil); err != nil {
            return err
        }

        if user.Ok() {
            if user.Name == data.Name &&
                user.Bio == data.Bio &&
                user.Avatar == data.Avatar &&
                user.Username == data.Username &&
                user.AccessToken == data.AccessToken {
                // nothing changed
                return nil
            }

            // username changed?
            // need to release existing username
            if user.Username != "" && strings.ToLower(user.Username) != strings.ToLower(data.Username) {
                if err := db.Delete(UsernameKey(db, user.Username), nil); err != nil {
                    return err
                }
            }
        }

        // update user
        user.Name = data.Name
        user.Bio = data.Bio
        user.Avatar = data.Avatar
        user.Username = data.Username
        user.AccessToken = data.AccessToken

        // save the user
        if err := db.Put(user, nil); err != nil {
            return err
        }

        // no username needed
        if user.Username == "" {
            return nil
        }

        // get existing username
        username := &Username{}
        username.SetKey(UsernameKey(db, user.Username))
        if err := db.Get(username, nil); err != nil {
            return err
        }

        if username.Ok() && !username.UserKey.Equal(user.Key()) {
            // username belonged to someone else
            // in this case, we need to remove the username first
            e := &User{}
            e.SetKey(username.UserKey)
            if err := db.Get(e, nil); err != nil {
                return err
            }

            if !e.Ok() {
                panic("models: username previous owner should exist")
            }

            e.Username = ""
            if err := db.Put(e, nil); err != nil {
                return err
            }
        }

        // save username
        if !user.Key().Equal(username.UserKey) {
            username.UserKey = user.Key()
            if err := db.Put(username, nil); err != nil {
                return err
            }
        }

        return nil
    }, true); err != nil {
        return nil, err
    }

    return user, nil
}

