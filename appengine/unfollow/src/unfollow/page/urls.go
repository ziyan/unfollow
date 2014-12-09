package page

import (
    "unfollow/urls"
    "unfollow/web"
)

func init() {
    urls.Routers[0].Handle("/", web.MakeView(WelcomeView)).Methods("GET").Name("page:welcome")
    urls.Routers[0].Handle("/about", web.MakeView(AboutView)).Methods("GET").Name("page:about")
    urls.Routers[0].Handle("/network", web.MakeView(NetworkView)).Methods("GET").Name("page:network")
}
