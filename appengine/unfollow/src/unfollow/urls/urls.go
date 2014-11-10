package urls

import (
    "unfollow/settings"
    "github.com/gorilla/mux"
    "net/http"
    "net/url"
)

// Global router.
// All packages should register their url to this router.
var Router = mux.NewRouter().StrictSlash(true)
var Routers = make([]*mux.Router, 2)

func init() {
    for i := 0; i < len(Routers); i++ {
        Routers[i] = Router.NewRoute().Subrouter()
    }
    http.Handle("/", Router)
}

func Reverse(name string, args ...string) *url.URL {
    query := args[len(args):len(args)]
    pairs := args[0:len(args)]
    for i := 0; i < len(args); i += 2 {
        if args[i] == "?" {
            pairs = args[0:i]
            query = args[i+1 : len(args)]
            break
        }
    }

    u, err := Router.Get(name).URL(pairs...)
    if err != nil {
        panic("urls: route not found")
    }

    if len(query) > 0 {
        q := url.Values{}
        for i := 0; i < len(query); i += 2 {
            q.Add(query[i], query[i+1])
        }
        u.RawQuery = q.Encode()
    }

    return u
}

func ReverseByRequest(request *http.Request, name string, args ...string) *url.URL {
    url := Reverse(name, args...)
    url.Host = request.Host
    url.Scheme = "http"
    if request.TLS != nil {
        url.Scheme = "https"
    }
    return url
}

func ReverseBySettings(name string, args ...string) *url.URL {
    url := Reverse(name, args...)
    url.Host = settings.URL.Host
    url.Scheme = settings.URL.Scheme
    return url
}
