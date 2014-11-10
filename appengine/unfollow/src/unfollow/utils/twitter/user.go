package twitter

import (
    "net/url"
)

type User struct {
    ID       int64  `json:"id"`
    ScreenName string `json:"screen_name"`
    Name     string `json:"name"`
    Description      string `json:"description"`
    ProfileImageUrlHttps string `json:"profile_image_url_https"`
}


func (twitter *Twitter) VerifyCredentials() (*User, error) {
    user := User{}
    values := url.Values{
        "skip_status": {"true"},
        "include_entities": {"false"},
    }
    if err := twitter.Get("/1.1/account/verify_credentials.json", values, &user); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: user: %v", user)
    return &user, nil
}
