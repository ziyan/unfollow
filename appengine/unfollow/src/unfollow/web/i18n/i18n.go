package i18n

import (
    "appengine"
    "fmt"
    "github.com/gorilla/i18n/gettext"
    "net/http"
    "os"
    "path"
    "unfollow/settings"
    "unfollow/web/sessions"
)

const (
    COOKIE_NAME = "locale"
    HEADER_NAME = "Accept-Language"
)

type Translation struct {
    Locale  string
    catalog *gettext.Catalog
}

func (translation *Translation) GetText(key string, args ...interface{}) string {
    return fmt.Sprintf(translation.catalog.Singular(key), args...)
}

var catalogs = make(map[string]*gettext.Catalog)

func load(locale string) error {
    file, err := os.Open(path.Join("locale", locale, "LC_MESSAGES", "unfollow.mo"))
    if err != nil {
        return err
    }
    defer file.Close()

    catalog := gettext.NewCatalog()
    err = catalog.ReadMo(file)
    if err != nil {
        return err
    }

    catalogs[locale] = catalog
    return nil
}

func init() {
    for locale, _ := range settings.LOCALES {
        if err := load(locale); err != nil {
            // panic(err)
        }
    }
}

func pick(context appengine.Context, request *http.Request, session *sessions.Session) (string, *gettext.Catalog) {
    // try using cookie first
    if cookie, _ := request.Cookie(COOKIE_NAME); cookie != nil {
        if catalog := catalogs[cookie.Value]; catalog != nil {
            if session.Locale != cookie.Value {
                session.Locale = cookie.Value
                session.Save()
            }
            return cookie.Value, catalog
        }
    }

    // try session
    if session.Locale != "" {
        if catalog := catalogs[session.Locale]; catalog != nil {
            return session.Locale, catalog
        }
    }

    // use default
    session.Locale = settings.DEFAULT_LOCALE
    session.Save()
    return session.Locale, catalogs[session.Locale]
}

func Do(context appengine.Context, request *http.Request, session *sessions.Session) *Translation {
    locale, catalog := pick(context, request, session)
    if catalog == nil {
        panic("unable to find i18n catalog.")
    }
    return &Translation{locale, catalog}
}

func Get(locale string) *Translation {
    catalog := catalogs[locale]
    if catalog == nil {
        panic("unable to find i18n catalog.")
    }
    return &Translation{locale, catalog}
}
