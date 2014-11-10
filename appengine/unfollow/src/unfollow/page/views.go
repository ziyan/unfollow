package page

import (
    "unfollow/web"
)

func WelcomeView(view *web.View) error {
    return view.Render(nil, "html.html", "base.html", "page/welcome.html")
}

func AboutView(view *web.View) error {
    view.Title = view.GetText("About Unfollow")
    return view.Render(nil, "html.html", "base.html", "page/about.html")
}

func HomeView(view *web.View) error {
    if err := view.LoginRequired(); err != nil {
        return err
    }

    return view.Render(nil, "html.html", "base.html", "page/home.html")
}
