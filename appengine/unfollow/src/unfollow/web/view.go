package web

import (
    "appengine"
    "unfollow/settings"
    "unfollow/utils/cache"
    "unfollow/utils/db"
    "unfollow/web/csrf"
    "unfollow/web/i18n"
    "unfollow/web/sessions"
    "unfollow/web/templates"
    "github.com/gorilla/mux"
    "github.com/mjibson/appstats"
    "net/http"
)

type View struct {
    Handler

    Session     *sessions.Session
    Translation *i18n.Translation
    Settings    struct {
        DEBUG     bool
        VERSION   string
        STATIC    string
        ANALYTICS string
        SOCKET    string
        LOCALES   map[string]string
    }
    Title string
}

type ViewFunction func(*View) error

func MakeView(function ViewFunction) appstats.Handler {
    return appstats.NewHandler(func(context appengine.Context, response http.ResponseWriter, request *http.Request) {
        view := new(View)

        // initialize handler
        view.Context = context
        view.Request = request
        view.Response = response
        view.Variables = mux.Vars(request)

        // copy over some important settings exposed to the template system
        view.Settings.DEBUG = settings.DEBUG
        view.Settings.VERSION = settings.VERSION
        view.Settings.STATIC = settings.STATIC
        view.Settings.ANALYTICS = settings.ANALYTICS
        view.Settings.LOCALES = settings.LOCALES

        // headers
        view.configureHeaders()

        // cache
        view.Cache = cache.New(context)
        defer view.Cache.Clear()

        // db
        view.Database = db.New(context, view.Cache)

        // session
        view.Session = sessions.Load(context, view.Cache, view.Database, request, response)
        defer sessions.Save(context, view.Cache, view.Database, request, view.Session)

        // csrf
        if err := csrf.Do(context, request, response, view.Session); err != nil {
            return
        }

        // i18n
        view.Translation = i18n.Do(context, request, view.Session)

        // page title
        view.Title = view.GetText("Unfollow")

        // set up panic handling
        defer RecoverError(context, request, response)

        // do the work
        err := function(view)

        // handle error
        HandleError(context, request, response, err)
    })
}

func (view *View) Render(data interface{}, filenames ...string) error {
    view.Data = data
    return templates.Render(view.Response, view, filenames...)
}

func (view *View) RenderToString(data interface{}, filenames ...string) (string, error) {
    view.Data = data
    return templates.RenderToString(view, filenames...)
}

func (view *View) LoginRequired() error {
    if view.Session.User == nil {
        return ErrLoginRequired
    }
    return nil
}

func (view *View) GetText(key string, args ...interface{}) string {
    return view.Translation.GetText(key, args...)
}

func (view *View) GetLocalized(str map[string]string) string {
    if s, ok := str[view.Translation.Locale]; ok {
        return s
    }
    if s, ok := str[view.Translation.Locale[:2]]; ok {
        return s
    }
    if s, ok := str[settings.DEFAULT_LOCALE]; ok {
        return s
    }
    if s, ok := str[settings.DEFAULT_LOCALE[:2]]; ok {
        return s
    }
    return ""
}

var viewHeaders = map[string]string{
    "Content-Security-Policy": "default-src *; script-src 'self' 'unsafe-eval' www.google-analytics.com; style-src 'self' 'unsafe-inline' fonts.googleapis.com",
    "X-Frame-Options":         "SAMEORIGIN",
    "X-Xss-Protection":        "1; mode=block",
    "X-Content-Type-Options":  "nosniff",
    "Cache-Control":           "private, no-cache, no-store, must-revalidate",
}

func (view *View) configureHeaders() {
    // setting some security related headers
    if settings.SECURE {
        view.Response.Header().Set("Strict-Transport-Security", "max-age=31536000")
    }

    for key, value := range viewHeaders {
        view.Response.Header().Set(key, value)
    }
}
