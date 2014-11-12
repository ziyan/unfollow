// TODO:
// - Renew session cookie
// - Clean up expired sessions
package sessions

import (
    "appengine"
    "appengine/datastore"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base64"
    "errors"
    "net/http"
    "strconv"
    "strings"
    "time"
    "unfollow/models"
    "unfollow/settings"
    "unfollow/utils/cache"
    "unfollow/utils/db"
    "unfollow/utils/security"
)

const (
    COOKIE_NAME    = "session"
    COOKIE_PATH    = "/"
    COOKIE_MAX_AGE = 14 * 24 * 60 * 60
    CACHE_PREFIX   = "session"
)

var (
    ErrNotFound = errors.New("session: not found")
)

type Session struct {
    id         string
    User       *models.User
    CSRFToken  string
    Locale     string
    Values     map[string]interface{}
    isModified bool
}

func (session *Session) Save() {
    session.isModified = true
}

func (session *Session) Abandon(context appengine.Context, cache *cache.Cache, db *db.Database, response http.ResponseWriter) {
    if session.User != nil {
        deleteFromDatabase(db, session.id)
    }

    deleteFromMemcache(cache, session.id)
    session.reset(cache, db, response)
}

func (session *Session) reset(cache *cache.Cache, db *db.Database, response http.ResponseWriter) {

    for {
        // loop until we find an unused session id
        session.id = security.GenerateRandomID()

        if err := loadFromCache(cache, db, session); err == nil {
            continue
        }

        if err := loadFromDatabase(db, session); err == nil {
            continue
        }

        break
    }

    // Locale and CSRFToken do not need to be reset
    session.Values = make(map[string]interface{})
    session.isModified = true
    session.User = nil

    http.SetCookie(response, &http.Cookie{
        Name:     COOKIE_NAME,
        Value:    encodeCookie(session.id),
        Path:     COOKIE_PATH,
        Expires:  time.Now().Add(time.Duration(COOKIE_MAX_AGE) * time.Second),
        HttpOnly: true,
        Secure:   settings.SECURE,
    })
}

func Load(context appengine.Context, cache *cache.Cache, db *db.Database, request *http.Request, response http.ResponseWriter) *Session {
    session := &Session{}

    if cookie, err := request.Cookie(COOKIE_NAME); err == nil {
        session.id = decodeCookie(cookie.Value)
    }

    if session.id != "" {
        if err := loadFromCache(cache, db, session); err == nil {
            return session
        }
        if err := loadFromDatabase(db, session); err == nil {
            return session
        }
    }

    session.reset(cache, db, response)
    return session
}

func Save(context appengine.Context, cache *cache.Cache, db *db.Database, request *http.Request, session *Session) {

    // only save session if it is changed
    if !session.isModified {
        return
    }

    saveToCache(cache, session)

    // only save session to datastore if it is logged in
    if session.User == nil {
        return
    }

    saveToDatabase(db, request, session)
}

func decodeCookie(value string) string {
    parts := strings.Split(value, ":")
    if len(parts) != 3 {
        return ""
    }

    if calculateSignature(parts[0], parts[1]) != parts[2] {
        return ""
    }

    timestamp, err := strconv.ParseInt(parts[1], 10, 64)
    if err != nil {
        return ""
    }

    now := time.Now().UTC().Unix()
    if timestamp > now || now-timestamp > COOKIE_MAX_AGE {
        return ""
    }

    return parts[0]
}

func encodeCookie(id string) string {
    timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
    signature := calculateSignature(id, timestamp)
    return strings.Join([]string{id, timestamp, signature}, ":")
}

func calculateSignature(id string, timestamp string) string {
    hasher := hmac.New(sha1.New, settings.SECRET)
    hasher.Write([]byte(id))
    hasher.Write([]byte(timestamp))
    signature := hasher.Sum(nil)
    return strings.TrimRight(base64.URLEncoding.EncodeToString(signature), "=")
}

type cacheSession struct {
    UserKey   *datastore.Key
    CSRFToken string
    Locale    string
    Values    map[string]interface{}
}

func loadFromCache(cache *cache.Cache, db *db.Database, session *Session) error {
    key := strings.Join([]string{CACHE_PREFIX, session.id}, ":")
    object := cacheSession{}

    exist, err := cache.Get(key, &object, nil)
    if err != nil {
        return err
    }
    if !exist {
        return ErrNotFound
    }

    if object.UserKey != nil {
        session.User, _ = models.GetUser(db, object.UserKey)
    }

    session.CSRFToken = object.CSRFToken
    session.Locale = object.Locale
    session.Values = object.Values
    session.isModified = false

    return nil
}

func saveToCache(cache *cache.Cache, session *Session) error {
    key := strings.Join([]string{CACHE_PREFIX, session.id}, ":")
    object := cacheSession{
        UserKey:   nil,
        CSRFToken: session.CSRFToken,
        Locale:    session.Locale,
        Values:    session.Values,
    }
    if session.User != nil {
        object.UserKey = session.User.Key()
    }

    return cache.Set(key, &object, nil)
}

func deleteFromMemcache(cache *cache.Cache, id string) error {
    key := strings.Join([]string{CACHE_PREFIX, id}, ":")
    return cache.Delete(key, nil)
}

func loadFromDatabase(db *db.Database, session *Session) error {
    entity, err := models.GetSessionByID(db, session.id)
    if err != nil {
        return err
    }
    if entity == nil {
        return ErrNotFound
    }

    session.User, _ = models.GetUser(db, entity.UserKey)
    session.CSRFToken = ""
    session.Locale = entity.Locale
    session.Values = make(map[string]interface{})
    session.isModified = true

    return nil
}

func saveToDatabase(db *db.Database, request *http.Request, session *Session) error {
    _, err := models.SaveSession(db, models.SessionKey(db, session.id), session.User.Key(), session.Locale, request.UserAgent(), request.RemoteAddr)
    return err
}

func deleteFromDatabase(db *db.Database, id string) error {
    return models.DeleteSession(db, models.SessionKey(db, id))
}
