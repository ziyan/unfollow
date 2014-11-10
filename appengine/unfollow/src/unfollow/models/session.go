package models

import (
    "appengine/datastore"
    "unfollow/utils/db"
    "time"
)

const (
    SESSION_KIND = "session"
)

type Session struct {
    key *datastore.Key `datastore:"-"`
    ok  bool           `datastore:"-"`

    UserKey   *datastore.Key `datastore:"user_key,noindex"`
    Locale    string         `datastore:"locale,noindex"`
    UserAgent string         `datastore:"user_agent,noindex"`
    IP        string         `datastore:"ip,noindex"`

    Time time.Time `datastore:"time,noindex"`
}

func (s *Session) Key() *datastore.Key {
    return s.key
}

func (s *Session) SetKey(key *datastore.Key) {
    s.key = key
}

func (s *Session) Ok() bool {
    return s.ok
}

func (s *Session) SetOk(ok bool) {
    s.ok = ok
}

func (s *Session) ID() string {
    return s.Key().StringID()
}

func SessionKey(db *db.Database, id string) *datastore.Key {
    return db.Key(SESSION_KIND, id, 0, nil)
}

func GetSession(db *db.Database, key *datastore.Key) (*Session, error) {
    session := &Session{}
    session.SetKey(key)
    if err := db.Get(session, nil); err != nil {
        return nil, err
    }
    if !session.Ok() {
        return nil, nil
    }
    return session, nil
}

func GetSessionByID(db *db.Database, id string) (*Session, error) {
    return GetSession(db, SessionKey(db, id))
}

func SaveSession(db *db.Database, key, userKey *datastore.Key, locale, userAgent, ip string) (*Session, error) {
    session := &Session{
        UserKey:   userKey,
        Locale:    locale,
        UserAgent: userAgent,
        IP:        ip,
        Time:      time.Now(),
    }
    session.SetKey(key)
    if err := db.Put(session, nil); err != nil {
        return nil, err
    }
    return session, nil
}

func DeleteSession(db *db.Database, key *datastore.Key) error {
    return db.Delete(key, nil)
}
