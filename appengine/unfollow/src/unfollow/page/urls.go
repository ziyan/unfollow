package page

import (
    "unfollow/urls"
    "unfollow/web"
)

func init() {
    urls.Routers[0].Handle("/", web.MakeView(WelcomeView)).Methods("GET").Name("page:welcome")
    urls.Routers[0].Handle("/about", web.MakeView(AboutView)).Methods("GET").Name("page:about")
    urls.Routers[0].Handle("/home", web.MakeView(HomeView)).Methods("GET").Name("page:home")
}
