package user

import (
    "github.com/ziyan/oauth"
    "strings"
    "unfollow/models"
    "unfollow/urls"
    "unfollow/utils/twitter"
    "unfollow/web"
)

const (
    sessionRequestTokenKey = "request_token"
    sessionCallbackKey     = "callback"
)

func createLogonUrl(view *web.View, next string) (string, error) {
    callback := ""
    if next != "" && strings.HasPrefix(next, "/") {
        callback = urls.ReverseByRequest(view.Request, "user:logon", "?", "next", next).String()
    } else {
        callback = urls.ReverseByRequest(view.Request, "user:logon").String()
    }

    requestToken, err := twitter.GetRequestToken(view.Context, callback)
    if err != nil {
        return "", err
    }

    url, err := twitter.CreateAuthorizeUrl(requestToken)
    if err != nil {
        return "", err
    }

    // save stuff to session
    view.Session.Values[sessionRequestTokenKey] = requestToken.Encode()
    view.Session.Values[sessionCallbackKey] = callback
    view.Session.Save()

    return url, nil
}

func handleLogon(view *web.View) error {
    rawToken, ok := view.Session.Values[sessionRequestTokenKey].(string)
    if !ok {
        return web.ErrBadRequest
    }

    requestToken, err := oauth.DecodeToken(rawToken)
    if err != nil {
        return err
    }

    if err := view.Request.ParseForm(); err != nil {
        return err
    }

    if requestToken.Key() != view.Request.FormValue("oauth_token") {
        return web.ErrBadRequest
    }

    verifier := view.Request.FormValue("oauth_verifier")
    if verifier == "" {
        // access denied
        return web.ErrBadRequest
    }

    accessToken, err := twitter.GetAccessToken(view.Context, requestToken, verifier)
    if err != nil {
        return err
    }

    t := twitter.New(view.Context, accessToken)
    user, err := t.VerifyCredentials()
    if err != nil {
        return err
    }

    // create or update the user
    u, err := models.PutUser(view.Database, models.UserKey(view.Database, user.ID), &models.User{
        Name:        user.Name,
        Bio:         user.Description,
        Username:    user.ScreenName,
        Avatar:      user.ProfileImageUrlHttps,
        AccessToken: accessToken.Encode(),
    })
    if err != nil {
        return err
    }

    // login the user
    view.Session.Abandon(view.Context, view.Cache, view.Database, view.Response)
    view.Session.User = u
    view.Session.Save()

    return nil
}
