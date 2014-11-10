package user

import (
    "unfollow/urls"
    "unfollow/web"
)

func init() {
    urls.Routers[0].Handle("/user/login", web.MakeView(LoginView)).Methods("GET").Name("user:login")
    urls.Routers[0].Handle("/user/logon", web.MakeView(LogonView)).Methods("GET").Name("user:logon")
    urls.Routers[0].Handle("/user/logout", web.MakeView(LogoutView)).Methods("GET").Name("user:logout")

    urls.Routers[0].Handle("/user/welcome", web.MakeView(WelcomeView)).Methods("GET").Name("user:welcome")
    urls.Routers[0].Handle("/user/{id:[0-9]+}", web.MakeView(UserView)).Methods("GET").Name("user:user:id")
    urls.Routers[1].Handle("/{username:[0-9a-zA-Z\\-]+}", web.MakeView(UserView)).Methods("GET").Name("user:user:username")
}
